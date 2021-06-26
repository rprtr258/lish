from typing import List, Any, Union
from dataclasses import dataclass


class Symbol(str): pass
class Keyword(str): pass

@dataclass
class Atom(object):
    value: Union[bool, int, float]

@dataclass
class Macro:
    args: List[Symbol]
    body: List[Any]

def get_atom_value(atom): return atom.value

def atom_or_symbol(token):
    if token[0] == '"' and token[-1] == '"' and len(token) >= 2:
        return Atom(token[1 : -1])
    if token in ["true", "false"]:
        return Atom(token == "true")
    if token in [True, False]:
        return Atom(token)
    try:
        return Atom(int(token))
    except ValueError:
        pass
    try:
        return Atom(float(token))
    except ValueError:
        pass
    return Symbol(token)
