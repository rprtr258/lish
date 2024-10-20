# JSON serde

map := import('functional.ink').map
{join: cat, ws?: ws?, digit?: digit?} := import('str.ink')

# string escape '"'
esc := c => point(c) :: {
  9 -> '\\t'
  10 -> '\\n'
  13 -> '\\r'
  34 -> '\\"'
  92 -> '\\\\'
  _ -> c
}
escape := s => (
  max := len(s)
  (sub := (i, acc) => i :: {
    max -> acc
    _ -> sub(i + 1, acc + esc(s.(i)))
  })(0, '')
)

# is this character a numeral digit or .?
num? := c => c :: {
  '' -> false
  '.' -> true
  _ -> digit?(c)
}

# reader implementation with internal state for deserialization
reader := s => (
  state := {
    idx: 0
    # has there been a parse error?
    err?: false
  }

  next := () => (
    state.idx := state.idx + 1
    c := s.(state.idx - 1) :: {
      () -> ''
      _ -> c
    }
  )

  peek := () => c := s.(state.idx) :: {
    () -> ''
    _ -> c
  }

  {
    next: next
    peek: peek
    # fast-forward through whitespace
    ff: () => (sub := () => ws?(peek()) :: {
      true -> (
        state.idx := state.idx + 1
        sub()
      )
    })()
    done?: () => ~(state.idx < len(s))
    err: () => state.err? := true
    err?: () => state.err?
  }
)

# deserialize string
deString := r => (
  next := r.next # TODO: remove this and similar
  peek := r.peek

  # known to be a '"'
  next()

  (sub := acc => peek() :: {
    '' -> (
      (r.err)()
      ()
    )
    '\\' -> (
      # eat backslash
      next()
      sub(acc + (c := next() :: {
        't' -> '\t'
        'n' -> '\n'
        'r' -> '\r'
        '"' -> '"'
        _ -> c
      }))
    )
    '"' -> (
      next()
      acc
    )
    _ -> sub(acc + next())
  })('')
)

# deserialize number
deNumber := r => (
  next := r.next
  peek := r.peek
  state := {
    # have we seen a '.' yet?
    negate?: false
    decimal?: false
  }

  peek() :: {
    '-' -> (
      next()
      state.negate? := true
    )
  }

  result := (sub := acc => num?(peek()) :: {
    true -> peek() :: {
      '.' -> state.decimal? :: {
        true -> (r.err)()
        _ -> (
          state.decimal? := true
          sub(acc + next())
        )
      }
      _ -> sub(acc + next())
    }
    _ -> acc
  })('')

  state.negate? :: {
    true -> ~number(result)
    _ -> number(result)
  }
)

# deserialize null
deNull := r => (
  next := r.next
  next() + next() + next() + next() :: {
    'null' -> ()
    _ -> (r.err)()
  }
)

# deserialize boolean
deTrue := r => (
  next := r.next
  next() + next() + next() + next() :: {
    'true' -> true
    _ -> (r.err)()
  }
)
deFalse := r => (
  next := r.next
  next() + next() + next() + next() + next() :: {
    'false' -> false
    _ -> (r.err)()
  }
)

# deserialize list
deList := r => (
  next := r.next
  peek := r.peek
  ff := r.ff
  state := {
    idx: 0
  }

  # known to be a '['
  next()
  ff()

  (sub := acc => (r.err?)() :: {
    true -> ()
    _ -> peek() :: {
      '' -> (
        (r.err)()
        ()
      )
      ']' -> (
        next()
        acc
      )
      _ -> (
        acc.(state.idx) := der(r)
        state.idx := state.idx + 1

        ff()
        peek() :: {
          ',' -> next()
        }

        ff()
        sub(acc)
      )
    }
  })([])
)

# deserialize composite
deComp := r => (
  next := r.next
  peek := r.peek
  ff := r.ff

  # known to be a '{'
  next()
  ff()

  (sub := acc => (r.err?)() :: {
    true -> ()
    _ -> peek() :: {
      '' -> (r.err)()
      '}' -> (
        next()
        acc
      )
      _ -> (
        key := deString(r)

        (r.err?)() :: {
          false -> (
            ff()
            peek() :: {
              ':' -> next()
            }

            ff()
            val := der(r)

            (r.err?)() :: {
              false -> (
                ff()
                peek() :: {
                  ',' -> next()
                }

                ff()
                acc.(key) := val
                sub(acc)
              )
            }
          )
        }
      )
    }
  })({})
)

# JSON string in reader to composite
der := r => (
  # trim preceding whitespace
  (r.ff)()

  result := ((r.peek)() :: {
    'n' -> deNull(r)
    '"' -> deString(r)
    't' -> deTrue(r)
    'f' -> deFalse(r)
    '[' -> deList(r)
    '{' -> deComp(r)
    _ -> deNumber(r)
  })

  # if there was a parse error, just return null result
  (r.err?)() :: {
    true -> ()
    _ -> result
  }
)

# JSON string to composite
parse := s => der(reader(s))

# composite to JSON string
serialize := c => type(c) :: {
  '()' -> 'null'
  'string' -> '"' + escape(c) + '"'
  'number' -> string(c)
  'boolean' -> c :: {
    true -> 'true'
    _ -> 'false'
  }
  'function' -> 'null' # do not serialize functions
  'composite' -> '{' + cat(map(keys(c), k => '"' + escape(k) + '":' + serialize(c.(k))), ',') + '}'
}

{
  de: parse # TODO: remove
  ser: serialize # TODO: remove
}
