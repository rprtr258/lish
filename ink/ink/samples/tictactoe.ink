# interactive terminal tic tac toe in Ink

{mod, styled} := import('styled.ink')
{scan, slice} := import('std.ink')
{format: f} := import('str.ink')
{map, mapi, filter} := import('functional.ink')

log := s => out(s + '\n')

# async version of a while(... condition, ... predicate)
# that takes a callback
asyncWhile := (cond, do) => (sub := () => true :: {
  cond() -> do(sub)
  _ -> ()
})()

# shorthand tools for getting players and player labels
Player := {x: 1, o: 2}
Label := [' ', 'x', 'o']
# make letters appear bolder / fainter on the board
bold := c => styled(c, [mod.bold])
grey := c => styled(c, [mod.fg.yellow, mod.faint])

# create a new game board + state
newBoard := () => [
  1 # current player turn
  0, 0, 0
  0, 0, 0
  0, 0, 0
]

# format string to print board state
BoardFormat := '{{ 1 }} │ {{ 2 }} │ {{ 3 }}
──┼───┼──
{{ 4 }} │ {{ 5 }} │ {{ 6 }}
──┼───┼──
{{ 7 }} │ {{ 8 }} │ {{ 9 }}
'
# format-print board state
stringBoard := bd => f(
  BoardFormat
  mapi(bd, (player, idx) => Label.(player) :: {
    ' ' -> grey(string(idx))
    _ -> bold(Label.(player))
  })
)

# winning placement combinations for a single player
Combinations := [
  # horizontal
  [1, 2, 3]
  [4, 5, 6]
  [7, 8, 9]

  # vertical
  [1, 4, 7]
  [2, 5, 8]
  [3, 6, 9]

  # diagonal
  [1, 5, 9]
  [3, 5, 7]
]
# returns -1 if no win, 0 if tie, or winner player ID
Result := {
  None: ~1
  Tie: 0
  X: Player.x
  O: Player.o
}
checkBoard := bd => (
  checkIfPlayerWon := player => (
    isPlayer := row => row == [player, player, player]
    possibleRows := map(Combinations, combo => map(combo, idx => bd.(idx)))
    didWin := len(filter(possibleRows, (row, _) => isPlayer(row))) > 0

    didWin
  )

  true :: {
    checkIfPlayerWon(Player.x) -> Result.X
    checkIfPlayerWon(Player.o) -> Result.O
    _ -> (
      # check if game ended in a tie
      takenCells := filter(slice(bd, 1, 10), (val, _) => ~(val == 0))
      len(takenCells) :: {
        9 -> Result.Tie
        _ -> Result.None
      }
    )
  }
)

# take one player turn, mutates game state
stepBoard! := (bd, cb) => scan(s => idx := number(s) :: {
  # not a number, try again
  () -> stepBoard!(bd, cb)
  _ -> true :: {
    # number in range, make a move
    idx > 0 & idx < 10 -> bd.(idx) :: {
      # the given cell is empty, make a move
      0 -> (
        bd.(number(s)) := getPlayer(bd)
        setPlayer(bd, nextPlayer(bd))
        cb()
      )
      # the cell is already occupied, try again
      _ -> (
        log(f('{{ idx }} is already taken!', {idx}))
        out(f('Move for player {{ player }}: ', {
          player: Label.(getPlayer(bd))
        }))
        stepBoard!(bd, cb)
      )
    }
    # number not in range, try again
    _ -> (
      log('Enter a number 0 < n < 10.')
      out(f('Move for player {{ player }}: ', {
        player: Label.(getPlayer(bd))
      }))
      stepBoard!(bd, cb)
    )
  }
})

# get/set/modify player turn state from the game board
getPlayer := bd => bd.0
setPlayer := (bd, pl) => bd.0 := pl
nextPlayer := bd => Label.(getPlayer(bd)) :: {
  'x' -> Player.o
  _ -> Player.x
}

# divider used to delineate each turn in the UI
Divider := '
>---------------<
'

# run a single game
log('Welcome to Ink tic-tac-toe!')
bd := newBoard()
asyncWhile(
  () => checkBoard(bd) :: {
    Result.None -> true
    _ -> (
      log(Divider)
      checkBoard(bd) :: {
        Result.Tie -> log('x and o tied!')
        Result.X -> log('x won!')
        Result.O -> log('o won!')
      }
      log('')
      log(stringBoard(bd))

      false
    )
  }
  cb => (
    log(Divider)
    log(stringBoard(bd))
    out(f('Move for player {{ player }}: ', {
      player: Label.(getPlayer(bd))
    }))
    stepBoard!(bd, cb)
  )
)
