from dataclasses import dataclass
from typing import List, Any

from LispSH.env import global_env, Env
from LispSH.datatypes import Symbol, Atom, Macro
from LispSH.printer import schemestr


# TODO: rename to function
"A user-defined Lisp function."
@dataclass
class Procedure:
    args: List[str]
    body: Any
    env: Env
    def __call__(self, *args): 
        return eval(self.body, Env(self.args, args, self.env))
    def __str__(self):
        return schemestr([Symbol("lambda"), self.args, self.body])


# TODO: variadic defun (defn (& x) x) and macro (defmacro (& x) x)
# TODO: (loop ... recur) or Tail call optimisation(harder, non-recursive eval?)

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

#@log_eval
# TODO: remove global_env from global vars
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
            if not callable(proc) and not isinstance(proc, Procedure):
                raise RuntimeError(f"""{proc} (named {schemestr(x[0])}) is not a function call in {schemestr(x)}.""")
            try:
                if (res := proc(*args)) is None:
                    print("FUCK YOU,", schemestr(x), schemestr(args))
                return res
            except Exception as e:
                raise RuntimeError(f"""Error during evaluation ({proc} {" ".join(map(schemestr, args))}).
Error is:
    {"Recursed" if isinstance(e, RuntimeError) else e}""")