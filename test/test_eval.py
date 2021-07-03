import unittest

from definitions import A, B, C, NIL
from LiSH.datatypes import Symbol, Hashmap, Keyword
from LiSH.env import default_env
from LiSH.reader import READ
from LiSH.evaluator import EVAL


class TestEVAL(unittest.TestCase):
    def __EVAL_test__(self, test_program, expected_result):
        env = default_env()
        actual_result = EVAL(READ(test_program), env)
        self.assertEqual(actual_result, expected_result)

    def test_quote(self):
        self.__EVAL_test__("(quote a)", A)
        self.__EVAL_test__("(quote (a b c))", [A, B, C])
        self.__EVAL_test__("(quote (a b (c d (f e))))", [A, B, [C, Symbol("d"), [Symbol("f"), Symbol("e")]]])
        self.__EVAL_test__("(quote ())", NIL)

    def test_atom(self):
        self.__EVAL_test__("(atom? (quote a))", True)
        self.__EVAL_test__("(atom? (quote (a b c)))", False)
        self.__EVAL_test__("(atom? (quote ()))", True)
        self.__EVAL_test__("(atom? 'a)", True)
        self.__EVAL_test__("(atom? '(a b c))", False)
        self.__EVAL_test__("(atom? '())", True)
        self.__EVAL_test__("(atom? (atom? 'a))", True)
        self.__EVAL_test__("(atom? '(atom? 'a))", False)

    def test_eq(self):
        self.__EVAL_test__("(= 'a 'a)", True)
        self.__EVAL_test__("(= 'a 'b)", False)
        self.__EVAL_test__("(= '() '())", True)

    def test_cons(self):
        self.__EVAL_test__("(cons 'a '(b c))", [A, B, C])
        self.__EVAL_test__("(cons 'a (cons 'b (cons 'c '())))", [A, B, C])

    def test_car(self):
        self.__EVAL_test__("(car '(a b c))", A)
        self.__EVAL_test__("(car (cons 'a '(b c)))", A)

    def test_cdr(self):
        self.__EVAL_test__("(cdr (cons 'a '(b c)))", [B, C])
        self.__EVAL_test__("(cdr '(a))", NIL)
        self.__EVAL_test__("(cdr '(a b c))", [B, C])

    def test_cond(self):
        self.__EVAL_test__("(cond (= 'a 'a) 'first 'second)", Symbol("first"))
        self.__EVAL_test__("(cond (= 'a 'b) 'first 'second)", Symbol("second"))

    def test_lambda(self):
        self.__EVAL_test__("((fn (x) (cons x '(b))) 'a)", [A, B])
        self.__EVAL_test__("((fn (x y) (cons x (cdr y))) 'z '(a b c))", [Symbol("z"), B, C])
        self.__EVAL_test__("((fn (f) (f '(b c))) (fn (x) (cons 'a x)))", [A, B, C])
        self.__EVAL_test__("((fn (x & y) [x y]) 2)", [2, NIL])
        self.__EVAL_test__("((fn (x & y) [x y]) 2 3)", [2, [3]])
        self.__EVAL_test__("((fn (x & y) [x y]) 2 3 4)", [2, [3, 4]])
        self.__EVAL_test__("((fn (x y & z) [x y z]) 2 3 4)", [2, 3, [4]])
        self.__EVAL_test__("((fn (x & y) (* x (apply + y))) 2 3 4)", 14)

    def test_ariphmetic_operators(self):
        self.__EVAL_test__("(+ 1 2 3)", 6)
        self.__EVAL_test__("(+ 1 2 -3)", 0)
        self.__EVAL_test__("(* 1 2 -3)", -6)
        self.__EVAL_test__("(/ -12 2 3)", -2)

    def test_math_functions(self):
        env = default_env()
        random_result = EVAL(READ("(rand)"), env)
        self.assertTrue(isinstance(random_result, float))
        self.assertTrue(0 <= random_result and random_result < 1)
        self.__EVAL_test__("(abs -2)", 2)
        self.__EVAL_test__("(abs 2)", 2)
        self.__EVAL_test__("(cos 0)", 1.0)
        self.__EVAL_test__("(max 1 2 3)", 3)
        self.__EVAL_test__("(min 1 2 3)", 1)
        self.__EVAL_test__("(round 3.14)", 3)

    def test_comparison_operators(self):
        self.__EVAL_test__("(> 1 2 3)", False)
        self.__EVAL_test__("(> 3 2 1)", True)
        self.__EVAL_test__("(> 3 2 3)", False)
        self.__EVAL_test__("(> 3 2 2 1)", False)
        self.__EVAL_test__("(< 1 2 3)", True)
        self.__EVAL_test__("(< 3 2 1)", False)
        self.__EVAL_test__("(< 3 2 3)", False)
        self.__EVAL_test__("(< 3 2 2 1)", False)
        self.__EVAL_test__("(>= 1 2 3)", False)
        self.__EVAL_test__("(>= 3 2 1)", True)
        self.__EVAL_test__("(>= 3 2 2 1)", True)
        self.__EVAL_test__("(<= 1 2 3)", True)
        self.__EVAL_test__("(<= 3 2 1)", False)
        self.__EVAL_test__("(<= 3 2 3)", False)
        self.__EVAL_test__("(<= 1 2 2 3)", True)
        self.__EVAL_test__("(= 3 2 3)", False)
        self.__EVAL_test__("(= 2 (+ 1 1) (- 4 2))", True)
        self.__EVAL_test__("(nil? '())", True)
        self.__EVAL_test__("(nil? '(1))", False)
        self.__EVAL_test__("(list? '())", True)
        self.__EVAL_test__("(list? '(1))", True)
        self.__EVAL_test__("(number? '())", False)
        self.__EVAL_test__("(number? '(1))", False)
        self.__EVAL_test__("(number? 'a)", False)
        self.__EVAL_test__("(number? 11)", True)
        self.__EVAL_test__("(number? 11.22)", True)
        self.__EVAL_test__("(symbol? '())", False)
        self.__EVAL_test__("(symbol? '(1))", False)
        self.__EVAL_test__("(symbol? 'a)", True)
        self.__EVAL_test__("(symbol? 11)", False)
        self.__EVAL_test__("(symbol? 11.22)", False)

    def test_bool_functions(self):
        self.__EVAL_test__("(or true false true)", True)
        self.__EVAL_test__("(or false false)", False)

    def test_list_operations(self):
        self.__EVAL_test__("(cons 1 2 3 '())", [1, 2, 3])
        self.__EVAL_test__("(cons 1 '(2 3))", [1, 2, 3])
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
        self.__EVAL_test__("(int 3.14)", 3)
        self.__EVAL_test__("(int 3.14)", 3)
        self.__EVAL_test__("(prompt)", "lis.py> ")

    def test_let(self):
        self.__EVAL_test__("(let* (x 1 y 2) x)", 1)
        self.__EVAL_test__("(str (let* (x 1 y 2) x))", "1")
        self.__EVAL_test__("(str (let* (x 1 y 2) y))", "2")

    def test_vector(self):
        self.__EVAL_test__("[]", [])
        self.__EVAL_test__("[1 2 (+ 1 2)]", [1, 2, 3])

    def test_hashmap(self):
        self.__EVAL_test__('{}', Hashmap([]))
        self.__EVAL_test__('{"a" (+ 1 2)}', Hashmap(["a", 3]))
        self.__EVAL_test__('{:a (+ 1 2)}', Hashmap([Keyword("a"), 3]))
        self.__EVAL_test__('{(+ 1 2) (+ 1 2)}', Hashmap([3, 3]))

    def test_not_defined_function(self):
        with self.assertRaises(RuntimeError) as cm:
            self.__EVAL_test__("(abc 1 2 3)", None)
        self.assertEqual(str(cm.exception), "abc value not found")

    def test_sort(self):
        self.__EVAL_test__("(sort [9 6 5 8 7 6 3])", [3, 5, 6, 6, 7, 8, 9])
        self.__EVAL_test__("(apply < (sort [9 6 5 8 7 3]))", True)
        self.__EVAL_test__("(sort '(-1 1 5 -3 2 -4) abs)", [-1, 1, 2, -3, -4, 5])
        self.__EVAL_test__("(sort '(-1 1 5 -3 2 -4) (fn (x) x))", [-4, -3, -1, 1, 2, 5])


if __name__ == '__main__':
    unittest.main()
