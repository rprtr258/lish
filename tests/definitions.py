from context import LispSH
from LispSH import symbol

NIL = []
A = symbol("a")
B = symbol("b")
C = symbol("c")
ATOM_SYMBOL = symbol("atom")
QUOTE_SYMBOL = symbol("quote")
COND_SYMBOL = symbol("cond")
EQ_SYMBOL = symbol("eq?")
QA = [QUOTE_SYMBOL, A]
QB = [QUOTE_SYMBOL, B]
QC = [QUOTE_SYMBOL, C]
QNIL = [QUOTE_SYMBOL, NIL]
