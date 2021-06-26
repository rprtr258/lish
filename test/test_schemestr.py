import unittest

from definitions import NIL
from context import LispSH
from LispSH.printer import schemestr
from LispSH.datatypes import Symbol, Atom


class TestSchemestr(unittest.TestCase):
    def __test_schemestr__(self, expression, expected_str):
        self.assertEqual(schemestr(expression), expected_str)

    def test_nil(self):
        self.__test_schemestr__(NIL, "nil")

    def test_quote(self):
        self.__test_schemestr__([Symbol("quote"), Symbol("a")], "'a")

    def test_quote_list(self):
        self.__test_schemestr__([Symbol("quote"), [Symbol("a"), Symbol("b"), Symbol("c")]], "'(a b c)")

    def test_symbol(self):
        self.__test_schemestr__(Symbol("a"), "a")

    def test_atom_str(self):
        self.__test_schemestr__(Atom("a"), "\"a\"")

    def test_atom_int(self):
        self.__test_schemestr__(Atom(42), "42")

    def test_list(self):
        self.__test_schemestr__([Symbol("f"), Atom(42)], "(f 42)")

if __name__ == '__main__':
    unittest.main()
