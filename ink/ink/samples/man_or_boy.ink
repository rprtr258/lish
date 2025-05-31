# a := (k int, x1, x2, x3, x4, x5 func() int) int {
a := (k, x1, x2, x3, x4, x5) => (
  k1 := [k]
  b := () => (
    k1.0 := k1.0 - 1
    a(k1.0, b, x1, x2, x3, x4)
  )
  true :: {
    k1.0 > 0 -> b()
    _        -> x4() + x5()
  }
)

x := (i) => () => i
k := 8 # 10
out(string(a(k, x(1), x(~1), x(~1), x(1), x(0)))+'\n')
# TODO: assert -67 on k=10
