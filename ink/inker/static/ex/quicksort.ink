# minimal quicksort implementation using hoare partition

{map, range, log, clone} := import('std.ink')

swap! := (list, i, j) => (
  tmp := list.(i)
  list.(i) := list.(j)
  list.(j) := tmp
)

sortBy := (v, pred) => (
  vPred := map(v, (x, _) => pred(x))
  partition := (v, lo, hi) => (
    pivot := vPred.(lo)
    lsub := i => vPred.(i) < pivot :: {
      true -> lsub(i + 1)
      false -> i
    }
    rsub := j => vPred.(j) > pivot :: {
      true -> rsub(j - 1)
      false -> j
    }
    (sub := (i, j) => (
      i := lsub(i)
      j := rsub(j)
      i < j :: {
        false -> j
        true -> (
          swap!(v, i, j)
          swap!(vPred, i, j)
          sub(i + 1, j - 1)
        )
      }
    ))(lo, hi)
  )
  (quicksort := (v, lo, hi) => true :: {
    len(v) == 0 | ~(lo < hi) -> v
    _ -> (
      p := partition(v, lo, hi)
      quicksort(v, lo, p)
      quicksort(v, p + 1, hi)
    )
  })(v, 0, len(v) - 1)
)

sort! := v => sortBy(v, x => x)
sort := v => sort!(clone(v))

# TEST
L := map(range(0, 250, 1), (_, _) => floor(rand() * 500))
log('before quicksort: ' + string(L))
log('after quicksort: ' + string(sort(L)))
