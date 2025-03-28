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













new_coroutine_queue := (init) => (
  ctx := {queue: map(init, coro => () => coro(ctx))}
  ctx.spawn := thunk => ctx.queue.len(ctx.queue) = thunk
  resume_next_coroutine := () => true :: {
    len(ctx.queue) != 0 -> (
      next := ctx.queue.0 # TODO: pop
      ctx.queue = ctx.queue.tail
      reset(next)
    )
  }
  ctx.yield := () => shift(k => (
    ctx.spawn(() => k(void))
    resume_next_coroutine()
  ))
  ctx.run := () => resume_next_coroutine()
  ctx
)

# Example coroutines
coroutine1 := (ctx) => (
  out("Coroutine 1: Step 1\n")
  ctx.yield()
  out("Coroutine 1: Step 2\n")
  ctx.yield()
  out("Coroutine 1: Step 3\n")
)

coroutine2 := (ctx) => (
  out("Coroutine 2: Step 1\n")
  ctx.yield()
  out("Coroutine 2: Step 2\n")
  ctx.yield()
  out("Coroutine 2: Step 3\n")
)

new_coroutine_queue([coroutine1, coroutine2]).run()

handler_coroutine := id => ctx => (
  out(f("Handler %: Processing request...\n", id))
  ctx.yield() # Simulate async I/O
  out(f("Handler %: Sending response...\n", id))
  ctx.yield() # Simulate async I/O
  out(f("Handler %: Done.\n", id))
)

web_server := (ctx) => (loop := () => (
  out("Web Server: Waiting for connection...\n")
  ctx.yield() # Simulate waiting for I/O
  out("Web Server: Connection accepted. Spawning handler.\n")
  ctx.spawn(handler_coroutine(random(1000)))
  loop()
)())

new_coroutine_queue([web_server]).run() # Start the server
