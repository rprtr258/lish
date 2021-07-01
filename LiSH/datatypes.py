from typing import Any
# TODO: own list type with str = PRINT


class Symbol(str):
    @property
    def name(self):
        return str(self)

    def __hash__(self):
        return hash(("Symbol", str(self)))

    def __repr__(self):
        return f"<Symbol {str.__repr__(self)}>"


class Keyword(str):
    def __hash__(self):
        return hash(("Keyword", str(self)))

    def __repr__(self):
        return f":{self}"


def is_atom(x: Any):
    return isinstance(x, int) or isinstance(x, float) or isinstance(x, str)


class Hashmap(dict):
    def __init__(self, keys_vals):
        assert len(keys_vals) % 2 == 0
        i = 0
        while i < len(keys_vals):
            key, val = keys_vals[i: i + 2]
            self[key] = val
            i += 2

    def __repr__(self):
        return '{' + ",".join(map(lambda kv: f"{repr(kv[0])}: {kv[1]}", self.items())) + '}'
