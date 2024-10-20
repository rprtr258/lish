` command-line interface abstractions
  for [cmd] [verb] [options] form`

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
each := std.each
slice := std.slice
str := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/str.ink')
hasPrefix? := str.hasPrefix?

maybeOpt := part => true :: {
  hasPrefix?(part, '--') -> slice(part, 2, len(part))
  hasPrefix?(part, '-') -> slice(part, 1, len(part))
  _ -> ()
}

`
Supports:
  -opt val
  --opt val
  -opt=val
  --opt val
all other values are considered args
`
parsed := () => (
  as := args()

  verb := as.2
  rest := slice(as, 3, len(as))

  opts := {}
  args := []

  s := {
    lastOpt: ()
    onlyArgs: false
  }
  each(rest, part => [maybeOpt(part), s.lastOpt] :: {
    [(), ()] -> (
      ` not opt, no prev opt `
      args.len(args) := part
    )
    [(), _] -> (
      ` not opt, prev opt exists `
      opts.(s.lastOpt) := part
      s.lastOpt := ()
    )
    [_, ()] -> (
      ` is opt, no prev opt `
      s.lastOpt := maybeOpt(part)
    )
    _ -> (
      ` is opt, prev opt exists `
      opts.(s.lastOpt) := true
      s.lastOpt := maybeOpt(part)
    )
  })

  s.lastOpt :: {
    () -> ()
    _ -> opts.(s.lastOpt) := true
  }

  {
    verb: verb
    opts: opts
    args: args
  }
)

{parsed: parsed}
