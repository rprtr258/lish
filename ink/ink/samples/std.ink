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

# get a substring of a given string, or sublist of a given list
slice := (s, start, end) => (
  # clamp start and end numbers to ranges, such that
  # start < end. Utility used in slice
  clamp := (start, end, min, max) => (
    {max: mmax, min: mmin} := import('math.ink')
    start := mmax(start, min)
    end := mmax(end, min)
    end := mmin(end, max)
    start := mmin(start, end)

    {start, end}
  )

  # bounds checks
  x := clamp(start, end, 0, len(s))
  {start} := x
  max := x.end - start

  (sub := (i, acc) => i :: {
    max -> acc
    _ -> sub(i + 1, acc.(i) := s.(start + i))
  })(0, type(s) :: {
    'string' -> ''
    'list' -> []
  })
)

# clone a composite value
clone := x => (
  {reduce} := import('functional.ink')

  type(x) :: {
    'string' -> '' + x
    'composite' -> reduce(keys(x), (acc, k, _) => acc.(k) := x.(k), {})
    'list' -> reduce(keys(x), (acc, i, _) => acc.(i) := x.(i), [])
    _ -> x
  }
)

{scan, slice, clone}
