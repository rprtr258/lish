from typing import List, Any, Union
import re

from LispSH.datatypes import Symbol, Atom, Keyword, Vector, Hashmap


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
    begin = reader.next()
    constructor, end = {
        '(': (list, ')'),
        '[': (Vector, ']'),
        '{': (Hashmap, '}')
    }[begin]
    while reader.peek() != end:
        L.append(read_form(reader))
    reader.next() # remove end
    return constructor(L)

def read_atom(token):
    # str
    if token[0] == '"': return Atom(token[1:-1])
    # bool
    if token in ["true", "false"]: return Atom(token == "true")
    if token[0] == ':': return Keyword(token[1:])
    if re.match(r"-?[0-9]+$", token):
        return Atom(int(token))
    if re.match(r"-?[0-9][0-9.]*$", token):
        return Atom(float(token))
    return Symbol(token)

def read_form(reader):
    "Read an expression from a sequence of tokens"
    if not reader.has_next():
        raise SyntaxError("Unexpected end of input")
    token = reader.peek()
    # LIST
    if token in ['(', '[', '{']:
        return read_list(reader)
    # COMMENT
    if token[0] == ";": return []
    # READER MACROSES
    token = reader.next() # remove reader macro
    if token == "'": return ["quote", read_form(reader)]
    if token == '`': return ["quasiquote", read_form(reader)]
    if token == '~': return ["unquote", read_form(reader)]
    if token == "~@": return ["splice-unquote", read_form(reader)]
    if token == "@": return ["deref", read_form(reader)]
    if token == "^":
        meta = read_list(reader)
        data = read_form(reader)
        return ["with-meta", data, meta]
    if token == CLOSE_PAREN:
        raise SyntaxError(f"Unexpected {CLOSE_PAREN}")
    # ATOM
    return read_atom(token)

# TODO: make normal
def check_parens(tokens):
    if tokens.count('(') != tokens.count(')'):
        raise SyntaxError("Different number of open and close parens")
    if tokens.count('(') > 0 and tokens[-1] != ')':
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
