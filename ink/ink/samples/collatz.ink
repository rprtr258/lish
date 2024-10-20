# finding long collatz sequences

log := import('logging.ink').log
f := import('str.ink').format
pipe := import('functional.ink').pipe
iter := import('iter.ink')
range := iter.range
map := iter.map
fold := iter.foldl
generate := iter.generate
takeWhile := iter.takeWhile
collect := iter.collect
ternary := import('math.ink').ternary

sequence := start => pipe(start, [
  n => generate(n, n => n % 2 :: {
    0 -> n / 2
    1 -> 3 * n + 1
  })
  it => takeWhile(it, n => n > 1)
  collect
])

longestSequenceUnder := max => pipe(max, [
  n => range(1, n+1, 1),
  it => map(it, sequence)
  it => fold(it, (acc, x) => ternary(len(x) < len(acc), acc, x), [])
])

# run a search for longest collatz sequence under Max
Max := 1000
longest := longestSequenceUnder(Max)
log(f('Longest collatz seq under {{ max }} is {{ len }} items', {
  max: Max
  len: len(longest)
}))
log(string(longest))
