# september syntax highlighter command

{log, map, each, slice, cat} := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
{Gray, Red, Green, Yellow, Blue, Magenta, Cyan} := import('../vendor/ansi.ink')
Norm := s => s

{Tok, tokenizeWithComments: tokenize} := import('tokenize.ink')

# associating token types with their highlight colors
colorFn := tok => tok.type :: {
  Tok.Separator -> Norm

  Tok.Comment -> Gray

  Tok.Ident -> Norm
  Tok.EmptyIdent -> Norm

  Tok.NumberLiteral -> Magenta
  Tok.StringLiteral -> Yellow
  Tok.TrueLiteral -> Magenta
  Tok.FalseLiteral -> Magenta

  Tok.AccessorOp -> Red
  Tok.EqOp -> Red

  Tok.FunctionArrow -> Green

  # operators are all red
  Tok.KeyValueSeparator -> Red
  Tok.DefineOp -> Red
  Tok.MatchColon -> Red
  Tok.CaseArrow -> Red
  Tok.SubOp -> Red
  Tok.NegOp -> Red
  Tok.AddOp -> Red
  Tok.MulOp -> Red
  Tok.DivOp -> Red
  Tok.ModOp -> Red
  Tok.GtOp -> Red
  Tok.LtOp -> Red
  Tok.AndOp -> Red
  Tok.OrOp -> Red
  Tok.XorOp -> Red

  Tok.LParen -> Cyan
  Tok.RParen -> Cyan
  Tok.LBracket -> Cyan
  Tok.RBracket -> Cyan
  Tok.LBrace -> Cyan
  Tok.RBrace -> Cyan

  _ -> () # should error, unreachable
}

prog => (
  tokens := tokenize(prog)
  spans := map(tokens, (tok, i) => {
    colorFn: [tok.type, tokens.(i + 1)] :: {
      # direct function calls are marked green
      # on a best-effort basis
      [
        Tok.Ident
        {type: Tok.LParen, val: _, line: _, col: _, i: _}
      ] -> Green
      _ -> colorFn(tok)
    }
    start: tok.i
    end: tokens.(i + 1) :: {
      () -> len(prog)
      _ -> tokens.(i + 1).i
    }
  })
  pcs := map(
    spans
    span => (span.colorFn)(slice(prog, span.start, span.end))
  )
  cat(pcs, '')
)
