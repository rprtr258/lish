from typing import List, Union

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
    elif is_atom(ast) or callable(ast):
        # ast
        # but ast is atom (number or string) or callable
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


class Action(str):
    pass


# TODO: add try-catch
# TODO: add implicit progn-s
def EVAL(ast: Expression, env: Env, todo_stack=None, done_stack=None):
    todo_stack = ([(ast, env)] if todo_stack is None else todo_stack)
    done_stack = [] if done_stack is None else done_stack
    while len(todo_stack) > 0:
        # TODO: cmd arg for debug
        if True:
            from tabulate import tabulate
            from textwrap import wrap

            def transpose(m):
                n = max(map(len, [m[0], m[1]]))

                def pr_todo(todo):
                    if isinstance(todo[0], Action):
                        return f"{todo}"
                    else:
                        return f'{(PRINT(todo[0]), "ENV", todo[1])}'
                return [[m[0][i] if i < len(m[0]) else " ", "\n".join(wrap(pr_todo(m[1][i]))) if i < len(m[1]) else " "] for i in range(n)]
            errprint(tabulate(transpose([done_stack, todo_stack]), headers=["DONE", "TODO"], tablefmt="grid"))
            # errprint("DONE:")
            # for done in done_stack:
            #     errprint("  ", done)
            # errprint("TODO:")
            # for todo in todo_stack:
            #     if isinstance(todo[0], Action):
            #         errprint("  ", todo)
            #     else:
            #         errprint("  ", "EVAL", PRINT(todo[0]), "ENV", todo[1])
            # errprint()
            # errprint("TODO:", [todo if isinstance(todo[0], Action) else ("EVAL", PRINT(todo[0])) for todo in todo_stack])
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
                _, binds, ast, env = todo
                let_env = Env(outer=env)
                vals = done_stack[-len(binds):]
                done_stack = done_stack[: -len(binds)]
                for var, value in zip(binds, vals):
                    if isinstance(value, Procedure):
                        value.env = let_env
                    # let_env[var] = value
                    let_env.set(var, value)
                # TODO: test progn and continuation compatibility
                todo_stack.append(([Symbol("progn")] + ast, let_env))  # implicit progn
                continue  # tail call optimisation
            elif code == "FUNCTION CALL":
                _, args_cnt = todo
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
                    todo_stack.append((proc.body, Env(fun_args, fun_exprs, proc.env)))
                    continue  # tail call optimisation
                else:
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
            elif code == "IF":
                _, x, y, env = todo
                predicate = done_stack.pop()
                todo_stack.append((x if predicate else y, env))
                continue
            else:
                print(f"UNKNOWN STACK EXECUTION ACTION: {code} in {todo}")
        assert len(todo) == 2, f"Incorrect stack frame: {todo}"
        ast, env = todo
        # MACROEXPANSION
        # (macroname exps...)
        ast = macroexpand(ast, env)

        if type(ast) != list:
            done_stack.append(eval_ast(ast, env))
            continue  # return
        if ast == []:
            done_stack.append([])
            continue  # return

        form_word = ast[0]
        if form_word == Symbol("quote"):
            # (quote exp) or 'exp
            _, exp = ast
            done_stack.append(exp)
            continue  # return
        elif form_word == Symbol("quasiquote"):
            # (quasiquote exp) or `exp
            _, exp = ast
            todo_stack.append((quasiquote(exp), env))
            continue  # tail call optimisation
        elif form_word == Symbol("if"):
            # (if predicate then else)
            assert len(ast) == 4, f"Wrong if form: {PRINT(ast)}"
            _, p, x, y = ast
            todo_stack.append((Action("IF"), x, y, env))
            todo_stack.append((p, env))
            continue  # tail call optimisation
        elif form_word == Symbol("set!"):  # TODO: rename to set
            # (define var exp)
            _, var, exp = ast
            if not isinstance(var, Symbol):
                raise RuntimeError(f"""Definition name is not a symbol, but a {repr(var)}""")
            todo_stack.append((Action("SET VAR VALUE"), var))
            todo_stack.append((exp, env))
            continue  # ???
        elif form_word == Symbol("let*"):  # TODO: rename to let
            # (let* (v1 e1 v2 e2...) e)
            assert len(ast) >= 3, "Wrong args count to let*"
            _, bindings, *exp = ast
            if len(bindings) % 2 != 0:
                raise RuntimeError(f"let* has {len(bindings)} items as bindings which is not even")
            todo_stack.append((Action("LET*"), bindings[::2], exp, env))
            for i in range(len(bindings) - 2, -1, -2):
                var, var_exp = bindings[i: i + 2]
                todo_stack.append((var_exp, env))
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
            todo_stack.append((macrovalue, env))
            continue  # ???
        elif form_word == Symbol("lambda"):  # TODO: change to fn
            # (lambda (args...) body)
            # TODO: implicit progn
            assert len(ast) == 3, f"Wrong lambda definition: {PRINT(ast)}"
            _, args, body = ast
            for arg in args:
                if not isinstance(arg, Symbol):
                    raise RuntimeError(f"Argument name is not a symbol, but a {repr(arg)}")
            done_stack.append(Procedure(args, body, env))
            continue  # return
        elif form_word == Symbol("call/cc"):
            # (call/cc f)
            assert len(ast) == 2, "call/cc got no or more than one arguments"
            _, f = ast

            class Continuation:
                def __init__(self, env):
                    self.env: Env = env
                    nonlocal todo_stack, done_stack
                    self.todo_stack = todo_stack[::]
                    self.done_stack = done_stack[::]

                def __call__(self, arg):
                    nonlocal todo_stack, done_stack
                    todo_stack = []
                    done_stack = []
                    return EVAL(None, self.env, todo_stack=self.todo_stack[::] + [(arg, self.env)], done_stack=self.done_stack[::])

                def __repr__(self):
                    # return f"Continuation(TODO: {self.todo_stack}, DONE: {self.done_stack})"
                    # return f"Cont(TODOs={len(self.todo_stack)}, DONEs={len(self.done_stack)}, todo={self.todo_stack[-1]})"
                    return f"Cont(TODOs={len(self.todo_stack)}, DONEs={len(self.done_stack)})"
            todo_stack.append(([f, Continuation(env)], env))
            continue  # tail call optimization
        else:
            # (proc arg...)
            arg_exps = ast[1:]
            arg_exps_rev = arg_exps[::-1]
            n = len(arg_exps)
            todo_stack.append((Action("FUNCTION CALL"), n))
            for i, exp in enumerate(arg_exps_rev):
                # TODO: optimise arg_exps[: i] to args[: i] ???, mutability?
                form_word_copy = form_word  # noqa: F841 fuck you python scoping
                todo_stack.append((exp, env))
            todo_stack.append((form_word, env))
            continue  # ???
    assert len(done_stack) == 1
    return done_stack[0]
