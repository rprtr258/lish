fract := x => x - floor(x)

ceil := x => fract(x) :: {
  0 -> x
  _ -> fract(x + 1)
}
# floor := floor
round := x => floor(x + 0.5)
trunc := x => x > 0 :: {
  true -> floor(x)
  false -> ceil(x)
}

abs := x => x :: {
  x < 0 -> ~x
  _ -> x
}

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
  false -> y
}
min := (x, y) => ternary(x < y, x, y)
max := (x, y) => ternary(x > y, x, y)