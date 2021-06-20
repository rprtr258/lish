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

if __name__ == '__main__':
    unittest.main()
