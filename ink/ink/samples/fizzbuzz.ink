# ink fizzbuzz implementation

log := s => out(string(s) + '\n')
std := import('functional')
range := std.range
each := std.each

fizzbuzz := n => each(
  range(1, n + 1, 1)
  n => log([n % 3, n % 5] :: {
    [0, 0] -> 'FizzBuzz'
    [0, _] -> 'Fizz'
    [_, 0] -> 'Buzz'
    _ -> n
  })
)

fizzbuzz(100)
