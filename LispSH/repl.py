from sys import stdin, stdout

from LispSH.reader import remove_comment, OPEN_PAREN, CLOSE_PAREN, QUOTE, READ
from LispSH.datatypes import Symbol
from LispSH.evaluator import EVAL
from LispSH.printer import PRINT

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

def print_prompt():
    print(EVAL([Symbol("prompt")]).value, end="")
    stdout.flush()

def rep(line):
    "Read, Eval, Print line"
    return PRINT(EVAL(READ(line)))

# TODO: add Ctrl-D support
# TODO: Shift-Enter for multiline input
def repl():
    "A prompt-read-eval-print loop."
    while True:
        try:
            print_prompt()
            line = input()
            if line.strip() == "": continue
            line = fix_parens(line)
            print(rep(line))
        except Exception as e:
            print(e)

################ File load

def load_file(filename):
    with open(filename, "r") as fd:
        deg = 0
        cmd = ""
        for line in fd:
            line = line.strip("\n") # remove newline
            line = remove_comment(line)
            line = line.strip()
            cmd += ' ' + line
            deg += line.count(OPEN_PAREN) - line.count(CLOSE_PAREN)
            if deg == 0 and cmd.strip() != "":
                EVAL(READ(cmd))
                cmd = ""
        if deg == 0:
            if cmd.strip() != "":
                EVAL(READ(cmd))
        else:
            raise ValueError(f"There are {deg} close parens required")
