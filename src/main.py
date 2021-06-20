from typing import Union, List
import math
import operator as op

################ Types

Symbol = str          # A Lisp Symbol is implemented as a Python str
List   = list         # A Lisp List is implemented as a Python list
Number = Union[int, float] # A Lisp Number is implemented as a Python int or float

################ Constants

OPEN_PAREN = '('
CLOSE_PAREN = ')'
QUOTE = '\''

################ Parsing: parse, tokenize, and read_from_tokens

def tokenize(s: str) -> List[str]:
    "Convert a string into a list of tokens."
    return s \
        .replace(OPEN_PAREN, f" {OPEN_PAREN} ") \
        .replace(CLOSE_PAREN, f" {CLOSE_PAREN} ") \
        .replace(QUOTE, f" {QUOTE} ") \
        .split()

def atom(token: str) -> Union[Symbol, Number]:
    "Numbers become numbers; every other token is a symbol."
    try:
        return int(token)
    except ValueError:
        try:
            return float(token)
        except ValueError:
            return Symbol(token)

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
        return ["quote"] + L
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
            return self[var]
        if self.outer is None:
            return None # nil
        return self.outer.find(var)

def standard_env():
    "An environment with some Scheme standard procedures."
    env = Env()
    env.update(vars(math)) # sin, cos, sqrt, pi, ...
    env.update({
        '+':op.add, '-':op.sub, '*':op.mul, '/':op.truediv,
        '>':op.gt, '<':op.lt, '>=':op.ge, '<=':op.le, '=':op.eq,
        'abs':     abs,
        'append':  op.add,
        'apply':   lambda fx: fx[0](*fx[1:]),
        'begin':   lambda *x: x[-1],
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
        'null?':   lambda x: x == [],
        'number?': lambda x: isinstance(x, Number),
        'procedure?': callable,
        'round':   round,
        'symbol?': lambda x: isinstance(x, Symbol),
    })
    return env

global_env = standard_env()

################ Interaction: A REPL

def repl(prompt='lis.py> '):
    "A prompt-read-eval-print loop."
    while True:
        val = eval(parse(input(prompt)))
        if val is not None: 
            print(lispstr(val))

def lispstr(exp):
    "Convert a Python object back into a Lisp-readable string."
    if isinstance(exp, List):
        return OPEN_PAREN + ' '.join(map(lispstr, exp)) + CLOSE_PAREN
    else:
        return str(exp)

################ Procedures

class Procedure(object):
    "A user-defined Scheme procedure."
    def __init__(self, parms, body, env):
        self.parms, self.body, self.env = parms, body, env
    def __call__(self, *args): 
        return eval(self.body, Env(self.parms, args, self.env))

################ eval

def eval(x, env=global_env):
    "Evaluate an expression in an environment."
    if isinstance(x, Symbol):      # variable reference
        return env.find(x)
    elif not isinstance(x, List):  # constant literal
        return x
    elif x[0] == "quote":          # (quote exp)
        _, exp = x
        return exp
    elif x[0] == "atom":          # (atom exp)
        _, exp = x
        exp = eval(exp, env)
        return \
            isinstance(exp, str) or \
            isinstance(exp, bool) or \
            exp == []
    # TODO: define through cond
    # elif x[0] == "if":             # (if test conseq alt)
        # _, test, conseq, alt = x
        # exp = (conseq if eval(test, env) else alt)
        # return eval(exp, env)
    elif x[0] == "cond": # TODO: return default in (cond p1 e1 p2 e2 default)
        predicates_exps = x[1:]
        i = 0
        while i < len(predicates_exps):
            predicate, expression = predicates_exps[i : i + 2]
            i += 2
            if eval(predicate, env):
                return eval(expression, env)
    elif x[0] == "define":         # (define var exp)
        _, var, exp = x
        env[var] = eval(exp, env)
    elif x[0] == "set!":           # (set! var exp)
        _, var, exp = x
        env.find(var)[var] = eval(exp, env)
    elif x[0] == "lambda":         # (lambda (var...) body)
        _, parms, body = x
        return Procedure(parms, body, env)
    else:                          # (proc arg...)
        proc = eval(x[0], env)
        args = [eval(exp, env) for exp in x[1:]]
        if not callable(proc):
            print(proc)
            raise ValueError(f"{x} is not a function call")
        return proc(*args)

def fix_parens(cmd_line):
    if cmd_line[0] != OPEN_PAREN:
        cmd_line = OPEN_PAREN + cmd_line
    if cmd_line[-1] != CLOSE_PAREN:
        cmd_line = cmd_line + CLOSE_PAREN
    # TODO: don't count brackets in strings
    open_parens, close_parens = cmd_line.count(OPEN_PAREN), cmd_line.count(CLOSE_PAREN)
    return \
        OPEN_PAREN * max(0, close_parens - open_parens) + \
        cmd_line + \
        CLOSE_PAREN * max(0, open_parens - close_parens)

def schemestr(exp):
    "Convert a Python object back into a Scheme-readable string."
    if isinstance(exp, List):
        return '(' + ' '.join(map(schemestr, exp)) + ')' 
    else:
        return str(exp)

def repl(prompt='lis.py> '):
    "A prompt-read-eval-print loop."
    while True:
        input_line = input(prompt)
        input_line = fix_parens(input_line)
        val = eval(parse(input_line))
        if val is not None: 
            print(schemestr(val))

if __name__ == "__main__":
    repl()
