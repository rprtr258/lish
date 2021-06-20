import unittest

from definitions import NIL, A, B, C
from context import LispSH
from LispSH import parse, eval, default_env, atom


class TestRepl(unittest.TestCase):
    def __test_cmds__(self, cmds):
        env = default_env()
        for cmd_line in cmds:
            if isinstance(cmd_line, str):
                cmd, expected_result = cmd_line, NIL
            elif isinstance(cmd_line, tuple) and len(cmd_line) == 2:
                cmd, expected_result = cmd_line
            else:
                self.assertTrue(False)
            actual_result = eval(parse(cmd), env)
            self.assertEqual(actual_result, expected_result)

    def test_cadr(self):
        self.__test_cmds__([
            "(define cadr (lambda (x) (car (cdr x))))",
            ("(cadr '(a (b c) d))", [B, C])
        ])

    def test_if_macro(self):
        self.__test_cmds__([
            "(defmacro if (p x y) (list 'cond p x y))",
            ("(if (nil? '()) 1 2)", atom(1)),
            ("(if (nil? '(a)) 1 2)", atom(2))
        ])

    def test_defun_macro(self):
        self.__test_cmds__([
            "(defmacro if (p x y) (list 'cond p x y))",
            "(defmacro defun (f args body) (list 'define f (list 'lambda args body)))",
            "(defun rev (x) (if (nil? x) x (+ (rev (cdr x)) (list (car x)))))",
            ("(rev '(a b c))", [C, B, A])
        ])

    def test_if_in_recursive_function(self):
        self.__test_cmds__([
            "(defmacro if (p x y) (list 'cond p x y))",
            "(define fact (lambda (n) (cond (eq? n 1) 1 (* n (fact (- n 1))))))",
            ("(fact 1)", atom(1)),
            ("(fact 2)", atom(2)),
            ("(fact 3)", atom(6)),
            ("(fact 4)", atom(24)),
            "(define rev (lambda (x) (progn (echo x) (if (nil? x) x (+ (rev (cdr x)) (list (car x)))))))",
            ("(rev '(a b c))", [C, B, A])
        ])

    def test_self_combinator(self):
        self.__test_cmds__([
            "(define S (lambda (y) (y y)))",
            ("(S str)", "<fun str>")
        ])

    def test_recursion(self):
        self.__test_cmds__([
            "(define fact (lambda (n) (cond (eq? n 1) 1 (* n (fact (- n 1))))))",
            ("(fact 1)", atom(1)),
            ("(fact 2)", atom(2)),
            ("(fact 3)", atom(6)),
            ("(fact 4)", atom(24))
        ])
        
    # def test_y_combinator(self):
        # global_env = default_env()
        # self.assertEqual(eval(parse("(define S (lambda (y) (y y)))"), global_env), NIL)
        # self.assertEqual(eval(parse("(define Y (lambda (f) (S (lambda (z) (f (z z))))))"), global_env), NIL)
        # self.assertEqual(eval(parse("(define fact (lambda (f) (lambda (x) (if (eq? x 1) 1 (* x (f (- x 1)))))))"), global_env), NIL)
        # self.assertEqual(eval(parse("(define yfact (Y fact))"), global_env), NIL)
        # self.assertEqual(eval(parse("(yfact 1)"), global_env), 1)
        # self.assertEqual(eval(parse("(yfact 2)"), global_env), 2)
        # self.assertEqual(eval(parse("(yfact 3)"), global_env), 6)
        # self.assertEqual(eval(parse("(yfact 4)"), global_env), 24)

    # (define p3f (lambda (x) (progn (echo x) (echo x) (echo x))))
    # (defmacro p3m (x) (progn (echo x) (echo x) (echo x)))

if __name__ == '__main__':
    unittest.main()
