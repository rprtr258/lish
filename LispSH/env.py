import math
import operator as op

from LispSH.datatypes import is_atom, Symbol, Keyword
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
                # import pdb;pdb.set_trace()
                print(f"WARNING: {repr(var)} was not found")
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
        return f"<func {self.name}>"

def plus(*x):
    if is_atom(x[0]):
        if isinstance(x[0], int) or isinstance(x[0], float):
            return sum(x, 0)
        elif isinstance(x[0], str):
            return "".join(x) # sum(x, "")
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
        Symbol(k): v
        for k, v in {
            # ARIPHMETIC OPERATORS
            "+": NamedFunction("+", plus),
            "-": NamedFunction("-", lambda *x: reduce(lambda a, b: a - b, x[1:], x[0]) if len(x) > 1 else -x[0]),
            "*": NamedFunction("*", lambda *x: reduce(lambda a, b: a * b, x[1:], x[0])),
            "/": NamedFunction("/", lambda *x: reduce(lambda a, b: a // b, x[1:], x[0])),

            # MATH FUNCTIONS
            "rand": NamedFunction("rand", random),
            "abs": NamedFunction("abs", abs),
            "cos": NamedFunction("cos", cos),
            "max": NamedFunction("max", lambda *x: reduce(lambda a, b: max(a, b), x[1:], x[0])),
            "min": NamedFunction("min", lambda *x: reduce(lambda a, b: min(a, b), x[1:], x[0])),
            "round": NamedFunction("round", round),

            # COMPARISON OPERATORS
            ">": NamedFunction(">", lambda *x: all(map(lambda xy: op.gt(xy[0], xy[1]), zip(x, x[1:])))),
            "<": NamedFunction("<", lambda *x: all(map(lambda xy: op.lt(xy[0], xy[1]), zip(x, x[1:])))),
            ">=": NamedFunction(">=", lambda *x: all(map(lambda xy: op.ge(xy[0], xy[1]), zip(x, x[1:])))),
            "<=": NamedFunction("<=", lambda *x: all(map(lambda xy: op.le(xy[0], xy[1]), zip(x, x[1:])))),
            "=": NamedFunction("=", lambda *x: all(map(lambda xy: op.eq(xy[0], xy[1]), zip(x, x[1:])))),
            "nil?": NamedFunction("nil?", lambda x: x == []),
            "number?": NamedFunction("number?", lambda x: isinstance(x, int) or isinstance(x, float)),
            "procedure?": NamedFunction("procedure?", callable),
            "atom?": NamedFunction("atom?", lambda x: is_atom(x) or isinstance(x, Symbol) or isinstance(x, Keyword) or x == []),
            "symbol?": lambda x: isinstance(x, Symbol),
            "list?": NamedFunction("list?", lambda x: isinstance(x, list)),

            # BOOL FUNCTIONS
            "or": NamedFunction("or", lambda *x: reduce(lambda x, y: x or y, x)),
            "not": NamedFunction("not", lambda x: op.not_(x)),

            # LIST OPERATIONS
            "cons": NamedFunction("cons", cons),
            "map": NamedFunction("map", lambda f, x: list(map(f, x))),
            # TODO: f is identity by default
            "sorted-by": NamedFunction("sorted-by", lambda x, f: sorted(x, key=lambda x: f(x))),
            "len": NamedFunction("len", len),
            "car": NamedFunction("car", lambda x: x[0]),
            "cdr": NamedFunction("cdr", lambda x: x[1:]),
            "list": NamedFunction("list", lambda *x: list(x)),

            # STRING FUNCTIONS
            "join": NamedFunction("join", lambda d, x: d.join(x)),
            "str": NamedFunction("str", lambda *x: " ".join(map(PRINT, x))),

            # FILE OPERATIONS
            # TODO: tests
            "path-getsize": NamedFunction("path-getsize", path.getsize),

            # OTHER FUNCTIONS
            "ls-r": NamedFunction("ls-r", lambda x: [[dir_name, files] for dir_name, _, files in walk(x)]),
            "echo": NamedFunction("echo", echo),
            "name": NamedFunction("name", lambda x: x if isinstance(x, Symbol) else []),
            'apply':   lambda fx: fx[0](*fx[1:]),
            'progn':   NamedFunction("progn", lambda *x: x[-1]),
            # TODO: rename to parse-int? / str->int
            "int": NamedFunction("int", int),
            "exit": lambda: exit(0),
            "prompt": lambda: "lis.py> ",
    }.items()})
    return env

global_env = default_env()
