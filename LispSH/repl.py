from sys import stdin, stdout

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
    while True:
        try:
            print_prompt(env)
            line = input()
            if line.strip() == "": continue
            line = fix_parens(line)
            print(rep(line, env))
        except Exception as e:
            print(f"{type(e).__name__}: {e}")

################ File load
# TODO: move to some load function
def load_file(filename, env):
    with open(filename, "r") as fd:
        deg = 0
        cmd = ""
        for line in fd:
            line = remove_comment(line)
            line = line.strip("\n").strip()
            cmd += ' ' + line
            deg += line.count(OPEN_PAREN) - line.count(CLOSE_PAREN)
            if deg == 0 and cmd.strip() != "":
                rep(cmd, env)
                cmd = ""
        if deg == 0:
            if cmd.strip() != "":
                rep(cmd, env)
        else:
            raise ValueError(f"There are {deg} close parens required")
