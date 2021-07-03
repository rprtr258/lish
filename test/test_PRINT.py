import unittest

from definitions import NIL
from LiSH.printer import PRINT, pr_str
from LiSH.datatypes import Symbol


class TestPRINT(unittest.TestCase):
    def __test_PRINT_no_escape__(self, expression, expected_str):
        self.assertEqual(pr_str(expression, None), expected_str)

    def __test_PRINT__(self, expression, expected_str):
        self.assertEqual(PRINT(expression), expected_str)

    def test_pr_str_noescape(self):
        self.__test_PRINT_no_escape__("a", "a")
        self.__test_PRINT_no_escape__("a\nb", "a\nb")

    def test_nil(self):
        self.__test_PRINT__(NIL, "()")

    def test_quote(self):
        self.__test_PRINT__([Symbol("quote"), Symbol("a")], "'a")

    def test_quote_list(self):
        self.__test_PRINT__([Symbol("quote"), [Symbol("a"), Symbol("b"), Symbol("c")]], "'(a b c)")

    def test_symbol(self):
        self.__test_PRINT__(Symbol("a"), "a")

    def test_atom_str(self):
        self.__test_PRINT__("a", "\"a\"")
        self.__test_PRINT__("a\nb", "\"a\\nb\"")

    def test_atom_int(self):
        self.__test_PRINT__(42, "42")

    def test_list(self):
        self.__test_PRINT__([Symbol("f"), 42], "(f 42)")


if __name__ == '__main__':
    unittest.main()
