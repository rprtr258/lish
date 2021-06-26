from typing import List, Union

from LispSH.env import Env
from LispSH.datatypes import Symbol, Macro, Vector, Hashmap, is_atom
from LispSH.printer import PRINT


Body = Union[List["Body"], Symbol, Vector, Hashmap, int, float, str]


def FunctionCallError(proc, args, e):
    form = " ".join(map(PRINT, [proc] + args))
    error = "Recursed" if isinstance(e, RuntimeError) else e
    return RuntimeError(f"""Error during evaluation ({form}).
Error is:
{error}""")

# TODO: rename to function
"A user-defined Lisp function."
class Procedure:
    # TODO: add keyword args
    def __init__(self, args: List[str], body: Body, env: Env):
        i = 0
        self.pos_args = []
        while i < len(args) and args[i] != "&":
            self.pos_args.append(args[i])
            i += 1
        if i < len(args) and args[i] == "&":
            assert i + 1 < len(args), "Rest argument is not named"
            self.rest_arg: Symbol = args[i + 1]
        else:
            self.rest_arg = None
        self.body: Body = body
        self.env: Env = env

    def __call__(self, *args):
        if self.rest_arg is None:
            return EVAL(self.body, Env(self.pos_args, args, self.env))
        else:
            args = list(args)
            pos_len = len(self.pos_args)
            fun_args = self.pos_args + [self.rest_arg]
            fun_exprs = args[:pos_len] + [args[pos_len:]]
            return EVAL(self.body, Env(fun_args, fun_exprs, self.env))

    def __str__(self):
        if self.rest_arg is None:
            return PRINT([Symbol("lambda"), self.pos_args, self.body])
        else:
            return PRINT([Symbol("lambda"), self.pos_args + [Symbol("&"), self.rest_arg], self.body])


# TODO: variadic defun (defn (& x) x) and macro (defmacro (& x) x)
# TODO: (loop ... recur) or Tail call optimisation(harder, non-recursive eval?)

def macroexpand(macroform, env):
    macroname, *exps = macroform
    res = env.get(macroname)
    macroargs, macrobody = res.args, res.body
    return EVAL(macrobody, Env(macroargs, exps, env))

def eval_ast(ast, env):
    "Evaluate an expression in an environment."
    if isinstance(ast, Symbol):
        # ast
        # but ast is symbol
        res_env = env.find(ast)
        if res_env is None:
            raise RuntimeError(f"{ast} value not found")
        return res_env[ast]
    elif is_atom(ast):
        # ast
        # but ast is atom (number or string)
        return ast
    elif isinstance(ast, Vector):
        return Vector([EVAL(x, env) for x in ast])
    elif isinstance(ast, Hashmap):
        res = []
        for k, v in ast.items():
            res.append(EVAL(k, env))
            res.append(EVAL(v, env))
        return Hashmap(res)

def EVAL(ast, env):
    if type(ast) != list:
        return eval_ast(ast, env)
    if ast == []:
        return ast
    form_word = ast[0]
    if form_word == Symbol("quote"):
        # (quote exp)
        _, exp = ast
        return exp
    elif isinstance(ast[0], Symbol) and (res := env.get(ast[0])) and isinstance(res, Macro):
        # (macroname exps...)
        macroexpansion = macroexpand(ast, env)
        return EVAL(macroexpansion, env)
    elif form_word == Symbol("cond"): # TODO: test return default in (cond p1 e1 p2 e2 default)
        # (cond p1 e1 p2 e2 ... pn en)
        # or
        # (cond p1 e1 p2 e2 ... pn en default)
        predicates_exps = ast[1:]
        i = 0
        while i + 1 < len(predicates_exps):
            predicate, expression = predicates_exps[i : i + 2]
            i += 2
            if EVAL(predicate, env):
                return EVAL(expression, env)
        # if default value is given
        if len(predicates_exps) % 2 == 1:
            return EVAL(predicates_exps[-1], env)
    elif form_word == Symbol("set!"):
        # (define var exp)
        _, var, exp = ast
        assert isinstance(var, Symbol), f"""Definition name is not a symbol, but a {PRINT(var)}"""
        value = EVAL(exp, env)
        env.set(var, value)
        return value
    elif form_word == Symbol("let*"):
        # (let* (v1 e1 v2 e2...) e)
        _, bindings, exp = ast
        assert len(bindings) % 2 == 0
        i = 0
        let_env = Env(outer=env)
        while i < len(bindings):
            var, var_exp = bindings[i : i + 2]
            let_env.set(var, EVAL(var_exp, let_env))
            i += 2
        return EVAL(exp, let_env)
    elif form_word == Symbol("macroexpand"):
        # (macroexpand (macro exps...))
        _, macroform = ast
        return macroexpand(macroform, env)
    elif form_word == Symbol("defmacro"):
        # (defmacro macroname (args...) body)
        _, macroname, args, body = ast
        assert isinstance(macroname, Symbol), "Macro definition name is not a symbol"
        env[macroname] = Macro(args, body)
        return [] # TODO: nil
    elif form_word == Symbol("lambda"):
        # (lambda (args...) body)
        _, args, body = ast
        for arg in args:
            assert isinstance(arg, Symbol), f"Argument name is not a symbol, but a {PRINT(arg)}"
        return Procedure(args, body, env)
    elif form_word == Symbol("apply"):
        # (apply f (args...))
        _, proc, args = ast
        proc = EVAL(proc, env)
        args = EVAL(args, env)
        # TODO: fix
        # if len(proc.args) != len(args):
            # raise RuntimeError(f"{proc} expected {len(proc.args)} arguments, but got {len(args)}")
        try:
            return proc(*args)
        except Exception as e:
            raise FunctionCallError(proc, args, e)
    else:
        # (proc arg...)
        proc = EVAL(form_word, env)
        # TODO: fix
        # if len(proc.args) != len(ast[1:]):
            # raise RuntimeError(f"{proc} expected {len(proc.args)} arguments, but got {len(ast[1:])}")
        args = [EVAL(exp, env) for exp in ast[1:]]
        if not callable(proc) and not isinstance(proc, Procedure):
            raise RuntimeError(f"""{proc} (which is {PRINT(ast[0])}) is not a function call in {PRINT(ast)}.""")
        try:
            if (res := proc(*args)) is None:
                print("FUCK YOU,", PRINT(ast), PRINT(args))
            return res
        except Exception as e:
            raise FunctionCallError(proc, args, e)
