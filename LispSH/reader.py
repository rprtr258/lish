from typing import List, Any, Union
import re

from LispSH.datatypes import Symbol, Atom


# TODO: remove
OPEN_PAREN = '('
CLOSE_PAREN = ')'
QUOTE = '\''
# TODO: check , in the beginning of the line
# TODO: check "? in the end of string
TOKEN_REGEX = re.compile(r"""[\s,]*(~@|[\[\]{}()'`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"`,;)]+)""")


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
    reader.next() # remove '('
    while reader.peek() != ')':
        L.append(read_form(reader))
    reader.next() # remove ')'
    return L

def read_atom(reader):
    token = reader.next()
    if not reader.has_next():
        raise SyntaxError("Unexpected end of input")
    # str
    if token[0] == '"': return Atom(token[1:-1])
    # bool
    if token in ["true", "false"]: return Atom(token == "true")
    try:
        return Atom(int(token))
    except ValueError:
        pass
    try:
        return Atom(float(token))
    except ValueError:
        pass
    return Symbol(token)

def read_form(reader):
    "Read an expression from a sequence of tokens"
    token = reader.peek()
    if token == '(':
        return read_list(reader)
    elif token == "'":
        reader.next() # remove '
        return ["quote", read_form(reader)]
    elif token == CLOSE_PAREN:
        raise SyntaxError(f"Unexpected {CLOSE_PAREN}")
    else:
        return read_atom(reader)

def check_parens(tokens):
    if tokens[0] not in ['(', "'"]:
        raise SyntaxError("Not a form")
    if tokens[-1] != ')':
        raise SyntaxError("Form is not closed or there is garbage after form")
    parens = list(filter(lambda x: x in ['(', ')'], tokens))
    paren_degree = 0
    for i, paren in enumerate(parens):
        if paren == '(':
            paren_degree += 1
        elif paren == ')':
            paren_degree -= 1
            if paren_degree < 0:
                raise SyntaxError(f"Redundant close paren")
            if paren_degree == 0 and i < len(parens) - 1 and any(map(lambda x: x in ['('], parens[i + 1:])):
                raise SyntaxError(f"Another form found while parsing")
    if paren_degree != 0:
        if paren_degree > 0:
            raise SyntaxError(f"There are {paren_degree} open parens left unclosed")
        elif paren_degree < 0:
            raise SyntaxError(f"There are {-paren_degree} redundant closed parens")

def read_str(line):
    tokens = tokenize(line)
    check_parens(tokens)
    return read_form(Reader(tokens))

def READ(line):
    "Read an expression from a string"
    return read_str(line)
