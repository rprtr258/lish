import unittest

from definitions import NIL, A, B, C
from context import LispSH
from LispSH.reader import READ
from LispSH.evaluator import EVAL
from LispSH.env import default_env
from LispSH.datatypes import Symbol


class TestRepl(unittest.TestCase):
    def __test_cmds__(self, cmds):
        env = default_env()
        for cmd_line in cmds:
            if isinstance(cmd_line, str):
                EVAL(READ(cmd_line), env)
                return
            elif isinstance(cmd_line, tuple) and len(cmd_line) == 2:
                cmd, expected_result = cmd_line
            else:
                self.assertTrue(False)
            actual_result = EVAL(READ(cmd), env)
            self.assertEqual(actual_result, expected_result)

    def test_cadr(self):
        self.__test_cmds__([
            "(set! cadr (lambda (x) (car (cdr x))))",
            ("(cadr '(a (b c) d))", [B, C])
        ])

    def test_if_macro(self):
        self.__test_cmds__([
            "(defmacro if (p x y) (list 'cond p x y))",
            ("(if (nil? '()) 1 2)", 1),
            ("(if (nil? '(a)) 1 2)", 2)
        ])

    def test_defun_macro(self):
        self.__test_cmds__([
            "(defmacro if (p x y) (list 'cond p x y))",
            "(defmacro defun (f args body) (list 'set! f (list 'lambda args body)))",
            "(defun rev (x) (if (nil? x) x (+ (rev (cdr x)) (list (car x)))))",
            ("(rev '(a b c))", [C, B, A])
        ])

    def test_if_in_recursive_function(self):
        self.__test_cmds__([
            "(defmacro if (p x y) (list 'cond p x y))",
            "(set! fact (lambda (n) (cond (= n 1) 1 (* n (fact (- n 1))))))",
            ("(fact 1)", 1),
            ("(fact 2)", 2),
            ("(fact 3)", 6),
            ("(fact 4)", 24),
            "(set! rev (lambda (x) (if (nil? x) x (+ (rev (cdr x)) (list (car x))))))",
            ("(rev '(a b c))", [C, B, A])
        ])

    def test_self_combinator(self):
        self.__test_cmds__([
            "(set! S (lambda (y) (y y)))",
            ("(S str)", "<fun str>")
        ])

    def test_recursion(self):
        self.__test_cmds__([
            "(set! fact (lambda (n) (cond (= n 1) 1 (* n (fact (- n 1))))))",
            ("(fact 1)", 1),
            ("(fact 2)", 2),
            ("(fact 3)", 6),
            ("(fact 4)", 24)
        ])

    def test_anaphoric_lambda(self):
        self.__test_cmds__([
            "(defmacro # (& body) (list 'lambda '(%) (cons 'progn body)))",
            ("((# (+ % 2)) 3)", 5)
        ])

    def test_triple_print_function(self):
        env = default_env()
        EVAL(READ("(defmacro defun (f args body) (list 'set! f (list 'lambda args body)))"), env)
        EVAL(READ("(defun p3f (x) (list x x x))"), env)
        result = EVAL(READ("(p3f (rand))"), env)
        self.assertEqual(*result)

    def test_triple_print_macro(self):
        env = default_env()
        EVAL(READ("(defmacro p3m (x) (list 'list x x x))"), env)
        result = EVAL(READ("(p3m (rand))"), env)
        self.assertNotEqual(result[0], result[1])
        self.assertNotEqual(result[1], result[2])
        self.assertNotEqual(result[0], result[2])

    def test_rev_macro(self):
        self.__test_cmds__([
            "(defmacro rev (x) ((defun rev-helper (x) (cond (nil? x) x (+ (rev-helper (cdr x)) (list (car x))))) x))",
            ("(rev (1 str))", "1")
        ])


if __name__ == '__main__':
    unittest.main()
