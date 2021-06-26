from dataclasses import dataclass
from typing import List, Any

from LispSH.env import Env
from LispSH.datatypes import Symbol, Macro, Vector, Hashmap, is_atom
from LispSH.printer import PRINT


# TODO: rename to function
"A user-defined Lisp function."
@dataclass
class Procedure:
    args: List[str]
    body: Any
    env: Env
    def __call__(self, *args): 
        return eval_ast(self.body, Env(self.args, args, self.env))
    def __str__(self):
        return PRINT([Symbol("lambda"), self.args, self.body])


# TODO: variadic defun (defn (& x) x) and macro (defmacro (& x) x)
# TODO: (loop ... recur) or Tail call optimisation(harder, non-recursive eval?)

def macroexpand(macroform, env):
    macroname, *exps = macroform
    res = env.get(macroname)
    macroargs, macrobody = res.args, res.body
    return eval_ast(macrobody, Env(macroargs, exps, env))

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
        return Vector([eval_ast(x, env) for x in ast])
    elif isinstance(ast, Hashmap):
        res = []
        for k, v in ast.items():
            res.append(eval_ast(k, env))
            res.append(eval_ast(v, env))
        return Hashmap(res)
    elif isinstance(ast[0], Symbol) and (res := env.get(ast[0])) and isinstance(res, Macro):
        # (macroname exps...)
        macroexpansion = macroexpand(ast, env)
        return eval_ast(macroexpansion, env)
    else:
        form_word = ast[0]
        if form_word == Symbol("quote"):
            # (quote exp)
            _, exp = ast
            return exp
        elif form_word == Symbol("cond"): # TODO: test return default in (cond p1 e1 p2 e2 default)
            # (cond p1 e1 p2 e2 ... pn en)
            # or
            # (cond p1 e1 p2 e2 ... pn en default)
            predicates_exps = ast[1:]
            i = 0
            while i + 1 < len(predicates_exps):
                predicate, expression = predicates_exps[i : i + 2]
                i += 2
                if eval_ast(predicate, env):
                    return eval_ast(expression, env)
            # if default value is given
            if len(predicates_exps) % 2 == 1:
                return eval_ast(predicates_exps[-1], env)
        elif form_word == Symbol("set!"):
            # (define var exp)
            _, var, exp = ast
            assert isinstance(var, Symbol), f"""Definition name is not a symbol, but a {PRINT(var)}"""
            value = eval_ast(exp, env)
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
                let_env.set(var, eval_ast(var_exp, let_env))
                i += 2
            return eval_ast(exp, let_env)
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
            proc = eval_ast(proc, env)
            args = eval_ast(args, env)
            try:
                return proc(*args)
            except Exception as e:
                raise RuntimeError(f"""Error during evaluation ({proc} {" ".join(map(PRINT, args))}).
Error is:
    {"Recursed" if isinstance(e, RuntimeError) else e}""")
        else:
            # (proc arg...)
            proc = eval_ast(form_word, env)
            args = [eval_ast(exp, env) for exp in ast[1:]]
            if not callable(proc) and not isinstance(proc, Procedure):
                raise RuntimeError(f"""{proc} (which is {PRINT(ast[0])}) is not a function call in {PRINT(ast)}.""")
            try:
                if (res := proc(*args)) is None:
                    print("FUCK YOU,", PRINT(ast), PRINT(args))
                return res
            except Exception as e:
                raise RuntimeError(f"""Error during evaluation ({proc} {" ".join(map(PRINT, args))}).
Error is:
    {"Recursed" if isinstance(e, RuntimeError) else e}""")

def EVAL(ast, env):
    if ast == []:
        return ast
    return eval_ast(ast, env)
