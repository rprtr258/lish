# ink fizzbuzz implementation

{log, range, each} := import('std.ink')

fizzbuzz := n => each(
  range(1, n + 1, 1)
  n => [n % 3, n % 5] :: {
    [0, 0] -> log('FizzBuzz')
    [0, _] -> log('Fizz')
    [_, 0] -> log('Buzz')
    _ -> log(n)
  }
)

fizzbuzz(100)
