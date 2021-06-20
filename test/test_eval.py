import unittest

from definitions import A, B, C, NIL, TRUE, FALSE
from context import LispSH
from LispSH import parse, eval, symbol


class TestEval(unittest.TestCase):
    def __eval_test__(self, test_program, expected_result):
        actual_result = eval(parse(test_program))
        self.assertEqual(actual_result, expected_result)

    def test_quote(self):
        self.__eval_test__("(quote a)", A)

    def test_quote_list(self):
        self.__eval_test__("(quote (a b c))", [A, B, C])

    def test_quote_nested_list(self):
        self.__eval_test__("(quote (a b (c d (f e))))", [A, B, [C, symbol("d"), [symbol("f"), symbol("e")]]])

    def test_quote_nil(self):
        self.__eval_test__("(quote ())", NIL)

    def test_atom_quote(self):
        self.__eval_test__("(atom (quote a))", TRUE)

    def test_atom_quote_list(self):
        self.__eval_test__("(atom (quote (a b c)))", FALSE)

    def test_atom_quote_nil(self):
        self.__eval_test__("(atom (quote ()))", TRUE)

    def test_atom_quoted(self):
        self.__eval_test__("(atom 'a)", TRUE)

    def test_atom_quoted_list(self):
        self.__eval_test__("(atom '(a b c))", FALSE)

    def test_atom_quoted_nil(self):
        self.__eval_test__("(atom '())", TRUE)

    def test_atom_atom_quoted(self):
        self.__eval_test__("(atom (atom 'a))", TRUE)

    def test_atom_quoted_atom(self):
        self.__eval_test__("(atom '(atom 'a))", FALSE)

    def test_eq_equal_symbols(self):
        self.__eval_test__("(eq? 'a 'a)", TRUE)

    def test_eq_not_equal_symbols(self):
        self.__eval_test__("(eq? 'a 'b)", FALSE)

    def test_eq_nils(self):
        self.__eval_test__("(eq? '() '())", TRUE)

    def test_car(self):
        self.__eval_test__("(car '(a b c))", A)

    def test_cdr(self):
        self.__eval_test__("(cdr '(a b c))", [B, C])

    def test_cdr_nil(self):
        self.__eval_test__("(cdr '(a))", NIL)

    def test_cons(self):
        self.__eval_test__("(cons 'a '(b c))", [A, B, C])

    def test_cons_cons_cons(self):
        self.__eval_test__("(cons 'a (cons 'b (cons 'c '())))", [A, B, C])

    def test_car_cons(self):
        self.__eval_test__("(car (cons 'a '(b c)))", A)

    def test_cdr_cons(self):
        self.__eval_test__("(cdr (cons 'a '(b c)))", [B, C])

    def test_cond_true(self):
        self.__eval_test__("(cond (eq? 'a 'a) 'first 'second)", symbol("first"))

    def test_cond_false(self):
        self.__eval_test__("(cond (eq? 'a 'b) 'first 'second)", symbol("second"))

    def test_lambda_one_argument(self):
        self.__eval_test__("((lambda (x) (cons x '(b))) 'a)", [A, B])

    def test_lambda_two_arguments(self):
        self.__eval_test__("((lambda (x y) (cons x (cdr y))) 'z '(a b c))", [symbol("z"), B, C])

    def test_lambda_passing_lambda(self):
        self.__eval_test__("((lambda (f) (f '(b c))) (lambda (x) (cons 'a x)))", [A, B, C])

    # def test_label(self):
        # self.__eval_test__(
            # """((label subst
                # (lambda (x y z)
                    # (cond
                        # ((atom z) (cond
                            # ((eq z y) x)
                            # ('t z)))
                        # ('t (cons (subst x y (car z))
                            # (subst x y (cdr z)))))))
                # 'm 'b '(a b (a b c) d))""",
            # [A, symbol("m"), [A, symbol("m"), C], symbol("d")])

    # def test_defun_subst(self):
        # self.__eval_test__(
        # """((defun subst (x y z)
            # (cond ((atom z) (cond ((eq z y) x)
                    # ('t z)))
                # ('t (cons
                    # (subst x y (car z))
                    # (subst x y (cdr z))))))
            # 'm 'b '(a b (a b c) d))""",
            # [A, symbol("m"), [A, symbol("m"), C], symbol("d")])

    # def test_defun_cadr(self):
        # self.__eval_test__("""((defun cadr (x) (car (cdr x))) '((a b) (c d) e))""", [C, symbol("d")])

if __name__ == '__main__':
    unittest.main()
