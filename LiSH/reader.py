from typing import List, Any, Union
import re

from LiSH.datatypes import Symbol, Keyword, Hashmap


# TODO: remove
OPEN_PAREN = '('
CLOSE_PAREN = ')'
QUOTE = '\''
TOKEN_REGEX = re.compile(r"""[\s]*(~@|[\[\]{}()'`~^]|"(?:\\.|[^\\"])*"|;.*|[^\s\[\]{}('"`,;)]+)""")


# TODO: remove, use tokenize instead
def remove_comment(s: str) -> str:
    if ';' in s:
        return s[:s.find(';')] # TODO: check (echo "wOwOw ;;;; ") ; fuck you
    return s

class Reader:
    def __init__(self, tokens):
        self.tokens = tokens
        self.position = 0

    def has_next(self):
        return self.position < len(self.tokens)

    def peek(self):
        # skip comments
        while self.tokens[self.position][0] == ";":
            self.position += 1
        return self.tokens[self.position]

    def next(self):
        res = self.peek()
        self.position += 1
        return res

# TODO: test mirroring
def tokenize(s: str) -> List[str]:
    "Convert a string into a list of tokens"
    return TOKEN_REGEX.findall(s)

def read_list(reader):
    L = []
    begin = reader.next()
    constructor, end = {
        '(': (list, ')'),
        '[': (lambda x: [Symbol("list")] + x, ']'),
        '{': (Hashmap, '}')
    }[begin]
    while reader.peek() != end:
        L.append(read_form(reader))
    reader.next() # remove end
    return constructor(L)

def read_atom(token):
    # str
    if token[0] == '"': return token[1:-1]
    # bool
    if token in ["true", "false"]: return (token == "true")
    if token[0] == ':': return Keyword(token[1:])
    if re.match(r"-?[0-9]+$", token): return int(token)
    if re.match(r"-?[0-9][0-9.]*$", token): return float(token)
    return Symbol(token)

def read_form(reader):
    "Read an expression from a sequence of tokens"
    if not reader.has_next():
        raise SyntaxError("Unexpected end of input")
    token = reader.peek()
    # LIST
    if token in ['(', '[', '{']:
        return read_list(reader)
    # READER MACROSES
    # TODO: add reader macroses from config
    token = reader.next() # remove reader macro
    for reader_macro, macro_word in [
        ("'", "quote"),
        ('`', "quasiquote"),
        ('~', "unquote"),
        ("~@", "splice-unquote")]:
        if token == reader_macro:
            return [Symbol(macro_word), read_form(reader)]
    if token == "^":
        meta = read_list(reader)
        data = read_form(reader)
        return [Symbol("with-meta"), data, meta]
    if token == CLOSE_PAREN:
        # TODO: print place
        raise SyntaxError(f"Unexpected {CLOSE_PAREN} in {rest}")
    # ATOM
    return read_atom(token)

# TODO: make normal
def check_parens(tokens):
    if tokens.count('(') != tokens.count(')'):
        raise SyntaxError("Different number of open and close parens")
    stack = []
    for token in tokens:
        if token in ['(', '[', '{']:
            stack.append(token)
        elif token in [')', ']', '}']:
            close_paren = {
                ')': '(',
                ']': '[',
                '}': '{'
            }[token]
            if len(stack) == 0 or stack.pop() != close_paren:
                raise SyntaxError(f"Unexpected {token}")
    if len(stack) != 0:
        raise SyntaxError(f"There are {stack} parens unclosed")

def read_str(line):
    tokens = tokenize(line)
    check_parens(tokens)
    reader = Reader(tokens)
    ast = read_form(reader)
    if reader.has_next():
        # TODO: print error place
        raise SyntaxError(f"Tokens left after reading whole form, check parens")
    return ast

def READ(line):
    "Read an expression from a string"
    return read_str(line)
