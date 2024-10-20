# september translate command

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
log := std.log
map := std.map
each := std.each
cat := std.cat

Tokenize := import('tokenize.ink')
tokenize := Tokenize.tokenize
tkString := Tokenize.tkString

Parse := import('parse.ink')
parse := Parse.parse
ndString := Parse.ndString

analyze := import('analyze.ink')

gen := import('gen.ink')

Newline := char(10)

main := prog => (
  tokens := tokenize(prog)
  # each(tokens, tok => log(tkString(tok)))

  nodes := parse(tokens)

  type(nodes) :: {
    ` tree of nodes `
    'composite' -> (
      # each(nodes, node => log(ndString(node)))
      analyzed := map(nodes, analyze)
      cat(map(analyzed, gen), ';' + Newline) + Newline
    )
    ` parse err `
    'string' -> nodes
  }
)

{main: main}
