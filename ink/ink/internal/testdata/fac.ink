fac := n => true :: {
  n < 1 -> 1
  _ -> n * fac(n - 1)
}
{fac}
