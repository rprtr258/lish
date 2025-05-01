# prime sieve
# taken from thesephist/ink/samples/

{log, filter, stringList} := import('./runtime/std')
# {log, filter, stringList} := import('../../vendor/std')

# is a single number prime?
isPrime := n => (
  # is n coprime with nums < p?
  max := floor(pow(n, 0.5)) + 1
  (ip := p => p :: {
    max -> true
    _ -> n % p :: {
      0 -> false
      _ -> ip(p + 1)
    }
  })(2) # start with smaller # = more efficient
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
getPrimesUnder := n => filter(buildConsecutive(n), isPrime)

ps := getPrimesUnder(1250)
log(stringList(ps))
log('Total number of primes under 1250: ' + string(len(ps)))
