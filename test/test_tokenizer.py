import unittest

from context import LispSH
from definitions import A, B, C, QA, QB, QC, QNIL, ATOM_SYMBOL, QUOTE_SYMBOL, EQ_SYMBOL, COND_SYMBOL
from LispSH.reader import READ, tokenize
from LispSH.datatypes import Symbol, Keyword, Hashmap


class TestTokenizer(unittest.TestCase):
    def __tokenizer_test__(self, test_program, expected_result):
        actual_result = READ(test_program)
        self.assertEqual(actual_result, expected_result)

    def test_atom_quote_list(self):
        self.__tokenizer_test__("(atom '(a b c))", [ATOM_SYMBOL, [QUOTE_SYMBOL, [A, B, C]]])

    def test_list_spaces(self):
        self.__tokenizer_test__("(  a b )", [A, B])
        self.__tokenizer_test__("(a  b    c)", [A, B, C])
        self.__tokenizer_test__("   (a b)  ", [A, B])

    def test_eq(self):
        self.__tokenizer_test__("(eq? 'a 'a)", [EQ_SYMBOL, QA, QA])
        self.__tokenizer_test__("(eq? '() '())", [EQ_SYMBOL, QNIL, QNIL])

    def test_quoted_nil(self):
        self.__tokenizer_test__("'()", QNIL)

    def test_keyword(self):
        self.__tokenizer_test__(":ab", Keyword("ab"))
        self.__tokenizer_test__("(:ab :cd :kwwww)", [Keyword("ab"), Keyword("cd"), Keyword("kwwww")])

    def test_vector(self):
        self.__tokenizer_test__("[]", [Symbol("list")])
        self.__tokenizer_test__("   [    ]  ", [Symbol("list")])
        self.__tokenizer_test__("[+ 1 2]", [Symbol("list"), Symbol("+"), 1, 2])
        self.__tokenizer_test__("[[a b]]", [Symbol("list"), [Symbol("list"), A, B]])
        self.__tokenizer_test__("[+ 1 [* a b]]", [Symbol("list"), Symbol("+"), 1, [Symbol("list"), Symbol("*"), A, B]])
        self.__tokenizer_test__("   [   +   1   [ *   a   b  ]   ]    ", [Symbol("list"), Symbol("+"), 1, [Symbol("list"), Symbol("*"), A, B]])
        self.__tokenizer_test__("([])", [[Symbol("list")]])

    def test_hashmap(self):
        A_A, A_B, A_C = "a", "b", "c"
        self.__tokenizer_test__("{}", Hashmap([]))
        self.__tokenizer_test__("  {   } ", Hashmap([]))
        self.__tokenizer_test__('{"abc" 1}', Hashmap(["abc", 1]))
        self.__tokenizer_test__('{"a" {"b" 2}}', Hashmap([A_A, Hashmap([A_B, 2])]))
        self.__tokenizer_test__('{"a" {"b" {"c" 3}}}', Hashmap([A_A, Hashmap([A_B, Hashmap([A_C, 3])])]))
        self.__tokenizer_test__('{  "a"  {"b"   {  "cde"     3   }  }}', Hashmap([A_A, Hashmap([A_B, Hashmap(["cde", 3])])]))
        self.__tokenizer_test__('{"a1" 1 "a2" 2 "a3" 3}', Hashmap(["a1", 1, "a2", 2, "a3", 3]))
        self.__tokenizer_test__('{  :a  {:b   {  :cde     3   }  }}', Hashmap([Keyword("a"), Hashmap([Keyword("b"), Hashmap([Keyword("cde"), 3])])]))
        self.__tokenizer_test__('{"1" 1}', Hashmap(["1", 1]))
        self.__tokenizer_test__("({})", [Hashmap([])])

    def test_comments(self):
        self.__tokenizer_test__(";wow", [])
        self.__tokenizer_test__(" ;;ff", [])
        self.__tokenizer_test__("1 ;;ff", 1)
        self.__tokenizer_test__("1; ff", 1)

    def test_deref(self):
        self.__tokenizer_test__("@a", [Symbol("deref"), A])

    def test_deref(self):
        self.__tokenizer_test__('^{"a" 1} [1 2 3]', [Symbol("with-meta"), [Symbol("list"), 1, 2, 3], Hashmap(["a", 1])])

    def test_quote(self):
        self.__tokenizer_test__("'a", [QUOTE_SYMBOL, A])
        self.__tokenizer_test__("'(a 1 bc)", [QUOTE_SYMBOL, [A, 1, Symbol("bc")]])

    def test_quasiquote(self):
        self.__tokenizer_test__("`a", [Symbol("quasiquote"), A])
        self.__tokenizer_test__("`(a 1 bc)", [Symbol("quasiquote"), [A, 1, Symbol("bc")]])

    def test_unquote(self):
        self.__tokenizer_test__("~a", [Symbol("unquote"), A])
        self.__tokenizer_test__("~(a 1 bc)", [Symbol("unquote"), [A, 1, Symbol("bc")]])
        self.__tokenizer_test__("`(a 1 ~bc)", [Symbol("quasiquote"), [A, 1, [Symbol("unquote"), Symbol("bc")]]])

    def test_splice_unquote(self):
        self.__tokenizer_test__("~@(a 1 bc)", [Symbol("splice-unquote"), [A, 1, Symbol("bc")]])

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
        self.assertEqual(str(cm.exception), "Different number of open and close parens")

    def test_too_much_close_parens(self):
        with self.assertRaises(SyntaxError) as cm:
            READ("(a b))")
        self.assertEqual(str(cm.exception), "Different number of open and close parens")

    def test_unclosed_string(self):
        with self.assertRaises(SyntaxError) as cm:
            READ('(echo )"abc)')
        self.assertEqual(str(cm.exception), "Tokens left after reading whole form, check parens")

    def test_two_forms(self):
        with self.assertRaises(SyntaxError) as cm:
            READ('(+ 1 2)(* 2 3)')
        self.assertEqual(str(cm.exception), "Tokens left after reading whole form, check parens")

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
        self.__tokenizer_test__('(+ "a" "(a b))))")', [Symbol("+"), "a", "(a b))))"])

if __name__ == '__main__':
    unittest.main()
