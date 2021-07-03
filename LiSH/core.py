import operator as op
from random import random
from math import cos
from functools import reduce
from os import walk, path

from LiSH.datatypes import is_atom, Keyword, Symbol
from LiSH.reader import read_str
from LiSH.error_reporting.errors import FunctionCallError
from LiSH.printer import PRINT, pr_str


def plus(*x):
    if is_atom(x[0]):
        if isinstance(x[0], int) or isinstance(x[0], float):
            return sum(x, 0)
        elif isinstance(x[0], str):
            return "".join(x)  # sum(x, "")
    return sum(x, [])


def echo(*x):
    for y in x:
        print(pr_str(y, None), end="")
    print()
    return []


def cons(*x):
    x = list(x)
    return x[:-1] + x[-1]


def slurp(filename):
    with open(filename, "r") as fd:
        return fd.read()


def get(coll, key):
    if isinstance(coll, list):
        assert isinstance(key, int), f"Index {key} is not int"
        assert 0 <= key and key < len(coll), f"{key} is out of bounds"
        return coll[key]
    elif isinstance(coll, list):
        assert key in coll, f"{key} is not in hashmap"
        return coll[key]


def throw(message):
    raise RuntimeError(message)  # TODO: own exception type


def apply(proc, args):
    # TODO: fix
    # if len(proc.args) != len(args):
    #     raise RuntimeError(f"{proc} expected {len(proc.args)} arguments, but got {len(args)}")
    try:
        return proc(*args)
    except Exception as e:
        raise FunctionCallError(proc, args, e)


def sort_by(xs, key=lambda x: x):
    return sorted(xs, key=key)


def car(xs):
    assert len(xs) > 0, "Tried to car empty array"
    return xs[0]


def hash_map(*args):
    if len(args) % 2 != 0:
        raise RuntimeError(f"Wrong number of args to hash-map{len(args)}")
    return {
        args[i]: args[i + 1]
        for i in range(0, len(args), 2)
    }


ns = {
    # ARIPHMETIC OPERATORS
    "+": plus,
    "-": lambda *x: reduce(lambda a, b: a - b, x[1:], x[0]) if len(x) > 1 else -x[0],
    "*": lambda *x: reduce(lambda a, b: a * b, x[1:], x[0]),
    "/": lambda *x: reduce(lambda a, b: a // b, x[1:], x[0]),

    # MATH FUNCTIONS
    "rand": random,
    "abs": abs,
    "cos": cos,
    "max": lambda *x: reduce(lambda a, b: max(a, b), x[1:], x[0]),
    "min": lambda *x: reduce(lambda a, b: min(a, b), x[1:], x[0]),
    "round": round,

    # COMPARISON OPERATORS
    ">": lambda *x: all(map(lambda xy: op.gt(xy[0], xy[1]), zip(x, x[1:]))),
    "<": lambda *x: all(map(lambda xy: op.lt(xy[0], xy[1]), zip(x, x[1:]))),
    ">=": lambda *x: all(map(lambda xy: op.ge(xy[0], xy[1]), zip(x, x[1:]))),
    "<=": lambda *x: all(map(lambda xy: op.le(xy[0], xy[1]), zip(x, x[1:]))),
    "=": lambda *x: all(map(lambda xy: op.eq(xy[0], xy[1]), zip(x, x[1:]))),
    "number?": lambda x: isinstance(x, int) or isinstance(x, float),
    "fun?": callable,
    "atom?": lambda x: is_atom(x) or isinstance(x, Symbol) or isinstance(x, Keyword) or x == [],
    "symbol?": lambda x: isinstance(x, Symbol),

    # BOOL FUNCTIONS
    "or": lambda *x: reduce(lambda x, y: x or y, x),
    "and": lambda *x: reduce(lambda x, y: x and y, x),

    # LIST OPERATIONS
    "list?": lambda x: isinstance(x, list),
    "nil?": lambda x: x == [],
    "cons": cons,
    "sort": sort_by,
    "len": len,
    "car": car,
    "cdr": lambda x: x[1:],
    "list": lambda *x: list(x),

    # HASHMAP OPERATIONS
    "hash-map": hash_map,

    # STRING FUNCTIONS
    "join": lambda d, x: d.join(x),
    "str": lambda *x: " ".join(map(PRINT, x)),

    # FILE OPERATIONS
    # TODO: tests
    "path-getsize": path.getsize,
    "slurp": slurp,

    # OTHER FUNCTIONS
    "apply": apply,
    "throw": throw,
    "get": get,
    "ls-r": lambda x: [[dir_name, files] for dir_name, _, files in walk(x)],
    "echo": echo,
    "name": lambda x: x if isinstance(x, Symbol) else [],
    "progn": lambda *x: x[-1],
    "int": int,
    "exit": lambda: exit(0),
    "prompt": lambda: "lis.py> ",
    "read": read_str,
}
