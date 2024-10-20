# application utility tests

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/868f5e85c16fc7dbff6e630dd1e90d365a2c5546/std.ink')
f := std.format
each := std.each

formatNumber := import('../lib/util').formatNumber

run := (m, t) => (
  m('formatNumber')
  TestVals := [
    # normal cases
    [0, '0']
    [3, '3']
    [27, '27']
    [100, '100']
    [123, '123']
    [7331, '7,331']
    [14243, '14,243']
    [153243, '153,243']
    [8765432, '8,765,432']
    [87654321, '87,654,321']

    # regression tests
    [1007, '1,007']
    [1023, '1,023']
    [10234, '10,234']
    [10034, '10,034']
    [9000000, '9,000,000']
  ]

  each(TestVals, pair => (
    num := pair.0
    result := pair.1

    t(f('correctly formats {{ 0 }}', [num])
      formatNumber(num), result)
  ))
)
