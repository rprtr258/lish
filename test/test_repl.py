import unittest

from definitions import NIL, A, B, C
from context import LispSH
from LispSH import parse, eval, default_env, atom


class TestRepl(unittest.TestCase):
    def test_cadr(self):
        global_env = default_env()
        self.assertEqual(eval(parse("(define cadr (lambda (x) (car (cdr x))))"), global_env), NIL)
        self.assertEqual(eval(parse("(cadr '(a (b c) d))"), global_env), [B, C])

    def test_if_macro_true(self):
        global_env = default_env()
        self.assertEqual(eval(parse("(defmacro if (p x y) (cond p x y))"), global_env), NIL)
        self.assertEqual(eval(parse("(if (nil? '()) 1 2)"), global_env), atom(1))

    def test_if_macro_false(self):
        global_env = default_env()
        self.assertEqual(eval(parse("(defmacro if (p x y) (cond p x y))"), global_env), NIL)
        self.assertEqual(eval(parse("(if (nil? '(a)) 1 2)"), global_env), atom(2))

    def test_if_in_recursive_function(self):
        global_env = default_env()
        self.assertEqual(eval(parse("(defmacro if (p x y) (cond p x y))"), global_env), NIL)
        self.assertEqual(eval(parse("(define fact (lambda (n) (cond (eq? n 1) 1 (* n (fact (- n 1))))))"), global_env), NIL)
        self.assertEqual(eval(parse("(fact 1)"), global_env), atom(1))
        self.assertEqual(eval(parse("(fact 2)"), global_env), atom(2))
        self.assertEqual(eval(parse("(fact 3)"), global_env), atom(6))
        self.assertEqual(eval(parse("(fact 4)"), global_env), atom(24))
        self.assertEqual(eval(parse("(define rev (lambda (x) (progn (echo x) (if (nil? x) x (+ (rev (cdr x)) (list (car x)))))))"), global_env), NIL)
        self.assertEqual(eval(parse("(rev '(a b c))"), global_env), [C, B, A])

    def test_self_combinator(self):
        global_env = default_env()
        self.assertEqual(eval(parse("(define S (lambda (y) (y y)))"), global_env), NIL)
        self.assertEqual(eval(parse("(S str)"), global_env), "<fun str>")

    def test_recursion(self):
        global_env = default_env()
        # TODO: replace eq? with =
        self.assertEqual(eval(parse("(define fact (lambda (n) (cond (eq? n 1) 1 (* n (fact (- n 1))))))"), global_env), NIL)
        self.assertEqual(eval(parse("(fact 1)"), global_env), atom(1))
        self.assertEqual(eval(parse("(fact 2)"), global_env), atom(2))
        self.assertEqual(eval(parse("(fact 3)"), global_env), atom(6))
        self.assertEqual(eval(parse("(fact 4)"), global_env), atom(24))
        
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
