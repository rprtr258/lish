from typing import List, Union
import re

from LiSH.datatypes import Symbol, Keyword, Hashmap


TOKEN_REGEX = re.compile(r"""[\s]*(~@|[\[\]{}()'`~^]|"(?:\\.|[^\\"])*"|;.*|[^\s\[\]{}('"`,;)]+)""")
CLOSE_2_OPEN_PARENS = {
    ')': '(',
    ']': '[',
    '}': '{'
}
OPEN_2_CLOSE_PARENS = {
    '(': ')',
    '[': ']',
    '{': '}'
}
Expression = Union[int, float, str, Symbol, Keyword, Hashmap, List["Expression"]]


class Reader:
    def __init__(self, tokens):
        # skip comments
        self.tokens = list(filter(lambda x: x[0] != ';', tokens))
        self.position = 0

    def has_next(self):
        return self.position < len(self.tokens)

    def peek(self):
        return self.tokens[self.position]

    def next(self):
        res = self.peek()
        self.position += 1
        return res


# TODO: test mirroring
def tokenize(s: str) -> List[str]:
    """Convert a string into a list of tokens

        Args:
            s: string to tokenize

        Returns:
            list of tokens"""
    return TOKEN_REGEX.findall(s)


def read_list(reader):
    L = []
    begin = reader.next()
    constructor, end = {
        '(': (list, ')'),
        '[': (lambda x: [Symbol("list")] + x, ']'),
        '{': (lambda x: [Symbol("hash-map")] + x, '}')
    }[begin]
    while reader.peek() != end:
        L.append(read_form(reader))
    reader.next()  # remove end
    return constructor(L)


def read_atom(token):
    # str
    if token[0] == '"':
        return token[1:-1]
    # bool
    if token in ["true", "false"]:
        return (token == "true")
    if token[0] == ':':
        return Keyword(token[1:])
    if re.match(r"-?[0-9]+$", token):
        return int(token)
    if re.match(r"-?[0-9][0-9.]*$", token):
        return float(token)
    return Symbol(token)


def read_form(reader: Reader) -> Expression:
    """Read an expression from a sequence of tokens

        Args:
            reader: Reader instance with tokens of form to read

        Returns:
            expression that was read

        Raises:
            SyntaxError: in these cases:
                - if unexpected end of input occured
                - if found unexpected close paren"""
    if not reader.has_next():
        raise SyntaxError("Unexpected end of input")
    token = reader.peek()
    # LIST
    if token in ['(', '[', '{']:
        return read_list(reader)
    # READER MACROSES
    # TODO: add reader macroses from config
    token = reader.next()  # remove reader macro
    for reader_macro, macro_word in [
        ("'", "quote"),
        ('`', "quasiquote"),
        ('~', "unquote"),
        ("~@", "splice-unquote")
    ]:
        if token == reader_macro:
            return [Symbol(macro_word), read_form(reader)]
    if token == "^":
        meta = read_list(reader)
        data = read_form(reader)
        return [Symbol("with-meta"), data, meta]
    if token in [')', ']', '}']:
        # TODO: print place
        raise SyntaxError(f"Unexpected {token}")
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
            close_paren = CLOSE_2_OPEN_PARENS[token]
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
        raise SyntaxError("Tokens left after reading whole form, check parens")
    return ast


def READ(line: str) -> Expression:
    """Read an expression from a string

        Args:
            line: string to read expression from

        Returns:
            form that was read"""
    return read_str(line)


def fix_parens(cmd_line):
    cmd_line = cmd_line.strip()
    if cmd_line[0] not in ['(']:
        cmd_line = '(' + cmd_line
    paren_stack = []
    tokens = tokenize(cmd_line)
    for token in tokens:
        if token in ['(', '[', '{']:
            paren_stack.append(token)
        elif token in [')', ']', '}']:
            if len(paren_stack) == 0:
                cmd_line = CLOSE_2_OPEN_PARENS[token] + cmd_line
            elif paren_stack.pop() != CLOSE_2_OPEN_PARENS[token]:
                raise RuntimeError(f"Unexpected {token}")
    while len(paren_stack) > 0:
        cmd_line = cmd_line + OPEN_2_CLOSE_PARENS[paren_stack.pop()]
    return cmd_line
