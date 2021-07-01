from sys import stdout, argv
import os.path
import inspect

from LiSH.env import Env
from LiSH.reader import OPEN_PAREN, CLOSE_PAREN, QUOTE, READ
from LiSH.datatypes import Symbol
from LiSH.evaluator import EVAL
from LiSH.printer import PRINT


# TODO: move to reader
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


def print_prompt(env):
    print(EVAL([Symbol("prompt")], env), end="")
    stdout.flush()


def rep(line: str, env: Env) -> str:
    """Read, Eval, Print line using provided environment

        Args:
            line: line of Lish to execute
            env: environment to run EVAL in

        Returns:
            result of executing line pretty printed"""
    return PRINT(EVAL(READ(line), env))


def trace_error(e, debug=False):
    if isinstance(e, RuntimeError):
        print(f"Error: {e}")
    else:
        print(f"{type(e).__name__}: {e}")
    if debug:
        print()
        print("=" * 34 + "STACK FRAMES" + "=" * 34)
        print()
        frames = inspect.trace()
        for frame in frames:
            args = inspect.getargvalues(frame.frame)
            locals = [arg for arg in args.locals.keys() if arg not in args.args]
            print("=" * 80)
            filename = os.path.relpath(os.path.abspath(frame.filename), os.path.abspath("."))
            print(f'  File "{filename}", line {frame.lineno}, in {frame.function}(')
            for arg in args.args:
                arg_value = "\n      ".join(map(lambda x: x.strip(), str(args.locals[arg]).split('\n')))
                print(f"    {arg}={arg_value}")
            print("  )")
            if len(locals) > 0:
                print("  Locals:")
                for arg in locals:
                    arg_value = args.locals[arg]
                    print(f"    {arg}={repr(arg_value)}")
            for code_line in frame.code_context:
                print(f"{'(' + str(frame.lineno) + ')':6s}" + code_line.strip())
            print("=" * 80)
            print()


# TODO: add Ctrl-D support
# TODO: Shift-Enter for multiline input
# TODO: line editing, parens?
def repl(env: Env):
    """A prompt-read-eval-print loop.

        Args:
            env: environment to run repl with"""
    try:
        env[Symbol("*argv*")] = argv
        env[Symbol("*debug*")] = False
        env[Symbol("eval")] = lambda ast: EVAL(ast, env)
        rep('(set! load-file (lambda (f) (eval (read (+ "(progn " (slurp f) "\n)")))))', env)
        rep('(load-file ".lishrc")', env)
        while True:
            try:
                print_prompt(env)
                line = input()
                if line.strip() == "":
                    continue
                line = fix_parens(line)
                print(rep(line, env))
            except Exception as e:
                trace_error(e, env[Symbol("*debug*")])
    except Exception as e:
        trace_error(e, env[Symbol("*debug*")])
