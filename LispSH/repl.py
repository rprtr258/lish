from sys import stdin, stdout, argv
import os.path
import inspect

from LispSH.reader import remove_comment, OPEN_PAREN, CLOSE_PAREN, QUOTE, READ
from LispSH.datatypes import Symbol
from LispSH.evaluator import EVAL
from LispSH.printer import PRINT

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

def rep(line, env):
    "Read, Eval, Print line"
    return PRINT(EVAL(READ(line), env))

# TODO: add Ctrl-D support
# TODO: Shift-Enter for multiline input
# TODO: line editing, parens?
def repl(env):
    "A prompt-read-eval-print loop."
    env[Symbol("*argv*")] = argv
    env[Symbol("eval")] = lambda ast: EVAL(ast, env)
    rep('(set! load-file (lambda (f) (eval (read (+ "(progn " (slurp f) ")")))))', env)
    rep('(load-file ".lisprc")', env)
    while True:
        try:
            print_prompt(env)
            line = input()
            if line.strip() == "": continue
            line = fix_parens(line)
            print(rep(line, env))
        except Exception as e:
            print(f"{type(e).__name__}: {e}")
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
                    print("Locals:")
                    for arg in locals:
                        arg_value = args.locals[arg]
                        print(f"    {arg}={repr(arg_value)}")
                for code_line in frame.code_context:
                    print("    " + code_line.strip())
                print("=" * 80)
                print()
