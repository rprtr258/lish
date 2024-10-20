sqrt := x => pow(x, 0.5)
# pow := pow
E := 2.718281828459045
exp := x => pow(E, x)

# ln := ln
log := (x, base) => ln(x) / ln(base)

# sin := sin
# cos := cos
tan := x => sin(x) / cos(x)

# asin := asin
# acos := acos

ternary := (p, x, y) => p :: {
  true -> x
  _ -> y
}
min := (x, y) => ternary(x < y, x, y)
max := (x, y) => ternary(x > y, x, y)
iverson := p => ternary(p, 1, 0)
abs := x => x * ternary(x < 0, ~1, 1)

fract := x => x - floor(x)
ceil := x => fract(x) :: {
  0 -> x
  _ -> fract(x + 1)
}
# floor := floor
round := x => floor(x + 0.5)
trunc := x => ternary(x > 0, floor, ceil)(x)

{
  sqrt
  E
  exp
  log
  tan
  ternary
  min
  max
  iverson
  abs
  fract
  ceil
  round
  trunc
}