from typing import List, Union, Any

from LiSH.env import Env
from LiSH.datatypes import Symbol, Hashmap, is_atom
from LiSH.errors import FunctionCallError
from LiSH.printer import PRINT


Body = Union[List["Body"], Symbol, Hashmap, int, float, str]


# TODO: rename to function
"A user-defined Lisp function."
class Procedure:
    # TODO: add keyword args
    def __init__(self, args: List[Symbol], body: Body, env: Env):
        i = 0
        self.pos_args = []
        while i < len(args) and args[i] != "&":
            self.pos_args.append(args[i])
            i += 1
        if i < len(args) and args[i] == "&":
            if i + 1 >= len(args):
                raise RuntimeError("Rest argument is not named")
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

    def __type__(self):
        return "lambda"

    def __str__(self):
        special_args = [] if self.rest_arg is None else [Symbol("&"), self.rest_arg]
        return PRINT([Symbol(self.__type__()), self.pos_args + special_args, self.body])

class Macro(Procedure):
    def __init__(self, fn, env: Env):
        self.__call__ = fn.__call__
        self.rest_arg = fn.rest_arg
        self.pos_args = fn.pos_args
        self.body = fn.body
        self.fn = fn
        self.env = env
    def __str__(self):
        return "MACRO, based on " + str(self.fn)



def macroexpand(ast, env):
    while isinstance(ast, list) and \
        len(ast) > 0 and \
        isinstance(ast[0], Symbol) and \
        (macro := env.get(ast[0])) and \
        isinstance(macro, Macro):
        _, *exprs = ast
        ast = macro(*exprs)
    return ast

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
    elif isinstance(ast, Hashmap):
        res = []
        for k, v in ast.items():
            res.append(EVAL(k, env))
            res.append(EVAL(v, env))
        return Hashmap(res)

def quasiquote(ast):
    if isinstance(ast, list):
        if len(ast) == 2 and ast[0] == Symbol("unquote"):
            return ast[1]
        else:
            res = []
            for x in ast[-1::-1]:
                if isinstance(x, list) and len(x) == 2 and x[0] == Symbol("splice-unquote"):
                    res = [Symbol("+"), x[1], res]
                else:
                    res = [Symbol("cons"), quasiquote(x), res]
            return res
    elif is_atom(ast) or isinstance(ast, Hashmap):
        return [Symbol("quote"), ast]

# TODO: add try-catch
# TODO: add implicit progn-s
def EVAL(ast, env):
    while True:
        # MACROEXPANSION
        # (macroname exps...)
        ast = macroexpand(ast, env)

        if type(ast) != list:
            return eval_ast(ast, env)
        if ast == []:
            return ast

        form_word = ast[0]
        if form_word == Symbol("quote"):
            # (quote exp) or 'exp
            _, exp = ast
            return exp
        if form_word == Symbol("quasiquote"):
            # (quasiquote exp) or `exp
            _, exp = ast
            ast = quasiquote(exp)
            continue # tail call optimisation
        elif form_word == Symbol("cond"): # TODO: test return default in (cond p1 e1 p2 e2 default)
            # (cond p1 e1 p2 e2 ... pn en)
            # or
            # (cond p1 e1 p2 e2 ... pn en default)
            predicates_exps = ast[1:]
            if len(predicates_exps) % 2 == 1:
                predicates_exps, default_value = predicates_exps[:-1], predicates_exps[-1]
            found = False
            for i in range(0, len(predicates_exps), 2):
                predicate, expression = predicates_exps[i : i + 2]
                if EVAL(predicate, env):
                    found = True
                    ast = expression
                    break
            if found:
                continue # tail call optimisation
            # if default value is given
            ast = default_value
            continue # tail call optimisation
        elif form_word == Symbol("set!"): # TODO: rename to set
            # (define var exp)
            _, var, exp = ast
            if not isinstance(var, Symbol):
                raise RuntimeError(f"""Definition name is not a symbol, but a {repr(var)}""")
            value = EVAL(exp, env)
            env.set(var, value)
            return value
        elif form_word == Symbol("let*"): # TODO: rename to let
            # (let* (v1 e1 v2 e2...) e)
            assert len(ast) >= 3, "Wrong args count to let*"
            _, bindings, *exp = ast
            if len(bindings) % 2 != 0:
                raise RuntimeError(f"let* has {len(bindings)} items as bindings which is not even")
            let_env = Env(outer=env)
            for i in range(0, len(bindings), 2):
                var, var_exp = bindings[i : i + 2]
                let_env[var] = EVAL(var_exp, let_env)
            ast, env = [Symbol("progn")] + exp, let_env # implicit progn
            continue # tail call optimisation
        elif form_word == Symbol("macroexpand"):
            # (macroexpand (macro exps...))
            _, macroform = ast
            return macroexpand(macroform, env)
        elif form_word == Symbol("setmacro"):
            # (defmacro macroname (args...) body)
            _, macroname, macrovalue = ast
            if not isinstance(macroname, Symbol):
                raise RuntimeError("Macro definition name is not a symbol")
            macrofn = EVAL(macrovalue, env)
            env.set(macroname, Macro(macrofn, env))
            return [] # TODO: nil
        elif form_word == Symbol("lambda"): # TODO: change to fn
            # (lambda (args...) body)
            _, args, body = ast
            for arg in args:
                if not isinstance(arg, Symbol):
                    raise RuntimeError(f"Argument name is not a symbol, but a {repr(arg)}")
            return Procedure(args, body, env)
        else:
            # (proc arg...)
            proc = EVAL(form_word, env)
            if not callable(proc) and not isinstance(proc, Procedure):
                raise RuntimeError(f"""{proc} (which is {PRINT(ast[0])}) is not a function call in {PRINT(ast)}.""")
            args = [EVAL(exp, env) for exp in ast[1:]]
            try:
                if isinstance(proc, Procedure):
                    if not (len(proc.pos_args) == len(args) or len(proc.pos_args) < len(args) and proc.rest_arg):
                        raise RuntimeError(f"{proc} expected {len(proc.pos_args)} arguments, but got {len(ast[1:])}")
                    if proc.rest_arg is None:
                        fun_args = proc.pos_args
                        fun_exprs = args
                    else:
                        pos_len = len(proc.pos_args)
                        fun_args = proc.pos_args + [proc.rest_arg]
                        fun_exprs = args[:pos_len] + [args[pos_len:]]
                    ast, env = proc.body, Env(fun_args, fun_exprs, proc.env)
                    continue # tail call optimisation
                else:
                    # (proc arg...)
                    # TODO: fix
                    # if len(proc.args) != len(ast[1:]):
                        # raise RuntimeError(f"{proc} expected {len(proc.args)} arguments, but got {len(ast[1:])}")
                    if (res := proc(*args)) is None:
                        print("FUCK YOU,", PRINT(ast), PRINT(args))
                    return res
            except Exception as e:
                raise FunctionCallError(proc, args, e)
