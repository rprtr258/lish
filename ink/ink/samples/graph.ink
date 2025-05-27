# let's graph the sine / cosine functions in Ink!
{join} := import('str.ink')

log := (s) => out(s + '\n')

# loop a function F, N times
loop := (n, f) => true :: {
  n > 0 -> (
    f(n)
    loop(n - 1, f)
  )
}

# repeat a string n times
repeat := (s, n) => (
  res := []
  loop(n, (_) => res.len(res) := s)
  join(res, '')
)

# graph a single point
draw := (row, x, func, symbol) => (
  n := func(x / 4) + 1
  # some fuzzy math here to make spacing look decent
  row.(floor(20 * n)) := symbol
)

# recursively draw from a single value
drawRec := (max) => (
  loop(max, (n) => (
    row := repeat(' ', 60)
    draw(row, n, sin, '+')
    draw(row, n, (x) => cos(x) + 0.7, 'o')
    log(row)
  ))
)

# actually draw
drawRec(40)
