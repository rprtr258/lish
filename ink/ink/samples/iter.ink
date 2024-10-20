# Iter(T) = () => T | ()

# Iter(T)
empty := () => ()

# (number, number, number) => Iter(number)
range := (start, end, step) => (
  true :: {
    step = 0 -> empty
    step < 0 -> takeWhile(generate(start, x => x + step), x => x > end)
    step > 0 -> takeWhile(generate(start, x => x + step), x => x < end)
  }
)

generate := (x0, f) => (
  x := {value: x0}
  () => (
    res := x.value
    x.value := f(x.value)
    res
  )
)

# TODO: make methods of iterator
list := list => (
  i := {value: 0}
  () => i.value < len(list) :: {
    true -> (
      res := list.(i.value)
      i.value := i.value + 1
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
    x = () -> ()
    f(x) -> x
    _ -> filter(it, f)
  }
)

takeWhile := (it, f) => (
  stopped := {value: false}
  () => stopped.value :: {
    true -> ()
    _ -> (
      x := it()
      ~(x = ()) & f(x) :: {
        true -> x
        _ -> (
          stopped.value := true
          ()
        )
      }
    )
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

collect := it => foldl(it, (acc, x) => acc.(len(acc)) := x, [])