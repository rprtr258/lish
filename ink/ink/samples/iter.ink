# Iter(T) :: () => T | ()

# Iter(T)
empty := () => ()

generate := (x0, f) => (
  x := {value: x0}
  () => (
    res := x.value
    x.value = f(x.value)
    res
  )
)

# () => Iter(int)
count := () => generate(1, (n) => n + 1)

# (Iter(A), Iter(B)) => Iter([A, B])
zip := (l, r) => () => (
  a := l()
  b := r()
  true :: {
    (a == ()) | (b == ()) -> ()
    _ -> [a, b]
  }
)

# TODO: make methods of iterator
list := (list) => (
  i := {value: 0}
  () => true :: {
    i.value < len(list) -> (
      res := list.(i.value)
      i.value = i.value + 1
      res
    )
    _ -> ()
  }
)

map := (it, f) => () => (
  x := it() :: {
    () -> ()
    _ -> f(x)
  }
)

filter := (it, f) => () => (
  x := it()
  true :: {
    x == () -> ()
    f(x) -> x
    _ -> filter(it, f)
  }
)

takeWhile := (it, f) => (
  stopped := {value: false}
  () => true :: {
    stopped.value -> ()
    _ -> (
      x := it()
      true :: {
        ~(x == ()) & f(x) -> x
        _ -> (
          stopped.value := true
          ()
        )
      }
    )
  }
)

# (number, number, number) => Iter(number)
range := (start, end, step) => (
  true :: {
    step == 0 -> empty
    step < 0 -> takeWhile(generate(start, (x) => x + step), (x) => x > end)
    step > 0 -> takeWhile(generate(start, (x) => x + step), (x) => x < end)
  }
)

foldl := (it, f, acc) => (
  x := it()
  x :: {
    () -> acc
    _ -> foldl(it, f, f(acc, x))
  }
)

each := (it, f) => (
  x := it()
  x :: {
    () -> ()
    _ -> (
      f(x)
      each(it, f)
    )
  }
)

collect := (it) => foldl(it, (acc, x) => acc.(len(acc)) := x, [])

{
  range
  map
  foldl
  generate
  takeWhile
  collect
  each
}