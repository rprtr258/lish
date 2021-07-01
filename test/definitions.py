from context import LiSH
from LiSH.datatypes import Symbol

NIL = []
# TODO: inline
TRUE = True
FALSE = False
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
