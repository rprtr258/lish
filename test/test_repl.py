import unittest
from contextlib import redirect_stdout
from io import StringIO

from definitions import NIL, A, B, C
from context import LiSH
from LiSH.reader import READ
from LiSH.evaluator import EVAL
from LiSH.env import default_env
from LiSH.datatypes import Symbol


class TestRepl(unittest.TestCase):
    def __create_env__(self):
        env = default_env()
        env[Symbol("eval")] = lambda ast: EVAL(ast, env)
        EVAL(READ('(set! load-file (lambda (f) (eval (read (+ "(progn " (slurp f) "\n)")))))'), env)
        EVAL(READ('(load-file "compose.lish")'), env)
        return env

    def __test_cmds__(self, cmds):
        env = self.__create_env__()
        for cmd_line in cmds:
            if isinstance(cmd_line, str):
                EVAL(READ(cmd_line), env)
                continue
            elif isinstance(cmd_line, tuple) and len(cmd_line) == 2:
                cmd, expected_result = cmd_line
            else:
                self.assertTrue(False)
            actual_result = EVAL(READ(cmd), env)
            self.assertEqual(actual_result, expected_result)

    def __test_cmds_output__(self, cmds, expected_output):
        env = self.__create_env__()
        with redirect_stdout(StringIO()) as f:
            for cmd_line in cmds:
                EVAL(READ(cmd_line), env)
        self.assertEqual(f.getvalue(), expected_output)

    def test_cadr(self):
        self.__test_cmds__([
            ("(cadr '(a (b c) d))", [B, C])
        ])

    def test_if_macro(self):
        self.__test_cmds__([
            ("(if (nil? '()) 1 2)", 1),
            ("(if (nil? '(a)) 1 2)", 2)
        ])

    def test_letfun_macro(self):
        self.__test_cmds_output__([
            "(echo (letfun (f (x) (+ 2 x) g (x) (* 3 x)) (-> 4 g f)))"
        ], "14\n")

    def test_defun_macro(self):
        self.__test_cmds__([
            "(defun rev (x) (if (nil? x) x (+ (rev (cdr x)) (list (car x)))))",
            ("(rev '(a b c))", [C, B, A])
        ])

    def test_if_in_recursive_function(self):
        self.__test_cmds__([
            "(defun fact (n) (if (= n 1) 1 (* n (fact (- n 1)))))",
            ("(fact 1)", 1),
            ("(fact 2)", 2),
            ("(fact 3)", 6),
            ("(fact 4)", 24),
            "(defun rev (x) (if (nil? x) x (+ (rev (cdr x)) (list (car x)))))",
            ("(rev '(a b c))", [C, B, A])
        ])

    def test_self_combinator(self):
        self.__test_cmds__([
            "(set! S (lambda (y) (y y)))",
            ("(S str)", "#str")
        ])

    def test_recursion(self):
        self.__test_cmds__([
            "(defun fact (n) (if (= n 1) 1 (* n (fact (- n 1)))))",
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
        env = self.__create_env__()
        EVAL(READ("(defun p3f (x) (list x x x))"), env)
        result = EVAL(READ("(p3f (rand))"), env)
        self.assertEqual(*result)

    def test_triple_print_macro(self):
        env = self.__create_env__()
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
    
    def test_print_elochka(self):
        self.__test_cmds_output__([
            """(echo "el0chka")""",
            """(doseq (x '(1 1 1) y (range 1 9 2))
                (echo (* "*" y)))"""
        ], """el0chka
*
***
*****
*******
*
***
*****
*******
*
***
*****
*******
""")

    def test_print_almaz(self):
        self.__test_cmds_output__([
            """(echo "almaz")""",
            """(doseq (x (+ (range 1 9 2) (range 5 0 -2)))
                (echo (* " " (- 7 x)) (* "*" (* x 2))))"""
        ], """almaz
      **
    ******
  **********
**************
  **********
    ******
      **
""")

    def test_prompt_counter(self):
        self.__test_cmds_output__([
            """(set! prompt
                (let* (cnt 0)
                    (lambda () (progn
                    (swap! cnt inc)
                    (+ "lis.py(" (str cnt) ")> 123")))))""",
            "(echo (prompt))",
            "(echo (prompt))",
            "(echo (prompt))"
        ], """lis.py(1)> 123
lis.py(2)> 123
lis.py(3)> 123
""")

    def test_prompt_counter(self):
        self.__test_cmds_output__([
            """(doseq (x '(-1 0 1)
                       y '(1 2 3))
                (echo (* x y)))"""
        ], """-1
-2
-3
0
0
0
1
2
3
""")


if __name__ == '__main__':
    unittest.main()
