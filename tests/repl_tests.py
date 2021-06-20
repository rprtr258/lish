program = [
    "(defun cadr (x) (car (cdr x)))",
    "(cadr '(a (b c) d))"
]
for cmd in program:
    res = eval(parse(cmd))
    print(res)
