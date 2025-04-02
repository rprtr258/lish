Expr :: +{
  Number = {
    value = float
  },
  BinOp = {
    op = '*' | '+'
    left = Expr
    right: Expr
  },
  Shift = {
    func = ((Any) -> Any) -> Any
  },
  Reset = {
    body = Expr
  },
}

Trampoline :: +{
  Done = {
      result = Any
  }
  Call = {
      fn = () => Trampoline
  }
}

# Continuation for outer TCO
Continuation :: {
  k = Callable
  __call__ :: (self, value: Any): Trampoline => self.k(value)
}

# Evaluation inside reset (no TCO, returns value directly)
eval_inside_reset :: (expr: Expr, k: (Any) -> Any): Any => expr :: {
  Number -> k(expr.value)
  BinOp(op, left, right) -> (
      l := eval_inside_reset(left, x => x)
      r := eval_inside_reset(right, x => x)
      expr.op :: {
        '*' -> k(l * r)
      }
  )
  Shift(func) -> func(k)
  Reset(body) -> (
    val := eval_inside_reset(body, x => x)
    k(val)
  )
}

eval_expr :: (e: Expr, k: Callable): Trampoline => e :: {
  Number(value) -> .Call(() => k(e.value))
  BinOp(op, left, right) -> (
    op := (left_val, right_val) => op :: {
      '*' -> .Call(() => k(left_val * right_val))
    }
    eval_expr(left, l => r => op(l, r))
  )
  Reset(body) -> (
      val := eval_inside_reset(body, x => x)
      .Call(() => k(val))
  )
  Shift(func) -> .Call(() => func(Continuation(k)))
}

# Main evaluation with TCO
evaluate :: (expr: Expr): Any => (
  result := eval_expr(expr, Done)
  # TODO: no syntax for loops yet
  while result :: {Call(fn) -> (
    result = fn()
  )}
  result.result # result:Done here
)

# Helpers
mul :: (left: Expr, right: Expr): Expr => .BinOp('*', left, right)

nested_expr := Expr.Reset(mul(.Number(2), .Shift(k => k(k(5)))))
result := evaluate(nested_expr)
print("reset(2 * shift(lambda k: k(k(5)))) = %" % [result])

