` minimal quicksort implementation
  using hoare partition `

map := load('functional').map
clone := load('std').clone

sortBy := (v, pred) => (
  partition := (v, lo, hi) => (
    pivot := pred(v.(lo))
    lsub := i => pred(v.(i)) < pivot :: {
      true -> lsub(i + 1)
      false -> i
    }
    rsub := j => pred(v.(j)) > pivot :: {
      true -> rsub(j - 1)
      false -> j
    }
    (sub := (i, j) => (
      i := lsub(i)
      j := rsub(j)
      i < j :: {
        false -> j
        true -> (
          ` inlined swap! `
          tmp := v.(i)
          v.(i) := v.(j)
          v.(j) := tmp

          sub(i + 1, j - 1)
        )
      }
    ))(lo, hi)
  )
  (quicksort := (v, lo, hi) => true :: {
    len(v) = 0 -> v
    lo < hi -> (
      p := partition(v, lo, hi)
      quicksort(v, lo, p)
      quicksort(v, p + 1, hi)
    )
    _ -> v
  })(v, 0, len(v) - 1)
)

sort! := v => sortBy(v, x => x)

sort := v => sort!(clone(v))

` TEST `
range := load('functional').range
log := load('logging').log

rint := () => floor(rand() * 500)
L := map(range(0, 250, 1), rint)
Before := clone(L)
log('before quicksort: ' + string(L))
log('after quicksort: ' + string(sort(L)))
log('before intact?: ' + (L :: {Before -> 'yes', _ -> 'no'}))
sort!(L)
log('after mutable sort: ' + string(L))
