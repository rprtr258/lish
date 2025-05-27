{slice} := import('std.ink')
{flatten, reduce, map, each} := import('functional.ink')
{split} := import('str.ink')

# implementation of [A prettier printer](https://homepages.inf.ed.ac.uk/wadler/papers/prettier/prettier.pdf)
# cons :: []Doc -> Doc
# nil :: Doc # TODO: text('')
# text :: string -> Doc
# line :: Doc # TODO: text('\n') ?
# nest :: (int, Doc) -> Doc # TODO: rename to indent
# layout :: Doc -> string
# group :: Doc -> Doc
#   union :: (Doc, Doc) -> Doc
#   flatten :: Doc -> Doc
# pretty :: (int, Doc) -> string

dbg := (x) => (out(string(x)+'\n'), x)

{cons, nil, text, line, nest, layout, group, pretty} := (
  # Doc :: Nil | Text (string, Doc) | Line (int, Doc)
  Nil  :=           {kind: 'nil'}
  Text := (s, x) => {kind: 'text', s, x}
  Line := (i, x) => {kind: 'line', i, x}

  # DOC :: NIL | CONS (DOC, DOC) | NEST (int, DOC) | TEXT string | LINE | UNION (DOC, DOC)
  NIL   :=           {kind: 'NIL'}
  CONS  := (l, r) => {kind: 'CONS', l, r}
  NEST  := (i, x) => {kind: 'NEST', i, x}
  TEXT  := s(      )=> {kind: 'TEXT', s}
  LINE  :=           {kind: 'LINE'}
  UNION := (l, r) => {kind: 'UNION', l, r}

  repeat := (n, x) => n :: {0 -> '', _ -> x + repeat(n-1, x)}
  # group := (x) => flatten(UNION(x, x))
  flatten := (x) => x.kind :: {
    'NIL' -> NIL
    'CONS' -> CONS(flatten(x.l), flatten(x.r))
    'NEST' -> NEST(x.i, flatten(x.x))
    'TEXT' -> x
    'LINE' -> TEXT(' ')
    'UNION' -> flatten(x.l)
  }
  # fits : (int, Doc) -> bool
  fits := (w, d) => true :: {
    w < 0 -> false,
    _ -> d.kind :: {
      'nil' -> true
      'text' -> ({s, x} := d, fits(w-len(s), x))
      'line' -> true
    }
  }
  better := (w, k, x, y) => true :: {fits(w-k, x) -> x, _ -> y}
  # be : (int, int, DOC) -> Doc
  be := (w, k, z) => true :: {
    len(z) == 0 -> Nil
    _ -> (
      zt := slice(z, 1, len(z)),
      [i, zz] := z.0,
      zz.kind :: {
        'NIL'   -> be(w, k, zt)
        'CONS'  -> be(w, k, [[i, zz.l], [i, zz.r]]+zt)
        'NEST'  -> be(w, k, [[i+zz.i, zz.x]]+zt)
        'TEXT'  -> Text(zz.s, be(w, k+len(zz.s), zt))
        'LINE'  -> Line(i, be(w, i, zt))
        'UNION' -> better(w, k, be(w, k, [[i, zz.l]]+zt),
                                be(w, k, [[i, zz.r]]+zt))
      })
  }
  best := (w, k, x) => be(w, k, [[0, x]])
  # layout : Doc -> string
  layout := (s) => s.kind :: {
    'nil' -> ''
    'text' -> ({s, x} := s, s + layout(x))
    'line' -> ({i, x} := s, '\n'+repeat(i, ' ')+layout(x))
  }
  {
    nil: NIL
    cons: (ss) => reduce(ss, (acc, x, _) => CONS(acc, x), NIL),
    text: TEXT
    line: LINE
    nest: NEST
    layout: (x) => layout(best(0, 0, x))
    group: (x) => UNION(flatten(x), x)
    pretty: (w, x) => layout(best(w, 0, x))
  }
)

# utility funcs
union := (l, r) => {kind: 'UNION', l, r} # TODO: hide
concat := (xs) => type(xs.0) :: {'list' -> flatten(xs), _ -> xs}
folddoc := (f, xs) => len(xs) :: {
  0 -> nil
  1 -> xs.0
  _ -> f(x, folddoc(f, slice(xs, 1, len(xs))))
}
spread := (xs) => folddoc((x, y) => cons([x, text(' '), y]), xs)
stack  := (xs) => folddoc((x, y) => cons([x,      line, y]), xs)
bracket := (l, x, r) => group(cons([
  text(l),
  nest(2, cons([line, x])), line,
  text(r),
]))
words := (s) => split(s, ' ')
# words := split(s, ' \n\t')
# fillwords := (s) => s | words | map(text) | folddoc(...)
fillwords := (s) => folddoc((x, y) => cons([x, group(line), y]), map(text, words))
fill := (xs) => len(xs) :: {
  0 -> nil
  1 -> xs.0
  _ -> union(
    spread([concat(xs), fill([concat(xs.1)]+slice(xs, 2, len(xs)))])
    stack([xs.0, fill(slice(xs, 1, len(xs)))])
  )
}

# Tree :: [string, []Tree]
showBracket := (ts) => ts :: {
  [] -> nil
  _ -> cons([text('['), nest(1, showTrees(ts)), text(']')])
}
#showTree := ([s, ts]) => cons([text(s), nest(len(s), showBracket(ts))])
showTree := (t) => (
  [s, ts] := t
  group(cons([text(s), nest(len(s), showBracket(ts))]))
)
showTrees := (ts) => true :: {
  len(ts) == 1 -> showTree(ts.0)
  _ -> cons([showTree(ts.0), text(','), line, showTrees(slice(ts, 1, len(ts)))])
}
l := (s) => [s, []] # leaf
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
out(pretty(30, showTree(exampleTree)) + '\n')

showTree2 := (t) => (
  [s, ts] := t
  cons([text(s), ts :: {
    [] -> nil
    _ -> bracket('[', showTrees2(ts), ']')
  }])
)
showTrees2 := (ts) => true :: {
  len(ts) == 1 -> showTree2(ts.0)
  _ -> cons([showTree2(ts.0), text(','), line, showTrees2(slice(ts, 1, len(ts)))])
}
out(layout(showTree2(exampleTree)) + '\n')

(
  g := (a, b, c) => group(cons([a, b, c]))
  example := g(
    g(
      g(
        g(
          text('hello'),
          line,
          text('a')
        ),
        line,
        text('b'),
      ),
      line,
      text('c'),
    ),
    line,
    text('d'),
  )
  each([13, 11, 9, 7, 5], (w, _) => out(pretty(w, example)+'\n'))
)

(
  # XML :: Elt (string, []Att, []XML) | Txt string
  elt := (tag, atts, xmls) => {kind: 'elt', n: tag, a: atts, c: xmls}
  txt := (s) => {kind: 'txt', s}
  # Att :: (string, string)
  quoted := (s) => '"' + s + '"'
  showAtts := (a) => ([n, v] := a, cons([text(n), text('='), text(quoted(v))]))
  showTag := (n, a) => cons([text(n), showFill(showAtts, a)])
  showFill := (f, xs) => len(xs) :: {
    0 -> nil
    _ -> bracket('', dbg(fill(dbg(concat(map(xs, f))))), '')
  }
  showXMLs := (x) => x.kind :: {
    'elt' -> len(x.c) :: {
      0 -> [cons([text('<'), showTag(x.n, x.a), text('/>')])]
      _ -> [cons([text('<'), showTag(x.n, x.a), text('>')
                  showFill(showXMLs, x.c),
                  text('</'), text(x.n), text('>')])]
    }
    'txt' -> map(words(x.s), text)
  }
  showXML := (x) => folddoc((x, y) => cons([x, y]), showXMLs(x))
  # example := showXML(elt('p', [['color', 'red'], ['font', 'Times'], ['size', '10']], [
  #   txt('Here is some'),
  #   elt('em', [], [txt('emphasized')]), txt('text.'),
  #   txt('Here is a'),
  #   elt('a', [['href', 'http://www.eg.com/']], [txt('link')]),
  #   txt('elsewhere.'),
  # ]))
  # TODO: fix
  # example := (folddoc((x, y) => cons([x, y]), showXMLs(elt('p', [['color', 'red'], ['size', '10']], []))))
  # example := showXML(elt('a', [['href', 'http://www.eg.com/']], [txt('link')]))
  # example := showXML(elt('a', [], []))
  # each([30, 60], (w, _) => out(pretty(w, example)+'\n'))
)
