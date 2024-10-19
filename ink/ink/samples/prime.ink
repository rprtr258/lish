`` prime sieve

log := load('logging').log
filter := load('functional').filter
range := load('functional').range

`` is a single number prime?
isPrime := n =>
  `` is n coprime with nums < p?
  (ip := p => true :: {
    p*p > n -> true
    n%p = 0 -> false
    _ -> ip(p + 1)
  })(2) `` start with smaller # = more efficient

N := 5000
`` primes under N are numbers 2 .. N, filtered by isPrime
ps := filter(range(2, N+1, 1), isPrime)
log(string(ps))
log('Total number of primes under 5000: ' + string(len(ps)))
