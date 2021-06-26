import unittest

from definitions import A, B, C, NIL, TRUE, FALSE
from context import LispSH
from LispSH.reader import READ
from LispSH.evaluator import EVAL
from LispSH.datatypes import Symbol


class TestEVAL(unittest.TestCase):
    def __EVAL_test__(self, test_program, expected_result):
        actual_result = EVAL(READ(test_program))
        self.assertEqual(actual_result, expected_result)

    def test_quote(self):
        self.__EVAL_test__("(quote a)", A)

    def test_quote_list(self):
        self.__EVAL_test__("(quote (a b c))", [A, B, C])

    def test_quote_nested_list(self):
        self.__EVAL_test__("(quote (a b (c d (f e))))", [A, B, [C, Symbol("d"), [Symbol("f"), Symbol("e")]]])

    def test_quote_nil(self):
        self.__EVAL_test__("(quote ())", NIL)

    def test_atom_quote(self):
        self.__EVAL_test__("(atom? (quote a))", TRUE)

    def test_atom_quote_list(self):
        self.__EVAL_test__("(atom? (quote (a b c)))", FALSE)

    def test_atom_quote_nil(self):
        self.__EVAL_test__("(atom? (quote ()))", TRUE)

    def test_atom_quoted(self):
        self.__EVAL_test__("(atom? 'a)", TRUE)

    def test_atom_quoted_list(self):
        self.__EVAL_test__("(atom? '(a b c))", FALSE)

    def test_atom_quoted_nil(self):
        self.__EVAL_test__("(atom? '())", TRUE)

    def test_atom_atom_quoted(self):
        self.__EVAL_test__("(atom? (atom? 'a))", TRUE)

    def test_atom_quoted_atom(self):
        self.__EVAL_test__("(atom? '(atom? 'a))", FALSE)

    def test_eq_equal_Symbols(self):
        self.__EVAL_test__("(= 'a 'a)", TRUE)

    def test_eq_not_equal_Symbols(self):
        self.__EVAL_test__("(= 'a 'b)", FALSE)

    def test_eq_nils(self):
        self.__EVAL_test__("(= '() '())", TRUE)

    def test_car(self):
        self.__EVAL_test__("(car '(a b c))", A)

    def test_cdr(self):
        self.__EVAL_test__("(cdr '(a b c))", [B, C])

    def test_cdr_nil(self):
        self.__EVAL_test__("(cdr '(a))", NIL)

    def test_cons(self):
        self.__EVAL_test__("(cons 'a '(b c))", [A, B, C])

    def test_cons_cons_cons(self):
        self.__EVAL_test__("(cons 'a (cons 'b (cons 'c '())))", [A, B, C])

    def test_car_cons(self):
        self.__EVAL_test__("(car (cons 'a '(b c)))", A)

    def test_cdr_cons(self):
        self.__EVAL_test__("(cdr (cons 'a '(b c)))", [B, C])

    def test_cond_true(self):
        self.__EVAL_test__("(cond (= 'a 'a) 'first 'second)", Symbol("first"))

    def test_cond_false(self):
        self.__EVAL_test__("(cond (= 'a 'b) 'first 'second)", Symbol("second"))

    def test_lambda_one_argument(self):
        self.__EVAL_test__("((lambda (x) (cons x '(b))) 'a)", [A, B])

    def test_lambda_two_arguments(self):
        self.__EVAL_test__("((lambda (x y) (cons x (cdr y))) 'z '(a b c))", [Symbol("z"), B, C])

    def test_lambda_passing_lambda(self):
        self.__EVAL_test__("((lambda (f) (f '(b c))) (lambda (x) (cons 'a x)))", [A, B, C])

    def test_ariphmetic_operators(self):
        self.__EVAL_test__("(+ 1 2 3)", 6)
        self.__EVAL_test__("(+ 1 2 -3)", 0)
        self.__EVAL_test__("(* 1 2 -3)", -6)
        self.__EVAL_test__("(/ -12 2 3)", -2)

    def test_math_functions(self):
        random_result = EVAL(READ("(rand)"))
        self.assertTrue(isinstance(random_result, float))
        self.assertTrue(0 <= random_result and random_result < 1)
        self.__EVAL_test__("(abs -2)", 2)
        self.__EVAL_test__("(abs 2)", 2)
        self.__EVAL_test__("(cos 0)", 1.0)
        self.__EVAL_test__("(max 1 2 3)", 3)
        self.__EVAL_test__("(min 1 2 3)", 1)
        self.__EVAL_test__("(round 3.14)", 3)

    def test_comparison_operators(self):
        self.__EVAL_test__("(> 1 2 3)", FALSE)
        self.__EVAL_test__("(> 3 2 1)", TRUE)
        self.__EVAL_test__("(> 3 2 3)", FALSE)
        self.__EVAL_test__("(> 3 2 2 1)", FALSE)
        self.__EVAL_test__("(< 1 2 3)", TRUE)
        self.__EVAL_test__("(< 3 2 1)", FALSE)
        self.__EVAL_test__("(< 3 2 3)", FALSE)
        self.__EVAL_test__("(< 3 2 2 1)", FALSE)
        self.__EVAL_test__("(>= 1 2 3)", FALSE)
        self.__EVAL_test__("(>= 3 2 1)", TRUE)
        self.__EVAL_test__("(>= 3 2 2 1)", TRUE)
        self.__EVAL_test__("(<= 1 2 3)", TRUE)
        self.__EVAL_test__("(<= 3 2 1)", FALSE)
        self.__EVAL_test__("(<= 3 2 3)", FALSE)
        self.__EVAL_test__("(<= 1 2 2 3)", TRUE)
        self.__EVAL_test__("(= 3 2 3)", FALSE)
        self.__EVAL_test__("(= 2 (+ 1 1) (- 4 2))", TRUE)
        self.__EVAL_test__("(nil? '())", TRUE)
        self.__EVAL_test__("(nil? '(1))", FALSE)
        self.__EVAL_test__("(list? '())", TRUE)
        self.__EVAL_test__("(list? '(1))", TRUE)
        self.__EVAL_test__("(number? '())", FALSE)
        self.__EVAL_test__("(number? '(1))", FALSE)
        self.__EVAL_test__("(number? 'a)", FALSE)
        self.__EVAL_test__("(number? 11)", TRUE)
        self.__EVAL_test__("(number? 11.22)", TRUE)
        self.__EVAL_test__("(symbol? '())", FALSE)
        self.__EVAL_test__("(symbol? '(1))", FALSE)
        self.__EVAL_test__("(symbol? 'a)", TRUE)
        self.__EVAL_test__("(symbol? 11)", FALSE)
        self.__EVAL_test__("(symbol? 11.22)", FALSE)

    def test_bool_functions(self):
        self.__EVAL_test__("(or true false true)", TRUE)
        self.__EVAL_test__("(or false false)", FALSE)
        self.__EVAL_test__("(not true)", FALSE)
        self.__EVAL_test__("(not false)", TRUE)

    def test_list_operations(self):
        self.__EVAL_test__("(cons 1 2 3 '())", [1, 2, 3])
        self.__EVAL_test__("(cons 1 '(2 3))", [1, 2, 3])
        self.__EVAL_test__("(map (lambda (x) (* 2 x)) '(2 3))", [4, 6])
        self.__EVAL_test__("(sorted-by '(-1 1 5 -3 2 -4) abs)", [-1, 1, 2, -3, -4, 5])
        self.__EVAL_test__("(sorted-by '(-1 1 5 -3 2 -4) (lambda (x) x))", [-4, -3, -1, 1, 2, 5])
        self.__EVAL_test__("(len '(-1 1 5 -3 2 -4))", 6)
        self.__EVAL_test__("(car '(-1 1 5 -3 2 -4))", -1)
        self.__EVAL_test__("(cdr '(-1 1 5 -3 2 -4))", [1, 5, -3, 2, -4])
        self.__EVAL_test__("(list -1 1 5 -3 2 -4)", [-1, 1, 5, -3, 2, -4])

    def test_string_functions(self):
        self.__EVAL_test__('(join " " \'("OH" "LOL" "YA"))', "OH LOL YA")
        self.__EVAL_test__('(str \'("OH" "LOL" "YA"))', '("OH" "LOL" "YA")')
        self.__EVAL_test__('(str "OH" "LOL" "YA" 1)', '"OH" "LOL" "YA" 1')

    def test_other_functions(self):
        self.__EVAL_test__("(name 'a)", "a")
        self.__EVAL_test__("(progn 1 2 3)", 3)
        self.__EVAL_test__('(int "234")', 234)
        self.__EVAL_test__('(int 3.14)', 3)
        self.__EVAL_test__('(int 3.14)', 3)
        self.__EVAL_test__('(prompt)', "lis.py> ")

    # def test_label(self):
        # self.__EVAL_test__(
            # """((label subst
                # (lambda (x y z)
                    # (cond
                        # ((atom? z) (cond
                            # ((eq z y) x)
                            # ('t z)))
                        # ('t (cons (subst x y (car z))
                            # (subst x y (cdr z)))))))
                # 'm 'b '(a b (a b c) d))""",
            # [A, Symbol("m"), [A, Symbol("m"), C], Symbol("d")])

    # def test_defun_subst(self):
        # self.__EVAL_test__(
        # """((defun subst (x y z)
            # (cond ((atom? z) (cond ((eq z y) x)
                    # ('t z)))
                # ('t (cons
                    # (subst x y (car z))
                    # (subst x y (cdr z))))))
            # 'm 'b '(a b (a b c) d))""",
            # [A, Symbol("m"), [A, Symbol("m"), C], Symbol("d")])

    # def test_defun_cadr(self):
        # self.__EVAL_test__("""((defun cadr (x) (car (cdr x))) '((a b) (c d) e))""", [C, Symbol("d")])

if __name__ == '__main__':
    unittest.main()
