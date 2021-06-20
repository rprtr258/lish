from context import LispSH
from LispSH import Symbol, Atom

NIL = []
TRUE = Atom(True)
FALSE = Atom(False)
A = Symbol("a")
B = Symbol("b")
C = Symbol("c")
ATOM_SYMBOL = Symbol("atom")
QUOTE_SYMBOL = Symbol("quote")
COND_SYMBOL = Symbol("cond")
EQ_SYMBOL = Symbol("eq?")
QA = [QUOTE_SYMBOL, A]
QB = [QUOTE_SYMBOL, B]
QC = [QUOTE_SYMBOL, C]
QNIL = [QUOTE_SYMBOL, NIL]
