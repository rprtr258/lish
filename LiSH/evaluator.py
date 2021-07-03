from typing import Callable, List, Union

from LiSH.env import Env
from LiSH.datatypes import Symbol, Hashmap, is_atom
from LiSH.errors import FunctionCallError
from LiSH.errprint import errprint
from LiSH.reader import Expression
from LiSH.printer import PRINT


Body = Union[List["Body"], Symbol, Hashmap, int, float, str]


# TODO: rename to function
class Procedure:
    """A user-defined Lisp function."""
    # TODO: add keyword args
    def __init__(self, args: List[Symbol], body: Body, env: Env, stack):
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
        self.stack = stack

    def __call__(self, *args):
        args = list(args)
        if self.rest_arg is None:
            return EVAL(self.body, Env(self.pos_args, args, self.env), [], self.stack + [[str(self), args]])
        else:
            pos_len = len(self.pos_args)
            fun_args = self.pos_args + [self.rest_arg]
            fun_exprs = args[:pos_len] + [args[pos_len:]]
            return EVAL(self.body, Env(fun_args, fun_exprs, self.env))

    def __str__(self):
        special_args = [] if self.rest_arg is None else [Symbol("&"), self.rest_arg]
        return PRINT([Symbol("lambda"), self.pos_args + special_args, self.body])


class Macro(Procedure):
    def __init__(self, fn, env: Env, stack):
        self.__call__ = fn.__call__
        self.rest_arg = fn.rest_arg
        self.pos_args = fn.pos_args
        self.body = fn.body
        self.fn = fn
        self.env = env
        self.stack = stack

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


def eval_ast(ast: Expression, env: Env) -> Expression:
    """Evaluate an expression in an environment.

        Args:
            ast: expression to evaluate
            env: environment to use to lookup variable values

        Returns:
            result of evaluating expression

        Raises:
            RuntimeError: if expression is a symbol which value is not in environment"""
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


def gensym():
    i = 1
    while True:
        yield Symbol(f"#SYM_{i}")
        i += 1


# TODO: add try-catch
# TODO: add implicit progn-s
def EVAL(ast: Expression, env: Env, zipper: List[Callable[[Expression], Expression]] = [], stack: List[Expression] = []):
    if False:  # TODO: cmd arg log cur eval
        errprint(f"=========== EVAL {PRINT(ast)} ===========")
    if type(ast) == list:
        if False:  # TODO: cmd arg log eval stack
            # DEBUG CALL STACK
            errprint("=========== STACK ===========")
            for frame in stack:
                print(frame)
                errprint(PRINT(frame))
            errprint("=========== STACK ===========")
    gs = gensym()
    while True:
        # MACROEXPANSION
        # (macroname exps...)
        ast = macroexpand(ast, env)

        if type(ast) != list:
            return eval_ast(ast, env)
        if ast == []:
            return ast
        if False:  # TODO: cmd arg log eval env
            errprint(env)
        if False:  # TODO: cmd arg log eval continuation
            # DEBUG CONTINUATION
            errprint("=========== CONTINUATION ===========")
            errprint(PRINT(ast), ":")
            for zipp in zipper:
                errprint("  ", PRINT(zipp(Symbol("<>"))))
            errprint("OR AS CONTINUATION:")
            res = Symbol("<>")
            for zipp in zipper[::-1]:
                res = zipp(res)
            errprint("  ", PRINT(res))
            errprint("=========== CONTINUATION ===========")

        form_word = ast[0]
        if form_word == Symbol("quote"):
            # (quote exp) or 'exp
            _, exp = ast
            return exp
        if form_word == Symbol("quasiquote"):
            # (quasiquote exp) or `exp
            _, exp = ast
            ast = quasiquote(exp)
            continue  # tail call optimisation
        elif form_word == Symbol("cond"):  # TODO: test return default in (cond p1 e1 p2 e2 default)
            # (cond p1 e1 p2 e2 ... pn en)
            # or
            # (cond p1 e1 p2 e2 ... pn en default)
            predicates_exps = ast[1:]
            if len(predicates_exps) % 2 == 1:
                predicates_exps, default_value = predicates_exps[:-1], predicates_exps[-1]
            found = False
            for i in range(0, len(predicates_exps), 2):
                predicate, expression = predicates_exps[i: i + 2]
                cur = lambda c: [Symbol("cond")] + predicates_exps[: i] + [c] + predicates_exps[i + 1:]  # noqa E731, zipper cond predicate
                if EVAL(predicate, env, zipper + [cur], stack + [ast]):
                    found = True
                    ast = expression
                    new_cur = lambda c: [Symbol("cond")] + predicates_exps[: i + 1] + [c] + predicates_exps[i + 2:]  # noqa E731, zipper cond expression
                    break
            if found:
                zipper = zipper + [new_cur]
                continue  # tail call optimisation
            # TODO: throw error if default value is not given
            new_cur = lambda c: [Symbol("cond")] + predicates_exps + [c]  # noqa E731, zipper cond default
            ast, zipper = default_value, zipper + [new_cur]
            continue  # tail call optimisation
        elif form_word == Symbol("set!"):  # TODO: rename to set
            # (define var exp)
            _, var, exp = ast
            if not isinstance(var, Symbol):
                raise RuntimeError(f"""Definition name is not a symbol, but a {repr(var)}""")
            cur = lambda c: [Symbol("set!")] + [var] + [c]  # noqa E731, zipper value
            value = EVAL(exp, env, zipper + [cur], stack + [ast])
            env.set(var, value)
            return value
        elif form_word == Symbol("let*"):  # TODO: rename to let
            # (let* (v1 e1 v2 e2...) e)
            assert len(ast) >= 3, "Wrong args count to let*"
            _, bindings, *exp = ast
            if len(bindings) % 2 != 0:
                raise RuntimeError(f"let* has {len(bindings)} items as bindings which is not even")
            let_env = Env(outer=env)
            for i in range(0, len(bindings), 2):
                var, var_exp = bindings[i: i + 2]
                cur = lambda c: [Symbol("let*")] + [bindings[: i + 1] + [c] + bindings[i + 2:]] + exp  # noqa E731, zipper var value
                let_env[var] = EVAL(var_exp, let_env, zipper + [cur], stack + [ast])
            new_cur = lambda c: [Symbol("let*")] + [bindings] + [c]  # noqa E731, zipper let expression
            # TODO: test progn and continuation compatibility
            ast, env, zipper = [Symbol("progn")] + exp, let_env, zipper + [new_cur]  # implicit progn
            continue  # tail call optimisation
        elif form_word == Symbol("macroexpand"):
            # (macroexpand (macro exps...))
            _, macroform = ast
            return macroexpand(macroform, env)
        elif form_word == Symbol("setmacro"):
            # (defmacro macroname (args...) body)
            _, macroname, macrovalue = ast
            if not isinstance(macroname, Symbol):
                raise RuntimeError("Macro definition name is not a symbol")
            macrofn = EVAL(macrovalue, env, zipper, stack + [ast])
            env.set(macroname, Macro(macrofn, env, stack))
            return []  # TODO: nil
        elif form_word == Symbol("lambda"):  # TODO: change to fn
            # (lambda (args...) body)
            _, args, body = ast
            for arg in args:
                if not isinstance(arg, Symbol):
                    raise RuntimeError(f"Argument name is not a symbol, but a {repr(arg)}")
            # TODO: pass zipper
            return Procedure(args, body, env, stack)
        elif form_word == Symbol("call/cc"):
            # (call/cc f)
            assert len(ast) == 2, "call/cc got no or more than one arguments"
            _, f = ast
            cont_arg = next(gs)
            continuation = cont_arg
            for cont in zipper[::-1]:
                continuation = cont(continuation)
            # TODO: skip original continuation
            ast = [f, [Symbol("lambda"), [cont_arg], continuation]]
            continue  # tail call optimization
        else:
            # (proc arg...)
            arg_exps = ast[1:]
            cur = lambda c: [c] + arg_exps  # noqa E731, zipper procedure
            proc = EVAL(form_word, env, zipper + [cur], stack + [ast])
            if not callable(proc) and not isinstance(proc, Procedure):
                raise RuntimeError(f"""{proc} (which is {PRINT(ast[0])}) is not a function call in {PRINT(ast)}.""")
            args = []
            for i, exp in enumerate(arg_exps):
                # TODO: optimise arg_exps[: i] to args[: i] ???, mutability?
                cur = lambda c: [form_word] + args + [c] + arg_exps[i + 1:]  # noqa E731, zipper procedure argument
                args.append(EVAL(exp, env, zipper + [cur], stack + [ast]))
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
                    continue  # tail call optimisation
                else:
                    # TODO: fix
                    # if len(proc.args) != len(ast[1:]):
                    #     raise RuntimeError(f"{proc} expected {len(proc.args)} arguments, but got {len(ast[1:])}")
                    if (res := proc(*args)) is None:
                        print("FUCK YOU,", PRINT(ast), PRINT(args))
                    return res
            except Exception as e:
                raise FunctionCallError(proc, args, e)
