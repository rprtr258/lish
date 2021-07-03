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
        args = list(args)
        if self.rest_arg is None:
            return EVAL(self.body, Env(self.pos_args, args, self.env))
        else:
            pos_len = len(self.pos_args)
            fun_args = self.pos_args + [self.rest_arg]
            fun_exprs = args[:pos_len] + [args[pos_len:]]
            return EVAL(self.body, Env(fun_args, fun_exprs, self.env))

    def __repr__(self):
        special_args = [] if self.rest_arg is None else [Symbol("&"), self.rest_arg]
        return PRINT([Symbol("lambda"), self.pos_args + special_args, self.body])

    def __str__(self):
        special_args = [] if self.rest_arg is None else [Symbol("&"), self.rest_arg]
        return PRINT([Symbol("lambda"), self.pos_args + special_args, self.body])


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


class Action(str): pass


def debug_continuation(zipper):
    errprint("=========== CONTINUATION ===========")
    continuation = Symbol("<>")
    for cont in zipper[::-1]:
        continuation = cont(continuation)
    errprint("CONTINUATION:", PRINT(continuation))
    errprint("BY STEPS CONTINUATION:")
    for cont in zipper[::-1]:
        errprint("  ", PRINT(cont(Symbol("<>"))))
    errprint("=========== CONTINUATION ===========")


# TODO: add try-catch
# TODO: add implicit progn-s
def EVAL(ast: Expression, env: Env, zipper: List[Callable[[Expression], Expression]] = []):
    if False:  # TODO: cmd arg log cur eval
        errprint(f"=========== EVAL {PRINT(ast)} ===========")
    gs = gensym()
    todo_stack = [(ast, env, zipper)]
    done_stack = []
    while len(todo_stack) > 0:
        # errprint("DONE:", done_stack)
        # errprint("TODO:", [todo if isinstance(todo[0], Action) else ("EVAL", PRINT(todo[0])) for todo in todo_stack])
        # errprint("TODO:")
        # for todo in todo_stack:
        #     if isinstance(todo[0], Action):
        #         errprint(todo)
        #     else:
        #         errprint("EVAL", PRINT(todo[0]), "ENV", todo[1])
        # errprint()
        todo = todo_stack.pop()
        if isinstance(todo[0], Action):
            code = todo[0]
            if code == "SET VAR VALUE":
                _, var = todo
                value = done_stack.pop()
                env.set(var, value)  # ???
                done_stack.append(value)
                continue  # return
            elif code == "LET*":
                _, binds, ast, env, zipper = todo
                # debug_continuation(zipper)
                let_env = Env(outer=env)
                vals = done_stack[-len(binds):]
                done_stack = done_stack[: -len(binds)]
                for var, value in zip(binds, vals):
                    if isinstance(value, Procedure):
                        value.env = let_env
                    let_env[var] = value
                new_cur = lambda c: [Symbol("let*")] + [binds] + [vals] + [c]  # noqa E731, zipper let expression
                # TODO: test progn and continuation compatibility
                todo_stack.append(([Symbol("progn")] + ast, let_env, zipper + [new_cur]))  # implicit progn
                continue  # tail call optimisation
            elif code == "FUNCTION CALL":
                _, args_cnt, zipper = todo
                # debug_continuation(zipper)
                proc = done_stack[-args_cnt - 1]
                args = done_stack[-args_cnt:] if args_cnt > 0 else []
                done_stack = done_stack[: -args_cnt - 1]
                if not callable(proc) and not isinstance(proc, Procedure):
                    raise RuntimeError(f"""{proc} (which is {PRINT(ast[0])}) is not a function call in {PRINT(ast)}.""")
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
                    todo_stack.append((proc.body, Env(fun_args, fun_exprs, proc.env), zipper))
                    continue  # tail call optimisation
                else:
                    # TODO: fix
                    # if len(proc.args) != len(ast[1:]):
                    #     raise RuntimeError(f"{proc} expected {len(proc.args)} arguments, but got {len(ast[1:])}")
                    try:
                        if (res := proc(*args)) is None:
                            print("FUCK YOU,", PRINT(ast), PRINT(args))
                    except Exception as e:
                        raise FunctionCallError(proc, args, e)
                    done_stack.append(res)
                    continue  # return
            elif code == "MACRO SET":
                _, macroname = todo
                macrofn = done_stack.pop()
                env.set(macroname, Macro(macrofn, env))
                done_stack.append([])
                continue
            else:
                print(f"UNKNOWN STACK EXECUTION ACTION: {code} in {todo}")
        assert len(todo) == 3, f"Incorrect stack frame: {todo}"
        ast, env, zipper = todo
        # MACROEXPANSION
        # (macroname exps...)
        ast = macroexpand(ast, env)

        if type(ast) != list:
            done_stack.append(eval_ast(ast, env))
            continue  # return
        if ast == []:
            done_stack.append([])
            continue  # return
        if False:  # TODO: cmd arg log eval env
            errprint(env)
        if False:  # TODO: cmd arg log eval continuation
            debug_continuation(zipper)

        form_word = ast[0]
        if form_word == Symbol("quote"):
            # (quote exp) or 'exp
            _, exp = ast
            done_stack.append(exp)
            continue  # return
        elif form_word == Symbol("quasiquote"):
            # (quasiquote exp) or `exp
            _, exp = ast
            todo_stack.append((quasiquote(exp), env, zipper))
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
                if EVAL(predicate, env, zipper + [cur]):  # TODO: move to todo_stack
                    found = True
                    ast = expression
                    new_cur = lambda c: [Symbol("cond")] + predicates_exps[: i + 1] + [c] + predicates_exps[i + 2:]  # noqa E731, zipper cond expression
                    break
            if found:
                todo_stack.append((ast, env, zipper + [new_cur]))
                continue  # tail call optimisation
            # TODO: throw error if default value is not given
            new_cur = lambda c: [Symbol("cond")] + predicates_exps + [c]  # noqa E731, zipper cond default
            todo_stack.append((default_value, env, zipper + [new_cur]))
            continue  # tail call optimisation
        elif form_word == Symbol("set!"):  # TODO: rename to set
            # (define var exp)
            _, var, exp = ast
            if not isinstance(var, Symbol):
                raise RuntimeError(f"""Definition name is not a symbol, but a {repr(var)}""")
            cur = lambda c: [Symbol("set!")] + [var] + [c]  # noqa E731, zipper value
            todo_stack.append((Action("SET VAR VALUE"), var))
            todo_stack.append((exp, env, zipper + [cur]))
            continue  # ???
        elif form_word == Symbol("let*"):  # TODO: rename to let
            # (let* (v1 e1 v2 e2...) e)
            assert len(ast) >= 3, "Wrong args count to let*"
            _, bindings, *exp = ast
            if len(bindings) % 2 != 0:
                raise RuntimeError(f"let* has {len(bindings)} items as bindings which is not even")
            todo_stack.append((Action("LET*"), bindings[::2], exp, env, zipper))
            for i in range(len(bindings) - 2, -1, -2):
                var, var_exp = bindings[i: i + 2]
                cur = lambda c: [Symbol("let*")] + [bindings[: i + 1] + [c] + bindings[i + 2:]] + [exp]  # noqa E731, zipper var value
                todo_stack.append((var_exp, env, zipper + [cur]))
            continue  # ???(tail call optimisation)
        elif form_word == Symbol("macroexpand"):
            # (macroexpand (macro exps...))
            _, macroform = ast
            return macroexpand(macroform, env)
        elif form_word == Symbol("setmacro"):
            # (setmacro macroname f)
            _, macroname, macrovalue = ast
            if not isinstance(macroname, Symbol):
                raise RuntimeError("Macro definition name is not a symbol")
            todo_stack.append((Action("MACRO SET"), macroname))
            todo_stack.append((macrovalue, env, zipper))
            continue  # ???
        elif form_word == Symbol("lambda"):  # TODO: change to fn
            # (lambda (args...) body)
            # TODO: implicit progn
            assert len(ast) == 3, f"Wrong lambda definition: {PRINT(ast)}"
            _, args, body = ast
            for arg in args:
                if not isinstance(arg, Symbol):
                    raise RuntimeError(f"Argument name is not a symbol, but a {repr(arg)}")
            # TODO: pass zipper
            done_stack.append(Procedure(args, body, env))
            continue  # return
        elif form_word == Symbol("call/cc"):
            # (call/cc f)
            assert len(ast) == 2, "call/cc got no or more than one arguments"
            _, f = ast
            cont_arg = next(gs)
            continuation = cont_arg
            for cont in zipper[::-1]:
                continuation = cont(continuation)
            # TODO: skip original continuation
            todo_stack = [([f, [Symbol("lambda"), [cont_arg], continuation]], env, zipper)]
            done_stack = []
            continue  # tail call optimization
        else:
            # (proc arg...)
            arg_exps = ast[1:]
            arg_exps_rev = arg_exps[::-1]
            n = len(arg_exps)
            todo_stack.append((Action("FUNCTION CALL"), n, zipper))
            for i, exp in enumerate(arg_exps_rev):
                # TODO: optimise arg_exps[: i] to args[: i] ???, mutability?
                form_word_copy = form_word # fuck you python scoping
                cur = lambda c: [form_word_copy] + arg_exps[: n - i - 1] + [c] + arg_exps[n - i:]  # noqa E731, zipper procedure argument
                todo_stack.append((exp, env, zipper + [cur]))
            cur = lambda c: [c] + arg_exps  # noqa E731, zipper procedure
            todo_stack.append((form_word, env, zipper + [cur]))
            continue  # ???
    assert len(done_stack) == 1
    return done_stack[0]
