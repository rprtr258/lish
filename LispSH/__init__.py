#!/usr/bin/env python
from typing import List, Any
import math
import operator as op

################ Types

List   = list         # A Lisp List is implemented as a Python list

################ Constants

OPEN_PAREN = '('
CLOSE_PAREN = ')'
QUOTE = '\''

################ Datatypes

def macro(args, body): return ["MACRO", args, body]
def is_macro(form): return isinstance(form, list) and len(form) == 3 and form[0] == "MACRO"

# Symbol is implemented as a ["SYMBOL", symbol_name]
def symbol(name): return ["SYMBOL", name]
def is_symbol(form): return isinstance(form, list) and len(form) == 2 and form[0] == "SYMBOL"
def symbol_name(symbol): return symbol[1]

# Number, Boolean is implemented as a ["ATOM", value]
def atom(token: str) -> List[Any]:
    "Numbers become numbers; every other token is a symbol."
    try:
        return ["ATOM", int(token)]
    except ValueError:
        try:
            return ["ATOM", float(token)]
        except ValueError:
            return symbol(token) # TODO: wtf atom returns symbol?
def is_atom(form): return isinstance(form, list) and len(form) == 2 and form[0] == "ATOM"
def atom_value(atom): return atom[1]

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
        return [symbol("quote")] + L
    elif token == CLOSE_PAREN:
        raise SyntaxError(f'unexpected {CLOSE_PAREN}')
    else:
        return atom(token)

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

def default_env():
    "An environment with some Scheme standard procedures."
    env = Env()
    env.update(vars(math)) # sin, cos, sqrt, pi, ...
    env.update({
        '+': NamedFunction("+", lambda *x: sum(x) if isinstance(x[0], int) else "".join(x)),
        '-': NamedFunction("-", lambda *x: x[0] - sum(x[1:]) if len(x) > 1 else -x[0]),
        '*':op.mul,
        '/':op.truediv,
        '>':op.gt,
        '<':op.lt,
        '>=':op.ge,
        '<=':op.le,
        '=':op.eq,
        'abs':     abs,
        "echo":    NamedFunction("echo", lambda *x: print(*x)),
        'append':  op.add,
        'apply':   lambda fx: fx[0](*fx[1:]),
        'progn':   NamedFunction("progn", lambda *x: x[-1]),
        'car':     lambda x: x[0],
        'cdr':     lambda x: x[1:],
        'cons':    lambda x,y: [x] + y,
        'eq?':     lambda x, y: x == y,
        'equal?':  op.eq,
        'length':  len,
        'list':    lambda *x: list(x),
        'list?':   lambda x: isinstance(x,list),
        'map':     map,
        'max':     max,
        'min':     min,
        'not':     op.not_,
        'nil?':    NamedFunction("nil?", lambda x: x == []),
        'number?': lambda x: is_atom(x) and (isinstance(val := atom_value(x), int) or isinstance(val, float)),
        'procedure?': callable,
        'round':   round,
        'symbol?': is_symbol,
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

# @log_eval
def eval(x, env=global_env):
    "Evaluate an expression in an environment."
    if is_symbol(x):      # variable reference
        return env.get(symbol_name(x))
        # return eval(env.get(symbol_name(x)), env)
    elif is_atom(x):      # atom aka constant
        return atom_value(x)
    elif not isinstance(x, List):  # constant literal
        return x
    elif x[0] == symbol("quote"):          # (quote exp)
        _, exp = x
        return exp
    elif x[0] == symbol("atom"):          # (atom exp)
        _, exp = x
        exp = eval(exp, env)
        return \
            is_atom(exp) or \
            is_symbol(exp) or \
            isinstance(exp, str) or \
            isinstance(exp, bool) or \
            exp == []
    elif x[0] == symbol("cond"): # TODO: test return default in (cond p1 e1 p2 e2 default)
        predicates_exps = x[1:]
        i = 0
        while i + 1 < len(predicates_exps):
            predicate, expression = predicates_exps[i : i + 2]
            i += 2
            if eval(eval(predicate, env), env):
                return eval(expression, env)
        # if default value is given
        if len(predicates_exps) % 2 == 1:
            return eval(predicates_exps[-1], env)
    elif x[0] == symbol("define"):         # (define var exp)
        _, var, exp = x
        assert is_symbol(var), "Definition name is not a symbol"
        env[symbol_name(var)] = eval(exp, env)
        return [] # TODO: nil
    elif x[0] == symbol("defmacro"):         # (defmacro macroname (args...) body)
        _, macroname, args, body = x
        assert is_symbol(macroname), "Macro definition name is not a symbol"
        env[symbol_name(macroname)] = macro(args, body)
        return [] # TODO: nil
    elif x[0] == symbol("set!"):           # (set! var exp)
        _, var, exp = x
        assert is_symbol(var), "Definition name is not a symbol"
        var_name = symbol_name(var)
        env.find(var_name)[var_name] = eval(exp, env)
    elif x[0] == symbol("lambda"):         # (lambda (args...) body)
        _, args, body = x
        return Procedure([symbol_name(arg) for arg in args], body, env)
    elif is_symbol(x[0]) and (res := env.get(symbol_name(x[0]))) and is_macro(res): # (macroname exps...)
        _, macroargs, macrobody = res
        exps = x[1:]
        # TODO: macroexpand
        macroexpansion = eval(macrobody, Env([symbol_name(arg) for arg in macroargs], exps, env))
        return eval(macroexpansion, env)
    else:                          # (proc arg...)
        proc = eval(x[0], env)
        args = [eval(exp, env) for exp in x[1:]]
        if not callable(proc):
            raise ValueError(f"{proc} is not a function call in {x}")
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
    if is_symbol(exp):
        return symbol_name(exp)
    elif is_atom(exp):
        return atom_value(exp)
    elif isinstance(exp, List):
        return '(' + ' '.join(map(schemestr, exp)) + ')'
    else:
        return str(exp)

def repl():
    "A prompt-read-eval-print loop."
    while True:
        input_line = input(eval(symbol("prompt")))
        input_line = fix_parens(input_line)
        val = eval(parse(input_line))
        if not val is None: 
            print(schemestr(val))

################ File load

def load_file(filename):
    with open(filename, "r") as fd:
        for line in fd:
            eval(parse(line))
