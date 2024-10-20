` september translate command `

std := import('../vendor/std')

log := std.log
map := std.map
each := std.each
cat := std.cat

Tokenize := import('tokenize')
tokenize := Tokenize.tokenize
tkString := Tokenize.tkString

Parse := import('parse')
parse := Parse.parse
ndString := Parse.ndString

Analyze := import('analyze')
analyze := Analyze.analyze

Gen := import('gen')
gen := Gen.gen

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
