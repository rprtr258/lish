#!/usr/bin/env python
from typing import List, Any, Union
from dataclasses import dataclass
import math
import operator as op


################ Constants

OPEN_PAREN = '('
CLOSE_PAREN = ')'
QUOTE = '\''

################ Datatypes

@dataclass
class Symbol:
    name: str

@dataclass
class Macro:
    args: List[Symbol]
    body: List[Any]

@dataclass
class Atom:
    value: Union[bool, int, float]

def atom_or_symbol(token):
    if token in ["True", "False"]:
        return Atom(token == "True")
    if token in [True, False]:
        return Atom(token)
    try:
        return Atom(int(token))
    except ValueError:
        pass
    try:
        return float(int(token))
    except ValueError:
        pass
    return Symbol(token)

################ Parsing: parse, tokenize, and read_from_tokens

def tokenize(s: str) -> List[str]:
    "Convert a string into a list of tokens."
    return s \
        .replace(OPEN_PAREN, f" {OPEN_PAREN} ") \
        .replace(CLOSE_PAREN, f" {CLOSE_PAREN} ") \
        .replace(QUOTE, f" {QUOTE} ") \
        .split()

def read_from_tokens(tokens):
    "Read an expression from a sequence of tokens."
    if len(tokens) == 0:
        raise SyntaxError('unexpected EOF while reading')
    token = tokens.pop(0)
    if token == OPEN_PAREN:
        L = []
        if len(tokens) == 0:
            raise ValueError("Not enough close parens found")
        while tokens[0] != CLOSE_PAREN:
            L.append(read_from_tokens(tokens))
            if len(tokens) == 0:
                raise ValueError("Not enough close parens found")
        tokens.pop(0) # pop off ')'
        return L
    elif token == QUOTE:
        L = []
        L.append(read_from_tokens(tokens))
        return [Symbol("quote")] + L
    elif token == CLOSE_PAREN:
        raise SyntaxError(f'unexpected {CLOSE_PAREN}')
    else:
        return atom_or_symbol(token)

def parse(program):
    "Read a Scheme expression from a string."
    return read_from_tokens(tokenize(program))

################ Environments

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
            return None # nil
        return self.outer.find(var)
    def get(self, var):
        env = self.find(var)
        if env is None:
            return []
        return env[var]
    def __repr__(self):
        return "{" + " ".join(f"{k}: {v}" for k, v in self.items()) + "}"

class NamedFunction:
    def __init__(self, name, body):
        self.name = name
        self.body = body
    def __call__(self, *args, **kwargs):
        return self.body(*args, **kwargs)
    def __repr__(self):
        return f"<fun {self.name}>"

def plus(*x):
    if isinstance(x[0], Atom) and isinstance(x[0].value, int):
        return Atom(sum(map(lambda x: x.value, x), 0))
    if isinstance(x[0], str):
        return "".join(x) # sum(x, "")
    return sum(x, [])

def default_env():
    "An environment with some Scheme standard procedures."
    from random import random
    env = Env()
    env.update(vars(math)) # sin, cos, sqrt, pi, ...
    env.update({
        '+': NamedFunction("+", plus),
        '-': NamedFunction("-", lambda *x: Atom(x[0].value - sum(map(lambda x: x.value, x[1:])) if len(x) > 1 else -x[0].value)),
        '*': NamedFunction("*", lambda x, y: Atom(x.value * y.value)),
        '/':op.truediv,
        '>':op.gt,
        '<': lambda x, y: Atom(x < y), # TODO: change ops to reduce
        '>=':op.ge,
        '<=':op.le,
        '=':op.eq,
        "rand": NamedFunction("rand", random),
        'abs':     abs,
        "echo":    NamedFunction("echo", lambda *x: print(*map(schemestr, x))),
        "str":    NamedFunction("str", lambda *x: Atom(" ".join(map(schemestr, x)))),
        'append':  op.add,
        'apply':   lambda fx: fx[0](*fx[1:]),
        'progn':   NamedFunction("progn", lambda *x: x[-1]),
        'car':     lambda x: x[0],
        'cdr':     lambda x: x[1:],
        'cons':    lambda x,y: [x] + y,
        "=": lambda x, y: Atom(x == y),
        'length':  len,
        'list':    lambda *x: list(x),
        'list?':   lambda x: isinstance(x,list),
        'map':     lambda f, xs: list(map(f, xs)),
        'max':     max,
        'min':     min,
        'not':     op.not_,
        'nil?':    NamedFunction("nil?", lambda x: Atom(x == [])),
        'number?': lambda x: isinstance(x, Atom) and (isinstance(val := x.value, int) or isinstance(val, float)),
        'procedure?': callable,
        'round':   round,
        'symbol?': lambda x: isinstance(x, Symbol),
        "prompt": "lis.py> ",
    })
    return env

global_env = default_env()

################ Procedures

class Procedure(object):
    "A user-defined Scheme procedure."
    def __init__(self, args, body, env):
        self.args, self.body, self.env = args, body, env
    def __call__(self, *args): 
        return eval(self.body, Env(self.args, args, self.env))
    def __repr__(self):
        return f"<({self.args}) -> {self.body}> with {self.env.keys()}"

################ eval

def log_eval(eval):
    def wrapped_eval(x, env=global_env):
        from copy import deepcopy
        x0 = deepcopy(x)
        res = eval(x, env)
        print(x0, {k:v for k, v in env.items() if k == "prompt" or k not in global_env}, '=', res)
        return res
    return wrapped_eval

def macroexpand(args, body, exps, env):
    # print("ARGS:", schemestr(args))
    # print("BODY:", schemestr(body))
    # print("EXPS:", schemestr(exps))
    macroexpansion = eval(body, Env([arg.name for arg in args], exps, env))
    # print("EXPANDED:", schemestr(macroexpansion))
    # print()
    return macroexpansion

# @log_eval
def eval(x, env=global_env):
    "Evaluate an expression in an environment."
    if isinstance(x, Symbol):
        # x
        # but x is symbol
        return env.get(x.name)
    elif isinstance(x, Atom):
        # x
        # but x is atom (e.g. number)
        return x
    elif isinstance(x[0], Symbol) and (res := env.get(x[0].name)) and isinstance(res, Macro): # (macroname exps...)
        macroargs, macrobody = res.args, res.body
        exps = x[1:]
        # TODO: macroexpand
        macroexpansion = macroexpand(macroargs, macrobody, exps, env)
        return eval(macroexpansion, env)
    else:
        form_word = x[0]
        if form_word == Symbol("quote"):
            # (quote exp)
            _, exp = x
            return exp
        elif form_word == Symbol("atom"):
            # (atom exp)
            _, exp = x
            exp = eval(exp, env)
            return Atom(\
                isinstance(exp, Atom) or \
                isinstance(exp, Symbol) or \
                isinstance(exp, str) or \
                isinstance(exp, bool) or \
                exp == [])
        elif form_word == Symbol("cond"): # TODO: test return default in (cond p1 e1 p2 e2 default)
            # (cond p1 e1 p2 e2 ... pn en)
            # or
            # (cond p1 e1 p2 e2 ... pn en default)
            predicates_exps = x[1:]
            i = 0
            while i + 1 < len(predicates_exps):
                predicate, expression = predicates_exps[i : i + 2]
                i += 2
                if eval(predicate, env).value:
                    return eval(expression, env)
            # if default value is given
            if len(predicates_exps) % 2 == 1:
                return eval(predicates_exps[-1], env)
        elif form_word == Symbol("define"):         # (define var exp)
            _, var, exp = x
            assert isinstance(var, Symbol), "Definition name is not a symbol"
            env[var.name] = eval(exp, env)
            return [] # TODO: nil
        elif form_word == Symbol("defmacro"):         # (defmacro macroname (args...) body)
            _, macroname, args, body = x
            assert isinstance(macroname, Symbol), "Macro definition name is not a symbol"
            env[macroname.name] = Macro(args, body)
            return [] # TODO: nil
        elif form_word == Symbol("set!"):           # (set! var exp)
            _, var, exp = x
            assert isinstance(var, Symbol), "Definition name is not a symbol"
            var_name = var.name
            env.find(var_name)[var_name] = eval(exp, env)
        elif form_word == Symbol("lambda"):         # (lambda (args...) body)
            _, args, body = x
            return Procedure([arg.name for arg in args], body, env)
        else:
            # (proc arg...)
            proc = eval(x[0], env)
            args = [eval(exp, env) for exp in x[1:]]
            if not callable(proc):
                raise ValueError(f"{proc} is not a function call in {schemestr(x)}")
            return proc(*args)

def fix_parens(cmd_line):
    cmd_line = cmd_line.strip()
    if cmd_line[0] not in [OPEN_PAREN, QUOTE]:
        cmd_line = OPEN_PAREN + cmd_line
    if cmd_line[-1] != CLOSE_PAREN:
        cmd_line = cmd_line + CLOSE_PAREN
    # TODO: don't count brackets in strings
    open_parens, close_parens = cmd_line.count(OPEN_PAREN), cmd_line.count(CLOSE_PAREN)
    return \
        OPEN_PAREN * max(0, close_parens - open_parens) + \
        cmd_line + \
        CLOSE_PAREN * max(0, open_parens - close_parens)

################ Interaction: A REPL

def schemestr(exp):
    "Convert a Python object back into a Scheme-readable string."
    if isinstance(exp, Symbol):
        return exp.name
    elif isinstance(exp, Atom):
        return str(exp.value)
    elif isinstance(exp, list):
        return '(' + ' '.join(map(schemestr, exp)) + ')'
    else:
        return str(exp)

def repl():
    "A prompt-read-eval-print loop."
    while True:
        input_line = input(eval(symbol("prompt")))
        input_line = fix_parens(input_line)
        val = eval(parse(input_line))
        if val != []: # TODO: nil
            print(schemestr(val))

################ File load

def load_file(filename):
    with open(filename, "r") as fd:
        for line in fd:
            eval(parse(line))
