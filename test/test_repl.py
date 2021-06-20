import unittest

from definitions import NIL, B, C
from context import LispSH
from LispSH import parse, eval, default_env


class TestRepl(unittest.TestCase):
    def test_define_use_cadr(self):
        global_env = default_env()
        self.assertEqual(eval(parse("(define cadr (lambda (x) (car (cdr x))))"), global_env), NIL)
        self.assertEqual(eval(parse("(cadr '(a (b c) d))"), global_env), [B, C])

    def test_define_use_if_macro_true(self):
        global_env = default_env()
        self.assertEqual(eval(parse("(defmacro if (p x y) (cond p x y))"), global_env), NIL)
        self.assertEqual(eval(parse("(if (nil? '()) 1 2)"), global_env), 1)

    def test_define_use_if_macro_false(self):
        global_env = default_env()
        self.assertEqual(eval(parse("(defmacro if (p x y) (cond p x y))"), global_env), NIL)
        self.assertEqual(eval(parse("(if (nil? '(a)) 1 2)"), global_env), 2)

    def test_define_use_self_combinator(self):
        global_env = default_env()
        self.assertEqual(eval(parse("(define S (lambda (y) (y y)))"), global_env), NIL)
        self.assertEqual(eval(parse("(S echo)"), global_env), "<fun echo>")
        
    def test_define_use_y_combinator(self):
        global_env = default_env()
        self.assertEqual(eval(parse("(define S (lambda (y) (y y)))"), global_env), NIL)
        self.assertEqual(eval(parse("(define Y (lambda (f) (S (lambda (z) (f (z z))))))"), global_env), NIL)
        self.assertEqual(eval(parse("(define fact (lambda (f) (lambda (x) (if (eq? x 1) 1 (* x (f (- x 1)))))))"), global_env), NIL)
        self.assertEqual(eval(parse("(define yfact (Y fact))"), global_env), NIL)
        self.assertEqual(eval(parse("(yfact 1)"), global_env), 1)
        self.assertEqual(eval(parse("(yfact 2)"), global_env), 2)
        self.assertEqual(eval(parse("(yfact 3)"), global_env), 6)
        self.assertEqual(eval(parse("(yfact 4)"), global_env), 24)

if __name__ == '__main__':
    unittest.main()
