#!/usr/bin/env python

# TODO: variadic defun (defn (& x) x) and macro (defmacro (& x) x)
# TODO: (loop ... recur) or Tail call optimisation(harder, non-recursive eval?)

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

################ Parsing: parse, tokenize, and read_from_tokens

def no_quote_replace(s: str, c: str, p: str):
    "Replaces c(char) to p(pattern) in s if c not in double quotes"
    i = 0
    res = ""
    quoted = False
    while i < len(s):
        if s[i] == '"':
            res += '"'
            quoted = not quoted
        elif s[i] == c and not quoted:
            res += p
        else:
            res += s[i]
        i += 1
    return res

def remove_comment(s: str) -> str:
    if ';' in s:
        return s[:s.find(';')] # TODO: check "wOwOw ;;;; " ; fuck you
    return s

# TODO: test mirroring
def tokenize(s: str) -> List[str]:
    "Convert a string into a list of tokens."
    s = remove_comment(s)
    word = None
    res = []
    quoted = False
    s = no_quote_replace(s, OPEN_PAREN, f" {OPEN_PAREN} ")
    s = no_quote_replace(s, CLOSE_PAREN, f" {CLOSE_PAREN} ")
    s = no_quote_replace(s, QUOTE, f" {QUOTE} ")
    mirrored = False
    for c in s:
        if c in [' ', '\n']:
            if not quoted:
                if not word is None:
                    res.append(word)
                    word = None
                mirrored = False
            else:
                word += c
        elif c == '\\' and quoted:
            if mirrored:
                word += '\\'
                mirrored = False
            else:
                mirrored = True
        else:
            word = c if word is None else word + c
            if c == '"' and not mirrored:
                if quoted:
                    res.append(word)
                    word = None
                quoted = not quoted
            mirrored = False
    if not word is None:
        res.append(word)
    return res

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
        return f"<fun {self.name}>"

def plus(*x):
    if isinstance(x[0], Atom):
        if isinstance(x[0].value, int) or isinstance(x[0].value, float):
            return Atom(sum(map(get_atom_value, x), 0))
        elif isinstance(x[0].value, str):
            return Atom("".join(map(get_atom_value, x))) # sum(x, "")
    return sum(x, [])

def echo(*x):
    def echo_helper(*x):
        if len(x) == 1:
            x = x[0]
            if isinstance(x, Atom):
                return str(x.value)
            else:
                return schemestr(x)
        else:
            res = ""
            for xi in x:
                if isinstance(xi, Atom):
                    res += str(xi.value)
                elif isinstance(xi, list):
                    res += "(" + echo_helper(*xi) + ")"
                else:
                    res += str(xi)
            return res
    print(echo_helper(*x))
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
        "list?": NamedFunction("list?", lambda x: Atom(isinstance(x, list))),
        "number?": NamedFunction("number?", lambda x: Atom(isinstance(x, Atom) and (isinstance(val := x.value, int) or isinstance(val, float)))),
        "procedure?": NamedFunction("procedure?", lambda x: Atom(callable(x))),
        "symbol?": lambda x: Atom(isinstance(x, Symbol)),
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
        "str": NamedFunction("str", lambda *x: Atom(" ".join(map(schemestr, x)))),
        # FILE OPERATIONS
        # TODO: tests
        "path-getsize": NamedFunction("path-getsize", lambda x: Atom(path.getsize(x.value))),
        # OTHER FUNCTIONS
        "ls-r": NamedFunction("ls-r", lambda x: [[Atom(dir_name), list(map(Atom, files))] for dir_name, _, files in walk(x.value)]),
        "echo":    NamedFunction("echo", echo),
        "name": NamedFunction("name", lambda x: Atom(x.name)),
        'apply':   lambda fx: fx[0](*fx[1:]),
        'progn':   NamedFunction("progn", lambda *x: x[-1]),
        # TODO: rename to parse-int? / str->int
        "int": NamedFunction("int", lambda x: Atom(int(x.value))),
        "prompt": lambda: Atom("lis.py> "),
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
        return schemestr([Symbol("lambda"), self.args, self.body])

################ eval

def log_eval(eval):
    def wrapped_eval(x, env=global_env):
        from copy import deepcopy
        x0 = deepcopy(x)
        res = eval(x, env)
        print(x0, {k:v for k, v in env.items() if k == "prompt" or k not in global_env}, '=', res)
        return res
    return wrapped_eval

def macroexpand(macroform, env):
    macroname, *exps = macroform
    res = env.get(macroname.name)
    macroargs, macrobody = res.args, res.body
    return eval(macrobody, Env([arg.name for arg in macroargs], exps, env))

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
    elif isinstance(x[0], Symbol) and (res := env.get(x[0].name)) and isinstance(res, Macro):
        # (macroname exps...)
        macroexpansion = macroexpand(x, env)
        return eval(macroexpansion, env)
    else:
        form_word = x[0]
        if form_word == Symbol("quote"):
            # (quote exp)
            _, exp = x
            return exp
        elif form_word == Symbol("atom?"):
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
            assert isinstance(var, Symbol), f"""Definition name is not a symbol, but a {schemestr(var)}"""
            env[var.name] = eval(exp, env)
            return env[var.name]
        elif form_word == Symbol("macroexpand"):
            # (macroexpand (macro exps...))
            _, macroform = x
            return macroexpand(macroform, env)
        elif form_word == Symbol("defmacro"):         # (defmacro macroname (args...) body)
            _, macroname, args, body = x
            assert isinstance(macroname, Symbol), "Macro definition name is not a symbol"
            env[macroname.name] = Macro(args, body)
            return [] # TODO: nil
        elif form_word == Symbol("set!"):
            # (set! var exp)
            _, var, exp = x
            assert isinstance(var, Symbol), "Definition name is not a symbol"
            var_name = var.name
            new_var_value = eval(exp, env)
            env.find(var_name)[var_name] = new_var_value
            return new_var_value
        elif form_word == Symbol("lambda"):
            # (lambda (args...) body)
            _, args, body = x
            for arg in args:
                assert isinstance(arg, Symbol), f"Argument name is not a symbol, but a {schemestr(arg)}"
            return Procedure([arg.name for arg in args], body, env)
        elif form_word == Symbol("apply"):
            # (apply f (args...))
            _, proc, args = x
            proc = eval(proc, env)
            args = eval(args, env)
            try:
                return proc(*args)
            except Exception as e:
                print(RuntimeError(f"""Error during evaluation ({proc} {" ".join(map(schemestr, args))}).
Error is:
    {"Recursed" if isinstance(e, RuntimeError) else e}"""))
        else:
            # (proc arg...)
            proc = eval(form_word, env)
            args = [eval(exp, env) for exp in x[1:]]
            if not callable(proc):
                print(RuntimeError(f"""{proc} (named {schemestr(x[0])}) is not a function call in {schemestr(x)}."""))
            try:
                if (res := proc(*args)) is None:
                    print("FUCK YOU,", schemestr(x), schemestr(args))
                return res
            except Exception as e:
                print(RuntimeError(f"""Error during evaluation ({proc} {" ".join(map(schemestr, args))}).
Error is:
    {"Recursed" if isinstance(e, RuntimeError) else e}"""))

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
    if exp == []:
        return "nil"
    elif isinstance(exp, Symbol):
        return exp.name
    elif isinstance(exp, Atom):
        if isinstance(exp.value, str):
            return repr(exp.value).replace("'", '"')
        else:
            return str(exp.value)
    elif isinstance(exp, list) and exp[0] == Symbol("quote") and len(exp) == 1:
        return "(quote)"
    elif isinstance(exp, list) and exp[0] == Symbol("quote"):
        assert len(exp) == 2, f"Quote has zero or more than one argument: {exp}"
        return "'" + schemestr(exp[1])
    elif isinstance(exp, list):
        return '(' + ' '.join(map(schemestr, exp)) + ')'
    elif exp is None:
        print("[FEAR AND LOATHING IN NONE VEGAS]")
    else:
        print("WTF IS THIS:", exp)
        return str(exp)

# TODO: add Ctrl-D support
# TODO: Shift-Enter for multiline input
def repl():
    "A prompt-read-eval-print loop."
    from sys import stdin, stdout
    print(eval([Symbol("prompt")]).value, end="")
    stdout.flush()
    for line in stdin:
        prompt = eval([Symbol("prompt")]).value
        if line.strip() == "":
            print(prompt, end="")
            stdout.flush()
            continue
        line = fix_parens(line)
        val = eval(parse(line))
        if val != []: # TODO: nil
            print(schemestr(val))
            print(prompt, end="")
            stdout.flush()

################ File load

def load_file(filename):
    with open(filename, "r") as fd:
        deg = 0
        cmd = ""
        for line in fd:
            line = line.strip("\n") # remove newline
            line = remove_comment(line)
            line = line.strip()
            cmd += ' ' + line
            deg += line.count(OPEN_PAREN) - line.count(CLOSE_PAREN)
            if deg == 0 and cmd.strip() != "":
                eval(parse(cmd))
                cmd = ""
        if deg == 0:
            if cmd.strip() != "":
                eval(parse(cmd))
        else:
            raise ValueError(f"There are {deg} close parens required")
