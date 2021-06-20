program = [
    "(defun cadr (x) (car (cdr x)))",
    "(cadr '(a (b c) d))"
]
for cmd in program:
    res, env = interpret(cmd, env)
    print(res)