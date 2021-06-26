from typing import List, Any, Union

from LispSH.datatypes import atom_or_symbol, Symbol


OPEN_PAREN = '('
CLOSE_PAREN = ')'
QUOTE = '\''

def no_quote_replace(s: str, c: str, p: str):
    "Replaces c(char) to p(pattern) in s if c not in double quotes"
    i = 0
    res = ""
    quoted = False
    while i < len(s):
        if s[i] == '"':
            res += '"'
            quoted = not quoted
        elif s[i] == c and not quoted:
            res += p
        else:
            res += s[i]
        i += 1
    return res

def remove_comment(s: str) -> str:
    if ';' in s:
        return s[:s.find(';')] # TODO: check "wOwOw ;;;; " ; fuck you
    return s

# TODO: test mirroring
def tokenize(s: str) -> List[str]:
    "Convert a string into a list of tokens."
    s = remove_comment(s)
    word = None
    res = []
    quoted = False
    s = no_quote_replace(s, OPEN_PAREN, f" {OPEN_PAREN} ")
    s = no_quote_replace(s, CLOSE_PAREN, f" {CLOSE_PAREN} ")
    s = no_quote_replace(s, QUOTE, f" {QUOTE} ")
    mirrored = False
    for c in s:
        if c in [' ', '\n']:
            if not quoted:
                if not word is None:
                    res.append(word)
                    word = None
                mirrored = False
            else:
                word += c
        elif c == '\\' and quoted:
            if mirrored:
                word += '\\'
                mirrored = False
            else:
                mirrored = True
        else:
            word = c if word is None else word + c
            if c == '"' and not mirrored:
                if quoted:
                    res.append(word)
                    word = None
                quoted = not quoted
            mirrored = False
    if not word is None:
        res.append(word)
    return res

def read_from_tokens(tokens):
    "Read an expression from a sequence of tokens."
    if len(tokens) == 0:
        raise SyntaxError('unexpected EOF while reading')
    token = tokens.pop(0)
    if token == OPEN_PAREN:
        L = []
        if len(tokens) == 0:
            raise ValueError("Not enough close parens found")
        while tokens[0] != CLOSE_PAREN:
            L.append(read_from_tokens(tokens))
            if len(tokens) == 0:
                raise ValueError("Not enough close parens found")
        tokens.pop(0) # pop off ')'
        return L
    elif token == QUOTE:
        L = []
        L.append(read_from_tokens(tokens))
        return [Symbol("quote")] + L
    elif token == CLOSE_PAREN:
        raise SyntaxError(f'unexpected {CLOSE_PAREN}')
    else:
        return atom_or_symbol(token)

def parse(program):
    "Read an expression from a string."
    return read_from_tokens(tokenize(program))
