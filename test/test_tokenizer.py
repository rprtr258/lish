import unittest

from context import LispSH
from definitions import A, B, C, QA, QB, QC, QNIL, ATOM_SYMBOL, QUOTE_SYMBOL, EQ_SYMBOL, COND_SYMBOL
from LispSH.reader import READ, tokenize
from LispSH.datatypes import Symbol, Atom


class TestTokenizer(unittest.TestCase):
    def __tokenizer_test__(self, test_program, expected_result):
        actual_result = READ(test_program)
        self.assertEqual(actual_result, expected_result)

    def test_atom_quote_list(self):
        self.__tokenizer_test__("(atom '(a b c))", [ATOM_SYMBOL, [QUOTE_SYMBOL, [A, B, C]]])

    def test_list_border_spaces(self):
        self.__tokenizer_test__("(  a b )", [A, B])

    def test_list_inner_spaces(self):
        self.__tokenizer_test__("(a  b    c)", [A, B, C])

    def test_list_outer_spaces(self):
        self.__tokenizer_test__("   (a b)  ", [A, B])

    def test_eq_quoted_equal(self):
        self.__tokenizer_test__("(eq? 'a 'a)", [EQ_SYMBOL, QA, QA])

    def test_eq_quoted_nils(self):
        self.__tokenizer_test__("(eq? '() '())", [EQ_SYMBOL, QNIL, QNIL])

    def test_quoted_nil(self):
        self.__tokenizer_test__("'()", QNIL)

    def test_cond(self):
        self.__tokenizer_test__(
            "(cond ((eq? 'a 'b) 'first) ((atom 'a) 'second))",
            [COND_SYMBOL,
                [[EQ_SYMBOL, QA, QB], [QUOTE_SYMBOL, Symbol("first")]],
                [[ATOM_SYMBOL, QA], [QUOTE_SYMBOL, Symbol("second")]]])

    def test_lambda_passing_lambda(self):
        self.__tokenizer_test__("((lambda (f) (f '(b c))) '(lambda (x) (cons 'a x)))", [
            [Symbol("lambda"), [Symbol("f")],
                [Symbol("f"), [QUOTE_SYMBOL, [B, C]]]],
            [QUOTE_SYMBOL,
                [Symbol("lambda"), [Symbol("x")],
                    [Symbol("cons"), QA, Symbol("x")]]]])

    def test_not_enough_close_parens(self):
        with self.assertRaises(SyntaxError) as cm:
            READ("(a b")
        self.assertEqual(str(cm.exception), "There are 1 open parens left unclosed")

    def test_too_much_close_parens(self):
        with self.assertRaises(SyntaxError) as cm:
            READ("(a b))")
        self.assertEqual(str(cm.exception), "There are 1 redundant closed parens")

    def test_tokenize(self):
        self.assertEqual(
            tokenize('(+ "a" "(a b))))")'),
            ['(', '+', '"a"', '"(a b))))"', ')'])
        M_SLASH = '\\\\'
        M_DQUOTE = '\\"'
        self.assertEqual(
            tokenize(f'(+ "{M_SLASH}{M_DQUOTE}" "abc")'),
            ['(', '+', f'"{M_SLASH}{M_DQUOTE}"', f'"abc"', ')'])
        self.assertEqual(
            tokenize(f'(+ "(" "{M_DQUOTE}" ")" "{M_SLASH}")'),
            ['(', '+', '"("', f'"{M_DQUOTE}"', '")"', f'"{M_SLASH}"', ')'])

    def test_string(self):
        self.__tokenizer_test__('(+ "a" "(a b))))")', [Symbol("+"), Atom("a"), Atom("(a b))))")])

if __name__ == '__main__':
    unittest.main()
