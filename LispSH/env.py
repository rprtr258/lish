import math
import operator as op

from LispSH.datatypes import get_atom_value, Atom, Symbol, Keyword
from LispSH.printer import PRINT, pr_str_no_escape


class Env(dict):
    "An environment: a dict of {'var':val} pairs, with an outer Env."
    def __init__(self, parms=(), args=(), outer=None):
        self.update(zip(parms, args))
        self.outer = outer
    def find(self, var):
        "Find the innermost Env where var appears."
        if var in self:
            return self
        if self.outer is None:
            if var not in ["cond", "quote", "atom?", "lambda", "define", "defmacro", "macroexpand", "set!"]:
                print(f"WARNING: {var} was not found")
            return None # nil
        return self.outer.find(var)
    def get(self, var):
        env = self.find(var)
        if env is None:
            return []
        return env[var]
    def __repr__(self):
        return "{\n  " + "  \n".join(f"{k}: {v}" for k, v in self.items()) + "\n}"

class NamedFunction:
    def __init__(self, name, body):
        self.name = name
        self.body = body
    def __call__(self, *args, **kwargs):
        return self.body(*args, **kwargs)
    def __repr__(self):
        return self.name

def plus(*x):
    if isinstance(x[0], Atom):
        if isinstance(x[0].value, int) or isinstance(x[0].value, float):
            return Atom(sum(map(get_atom_value, x), 0))
        elif isinstance(x[0].value, str):
            return Atom("".join(map(get_atom_value, x))) # sum(x, "")
    return sum(x, [])

def echo(*x):
    for y in x:
        print(pr_str_no_escape(y), end="")
    print()
    return []

def cons(*x):
    x = list(x)
    return x[:-1] + x[-1]

def default_env():
    "An environment with some Scheme standard procedures."
    from random import random
    from math import cos
    from functools import reduce
    from os import walk, path
    env = Env()
    # env.update(vars(math)) # sin, cos, sqrt, pi, ...
    env.update({
        # ARIPHMETIC OPERATORS
        "+": NamedFunction("+", plus),
        "-": NamedFunction("-", lambda *x: reduce(lambda a, b: Atom(a.value - b.value), x[1:], x[0]) if len(x) > 1 else Atom(-x[0].value)),
        "*": NamedFunction("*", lambda *x: reduce(lambda a, b: Atom(a.value * b.value), x[1:], x[0])),
        "/": NamedFunction("/", lambda *x: reduce(lambda a, b: Atom(a.value // b.value), x[1:], x[0])),
        # MATH FUNCTIONS
        "rand": NamedFunction("rand", lambda:Atom(random())),
        "abs": NamedFunction("abs", lambda x: Atom(abs(x.value))),
        "cos": NamedFunction("cos", lambda x: Atom(cos(x.value))),
        "max": NamedFunction("max", lambda *x: reduce(lambda a, b: Atom(max(a.value, b.value)), x[1:], x[0])),
        "min": NamedFunction("min", lambda *x: reduce(lambda a, b: Atom(min(a.value, b.value)), x[1:], x[0])),
        "round": NamedFunction("round", lambda x: Atom(round(x.value))),
        # COMPARISON OPERATORS
        ">": NamedFunction(">", lambda *x: Atom(all(map(lambda xy: op.gt(xy[0].value, xy[1].value), zip(x, x[1:]))))),
        "<": NamedFunction("<", lambda *x: Atom(all(map(lambda xy: op.lt(xy[0].value, xy[1].value), zip(x, x[1:]))))),
        ">=": NamedFunction(">=", lambda *x: Atom(all(map(lambda xy: op.ge(xy[0].value, xy[1].value), zip(x, x[1:]))))),
        "<=": NamedFunction("<=", lambda *x: Atom(all(map(lambda xy: op.le(xy[0].value, xy[1].value), zip(x, x[1:]))))),
        "=": NamedFunction("=", lambda *x: Atom(all(map(lambda xy: op.eq(xy[0], xy[1]), zip(x, x[1:]))))),
        "nil?": NamedFunction("nil?", lambda x: Atom(x == [])),
        "number?": NamedFunction("number?", lambda x: Atom(isinstance(x, Atom) and (isinstance(val := x.value, int) or isinstance(val, float)))),
        "procedure?": NamedFunction("procedure?", lambda x: Atom(callable(x))),
        "atom?": NamedFunction("atom?", lambda x: Atom(isinstance(x, Atom) or isinstance(x, Symbol) or isinstance(x, Keyword) or x == [])),
        "symbol?": lambda x: Atom(isinstance(x, Symbol)),
        "list?": NamedFunction("list?", lambda x: Atom(isinstance(x, list))),
        # BOOL FUNCTIONS
        "or": NamedFunction("or", lambda *x: reduce((lambda x, y: Atom(x.value or y.value)), x)),
        "not": NamedFunction("not", lambda x: Atom(op.not_(x.value))),
        # LIST OPERATIONS
        "cons": NamedFunction("cons", cons),
        "map": NamedFunction("map", lambda f, x: list(map(f, x))),
        # TODO: f is identity by default
        "sorted-by": NamedFunction("sorted-by", lambda x, f: sorted(x, key=lambda x: f(x).value)),
        "len": NamedFunction("len", lambda x: Atom(len(x if isinstance(x, list) else x.value))),
        "car": NamedFunction("car", lambda x: x[0]),
        "cdr": NamedFunction("cdr", lambda x: x[1:]),
        "list": NamedFunction("list", lambda *x: list(x)),
        # STRING FUNCTIONS
        "join": NamedFunction("join", lambda d, x: Atom(d.value.join(map(get_atom_value, x)))),
        "str": NamedFunction("str", lambda *x: Atom(" ".join(map(PRINT, x)))),
        # FILE OPERATIONS
        # TODO: tests
        "path-getsize": NamedFunction("path-getsize", lambda x: Atom(path.getsize(x.value))),
        # OTHER FUNCTIONS
        "ls-r": NamedFunction("ls-r", lambda x: [[Atom(dir_name), list(map(Atom, files))] for dir_name, _, files in walk(x.value)]),
        "echo":    NamedFunction("echo", echo),
        "name": NamedFunction("name", lambda x: Atom(x) if isinstance(x, Symbol) else []),
        'apply':   lambda fx: fx[0](*fx[1:]),
        'progn':   NamedFunction("progn", lambda *x: x[-1]),
        # TODO: rename to parse-int? / str->int
        "int": NamedFunction("int", lambda x: Atom(int(x.value))),
        "exit": lambda: exit(0),
        "prompt": lambda: Atom("lis.py> "),
    })
    return env

global_env = default_env()
