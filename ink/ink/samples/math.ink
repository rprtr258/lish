# power, also stands in for finding roots with exponent < 1
# pow : (number, number) => number
E := 2.718281828459045
exp := x => pow(E, x)
sqrt := x => pow(x, 0.5)

# ln : number => number # natural log
log := (x, base) => ln(x) / ln(base)

# sin : number => number  # sine
# cos : number => number  # cosine
tan := x => sin(x) / cos(x)

# asin : number => number # arcsine (inverse sin)
# acos : number => number # arccosine (inverse cos)
PI := acos(0) * 2

ternary := (p, x, y) => true :: {
  p -> x
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
# floor : number => number # floor / truncation
round := x => floor(x + 0.5)
trunc := x => ternary(x > 0, floor, ceil)(x)

{
  sqrt
  E
  exp
  log
  tan
  PI
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