# the ink standard library

scan := cb => (
  acc := ['']
  in(evt => evt.type :: {
    'end' -> cb(acc.0)
    'data' -> (
      acc.0 := acc.0 + slice(evt.data, 0, len(evt.data) - 1)
      false
    )
  })
)

# clamp start and end numbers to ranges, such that
# start < end. Utility used in slice
clamp := (start, end, min, max) => (
  m := load('math')
  start := (m.max)(start, min)
  end := (m.max)(end, min)
  end := (m.min)(end, max)
  start := (m.min)(start, end)

  {
    start: start
    end: end
  }
)

# get a substring of a given string, or sublist of a given list
slice := (s, start, end) => (
  # bounds checks
  x := clamp(start, end, 0, len(s))
  start := x.start
  max := x.end - start

  (sub := (i, acc) => i :: {
    max -> acc
    _ -> sub(i + 1, acc.(i) := s.(start + i))
  })(0, type(s) :: {
    'string' -> ''
    'composite' -> []
  })
)

# clone a composite value
clone := x => type(x) :: {
  'string' -> '' + x
  'composite' -> (load('functional').reduce)(keys(x), (acc, k) => acc.(k) := x.(k), {})
  _ -> x
}
