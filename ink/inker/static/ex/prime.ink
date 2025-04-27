# prime sieve

{log, filter, reduce} := import('std.ink')

# is a single number prime?
isPrime := n => (
  # is n coprime with nums < p?
  max := floor(pow(n, 0.5)) + 1
  (ip := p => p == max | ~(n % p == 0) & ip(p + 1))(2) # start with smaller = more efficient
)

# build a list of consecutive integers from 2 .. max
buildConsecutive := max => (
  peak := max + 1
  acc := []
  (bc := i => i :: {
    peak -> ()
    _ -> (
      acc.(i - 2) := i
      bc(i + 1)
    )
  })(2)
  acc
)

# primes under N are numbers 2 .. N, filtered by isPrime
getPrimesUnder := n => filter(buildConsecutive(n), (n, _) => isPrime(n))

ps := getPrimesUnder(10000)
log(ps)
log('Total number of primes under 10000: ' + string(len(ps)))
