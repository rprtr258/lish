# prime sieve

{log} := import('logging.ink')
{filter} := import('functional.ink')
{range} := import('functional.ink')

# is a single number prime?
isPrime := n =>
  # is n coprime with nums < p?
  (ip := p => (p*p > n) | ~(n%p == 0) & ip(p + 1))(2) # start with smaller # == more efficient

N := 5000
# primes under N are numbers 2 .. N, filtered by isPrime
ps := filter(range(2, N+1, 1), isPrime)
log(string(ps))
log('Total number of primes under 5000: ' + string(len(ps)))
