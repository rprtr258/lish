# cons :: []Doc -> Doc
# nil :: Doc # TODO: text('')
# text :: string -> Doc
# line :: Doc # TODO: text('\n') ?
# nest :: (int, Doc) -> Doc # TODO: rename to indent
# layout :: Doc -> string

string_impl := (
  {join, replace} := import('str.ink')
  repeat := (n, x) => n :: {0 -> '', _ -> x + repeat(n-1, x)}
  {
    cons: ss => join(ss, '')
    nil: ''
    text: s => s
    line: '\n'
    nest: (n, s) => replace(s, '\n', '\n' + repeat(n, ' '))
    layout: s => s
  }
)
{cons, nil, text, line, nest, layout} := string_impl

# Tree :: [string, []Tree]
showBracket := ts => ts :: {
  [] -> nil
  _ -> cons([text('['), nest(1, showTrees(ts)), text(']')])
}
#showTree := [s, ts] => cons([text(s), nest(len(s), showBracket(ts))])
showTree := t => (
  [s, ts] := t
  cons([text(s), nest(len(s), showBracket(ts))])
)
{slice} := import('std.ink')
showTrees := ts => true :: {
  len(ts) == 1 -> showTree(ts.0)
  _ -> cons([showTree(ts.0), text(','), line, showTrees(slice(ts, 1, len(ts)))])
}
l := s => [s, []] # leaf
exampleTree := ['aaa', [
  ['bbbbb', [
    l('ccc'),
    l('dd'),
  ]],
  l('eee'),
  ['ffff', [
    l('gg'),
    l('hhh'),
    l('ii')
  ]]
]]
out(showTree(exampleTree) + '\n')

showBracket2 := ts => ts :: {
  [] -> nil
  _ -> cons([text('['), nest(2, cons([line, showTrees2(ts)])), line, text(']')])
}
showTree2 := t => (
  [s, ts] := t
  cons([text(s), showBracket2(ts)])
)
showTrees2 := ts => true :: {
  len(ts) == 1 -> showTree2(ts.0)
  _ -> cons([showTree2(ts.0), text(','), line, showTrees2(slice(ts, 1, len(ts)))])
}
out(showTree2(exampleTree) + '\n')
