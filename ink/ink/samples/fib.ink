` fibonacci sequence generator `

log := s => out(s + char(10))

` naive implementation `
fib := n => n :: {
  0 -> 0
  1 -> 1
  _ -> fib(n - 1) + fib(n - 2)
}

` memoized / dynamic programming implementation `
memo := [0, 1]
fibMemo := n => (
  memo.(n) :: {
    () -> memo.(n) := fibMemo(n - 1) + fibMemo(n - 2)
  }
  memo.(n)
)

log('fib(20) is 6765:')
log('Naive solution: ' + string(fib(20)))
log('Dynamic solution: ' + string(fibMemo(20)))
