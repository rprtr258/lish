import unittest

from context import LispSH
from definitions import A, B, C, QA, QB, QC, QNIL, ATOM_SYMBOL, QUOTE_SYMBOL, EQ_SYMBOL, COND_SYMBOL
from LispSH import read_from_tokens, tokenize, symbol


class TestTokenizer(unittest.TestCase):
    def __tokenizer_test__(self, test_program, expected_result):
        actual_result = read_from_tokens(tokenize(test_program))
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
                [[EQ_SYMBOL, QA, QB], [QUOTE_SYMBOL, symbol("first")]],
                [[ATOM_SYMBOL, QA], [QUOTE_SYMBOL, symbol("second")]]])

    def test_lambda_passing_lambda(self):
        self.__tokenizer_test__("((lambda (f) (f '(b c))) '(lambda (x) (cons 'a x)))", [
            [symbol("lambda"), [symbol("f")],
                [symbol("f"), [QUOTE_SYMBOL, [B, C]]]],
                [QUOTE_SYMBOL,
                    [symbol("lambda"), [symbol("x")],
                        [symbol("cons"), QA, symbol("x")]]]])

    def test_not_enough_close_parens(self):
        with self.assertRaises(ValueError) as cm:
            read_from_tokens(tokenize("(a b"))
        self.assertEqual(str(cm.exception), "Not enough close parens found")

if __name__ == '__main__':
    unittest.main()
