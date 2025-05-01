dbg := x => (out(string(x)+'\n'), x)

run :=            k => k(v => v)
const :=     x => k => k(x)
plus := (x, y) => k => x(xv => y(yv => k(xv + yv)))
mul :=  (x, y) => k => x(xv => y(yv => k(xv * yv)))
str :=       x => k => x(xv => k(string(xv)))
print :=     x => k => str(x)(s => k(out(s)))
println :=   x => k => str(x)(xv => k(out(xv + '\n')))
list :=     xs => k => (sub := (i, ys) => i :: {
  len(xs) -> k(ys)
  _ -> (xs.(i))(yi => sub(i + 1, ys + [yi]))
})(0, [])

# println(2 * 3) # 6
run(println(mul(const(2), const(3))))

reset := f => k => (
  shifted := {box: false}
  v := f(fk => (
    shifted.box := true
    shifted.v := fk(x => f(_ => const(x)))
  ))
  (true :: {
    shifted.box -> shifted.v
    _ -> v
  })(k)
)
run(println(list([
  # reset(shift(k => k(3))) # 3
  reset(shift => shift(k => k(3)))
  # reset(1 * shift(k => k(3))) # 3
  reset(shift => mul(const(1), shift(k => k(3))))
  # reset(2 * shift(k => k(3))) # 6
  reset(shift => mul(const(2), shift(k => k(3))))
  reset(shift => shift(mul(const(2), const(3)))) # 6
  reset(_ => mul(const(2), const(3))) # 6
  # 1 + reset(2 * 3) # 7
  plus(const(1), reset(_ => mul(const(2), const(3))))
  # 1 + reset(2 * shift(k => 3)) # 4
  # plus(const(1), reset(shift => mul(const(2), shift(3))))
  # 1 + reset(2 * shift(k => k(3))) # 7
  plus(const(1), reset(shift => mul(const(2), shift(const(3)))))
  # 1 + reset(2 * shift(k => k(2) + k(2))) # 9
  plus(const(1), reset(shift => mul(const(2), shift(plus(const(2), const(2))))))
])))
