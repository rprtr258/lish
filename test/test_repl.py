import unittest

from definitions import NIL, A, B, C
from context import LispSH
from LispSH import parse, eval, default_env, atom, symbol


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
            "(define fact (lambda (n) (cond (= n 1) 1 (* n (fact (- n 1))))))",
            ("(fact 1)", atom(1)),
            ("(fact 2)", atom(2)),
            ("(fact 3)", atom(6)),
            ("(fact 4)", atom(24)),
            "(define rev (lambda (x) (if (nil? x) x (+ (rev (cdr x)) (list (car x))))))",
            ("(rev '(a b c))", [C, B, A])
        ])

    def test_self_combinator(self):
        self.__test_cmds__([
            "(define S (lambda (y) (y y)))",
            ("(S str)", atom("<fun str>"))
        ])

    def test_recursion(self):
        self.__test_cmds__([
            "(define fact (lambda (n) (cond (= n 1) 1 (* n (fact (- n 1))))))",
            ("(fact 1)", atom(1)),
            ("(fact 2)", atom(2)),
            ("(fact 3)", atom(6)),
            ("(fact 4)", atom(24))
        ])

    def test_let(self):
        self.__test_cmds__([
            "(defmacro defun (f args body) (list 'define f (list 'lambda args body)))",
            "(defun evens (x) (cond (nil? x) '() (cons (car x) (odds (cdr x)))))",
            "(defun odds (x) (cond (nil? x) '() (evens (cdr x))))",
            "(defmacro let (exps body) (cons (list 'lambda (evens exps) body) (odds exps)))",
            ("(evens '(x 1 y 2))", [symbol("x"), symbol("y")]),
            ("(odds '(x 1 y 2))", [atom(1), atom(2)]),
            ("(let (x 1 y 2) x)", atom(1)),
            ("(str (let (x 1 y 2) x))", atom("1")),
            ("(str (let (x 1 y 2) y))", atom("2"))
        ])

    def test_triple_print_function(self):
        env = default_env()
        eval(parse("(defmacro defun (f args body) (list 'define f (list 'lambda args body)))"), env)
        eval(parse("(defun p3f (x) (list x x x))"), env)
        result = eval(parse("(p3f (rand))"), env)
        self.assertEqual(*result)

    def test_triple_print_macro(self):
        env = default_env()
        eval(parse("(defmacro p3m (x) (list 'list x x x))"), env)
        result = eval(parse("(p3m (rand))"), env)
        self.assertNotEqual(result[0], result[1])
        self.assertNotEqual(result[1], result[2])
        self.assertNotEqual(result[0], result[2])

    # def test_y_combinator(self):
        # global_env = default_env()
        # self.assertEqual(eval(parse("(define S (lambda (y) (y y)))"), global_env), NIL)
        # self.assertEqual(eval(parse("(define Y (lambda (f) (S (lambda (z) (f (z z))))))"), global_env), NIL)
        # self.assertEqual(eval(parse("(define fact (lambda (f) (lambda (x) (if (= x 1) 1 (* x (f (- x 1)))))))"), global_env), NIL)
        # self.assertEqual(eval(parse("(define yfact (Y fact))"), global_env), NIL)
        # self.assertEqual(eval(parse("(yfact 1)"), global_env), 1)
        # self.assertEqual(eval(parse("(yfact 2)"), global_env), 2)
        # self.assertEqual(eval(parse("(yfact 3)"), global_env), 6)
        # self.assertEqual(eval(parse("(yfact 4)"), global_env), 24)

if __name__ == '__main__':
    unittest.main()
