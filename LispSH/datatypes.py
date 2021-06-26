from typing import List, Any, Union
from dataclasses import dataclass


class Symbol(str):
    def __hash__(self):
        return hash(("Symbol", str(self)))

    def __repr__(self):
        return f"<Symbol {str.__repr__(self)}>"

class Keyword(str):
    def __hash__(self):
        return hash(("Keyword", str(self)))

    def __repr__(self):
        return f":{self}"

@dataclass
class Atom:
    value: Union[bool, int, float]

    def __hash__(self):
        return hash(("Atom", self.value))

    def __repr__(self):
        return f"<Atom {repr(self.value)}>"

class Vector(list): pass

class Hashmap(dict):
    def __init__(self, keys_vals):
        assert len(keys_vals) % 2 == 0
        i = 0
        while i < len(keys_vals):
            key, val = keys_vals[i : i + 2]
            self[key] = val
            i += 2

    def __repr__(self):
        return '{' + ",".join(map(lambda kv: f"{repr(kv[0])}: {kv[1]}", self.items())) + '}'

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
