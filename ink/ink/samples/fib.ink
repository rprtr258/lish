# fibonacci sequence generator

log := s => out(s + '\n')

# naive implementation
fib := n => true :: {
  n < 2 -> n
  _     -> fib(n - 1) + fib(n - 2)
}

# memoized / dynamic programming implementation
#memo := [0, 1]
#fibMemo := n => (
#  memo.(n) :: {
#    () -> memo.(n) := fibMemo(n - 1) + fibMemo(n - 2)
#  }
#  memo.(n)
#)

N := 3
log('fib('+string(N)+') is:')
log('Naive solution: ' + string(fib(N)))
#log('Dynamic solution: ' + string(fibMemo(N)))
