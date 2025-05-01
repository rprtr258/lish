# september translate command

{log, map, each, cat} := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
{tokenize, tkString} := import('tokenize.ink')
{parse, ndString} := import('parse.ink')
analyze := import('analyze.ink')
gen := import('gen.ink')

prog => (
  tokens := tokenize(prog)
  # each(tokens, tok => log(tkString(tok)))

  nodes := parse(tokens)

  type(nodes) :: {
    # tree of nodes
    'composite' -> (
      # each(nodes, node => log(ndString(node)))
      analyzed := map(nodes, analyze)
      cat(map(analyzed, gen), ';\n') + '\n'
    )
    # parse err
    'string' -> nodes
  }
)
