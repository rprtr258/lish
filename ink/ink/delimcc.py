from typing import Callable, Any
from dataclasses import dataclass

# Expression classes
class Expr:
    pass

@dataclass
class Number(Expr):
    value: float

@dataclass
class BinOp(Expr):
    op: str
    left: Expr
    right: Expr

@dataclass
class Shift(Expr):
    func: Callable

@dataclass
class Reset(Expr):
    body: Expr

# Trampoline classes
class Trampoline:
    pass

@dataclass
class Done(Trampoline):
    result: Any

@dataclass
class Call(Trampoline):
    fn: Callable
    args: tuple

# Continuation for outer TCO
class Continuation:
    def __init__(self, k: Callable):
        self.k = k

    def __call__(self, value: Any) -> Trampoline:
        return self.k(value)

# Evaluation inside reset (no TCO, returns value directly)
def eval_inside_reset(expr: Expr, k: Callable) -> Any:
    if isinstance(expr, Number):
        return k(expr.value)
    elif isinstance(expr, BinOp):
        l = eval_inside_reset(expr.left, lambda x: x)
        r = eval_inside_reset(expr.right, lambda x: x)
        if expr.op == '*':
            return k(l * r)
        # Add other operators as needed
    elif isinstance(expr, Shift):
        return expr.func(k)
    elif isinstance(expr, Reset):
        val = eval_inside_reset(expr.body, lambda x: x)
        return k(val)
    raise ValueError(f"Unknown expression: {expr}")

def eval_expr(e: Expr, k: Callable) -> Trampoline:
    if isinstance(e, Number):
        return Call(k, (e.value,))
    elif isinstance(e, BinOp):
        def op(left_val, right_val):
            if e.op == '*':
                return Call(k, (left_val * right_val,))
            # Add other operators
        return eval_expr(e.left, lambda l: lambda r: op(l, r))
    elif isinstance(e, Reset):
        val = eval_inside_reset(e.body, lambda x: x)
        return Call(k, (val,))
    elif isinstance(e, Shift):
        return Call(e.func, (Continuation(k),))
    raise ValueError(f"Unknown expression: {e}")

# Main evaluation with TCO
def evaluate(expr: Expr) -> Any:
    result = eval_expr(expr, Done)
    while isinstance(result, Call):
        fn, args = result.fn, result.args
        result = fn(*args)
    return result.result

# Helpers
def mul(left: Expr, right: Expr) -> BinOp:
    return BinOp('*', left, right)

if __name__ == "__main__":
    nested_expr = Reset(mul(Number(2), Shift(lambda k: k(k(5)))))
    result = evaluate(nested_expr)
    print(f"reset(2 * shift(lambda k: k(k(5)))) = {result}")
