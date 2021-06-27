from sys import stdin, stdout, argv

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
    env.set(Symbol("*argv*"), argv)
    env.set(Symbol("eval"), lambda ast: EVAL(ast, env))
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
            import logging
            logging.exception(f"{type(e).__name__}: {e}")
