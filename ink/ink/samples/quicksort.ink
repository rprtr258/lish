# minimal quicksort implementation
# using hoare partition

{clone} := import('std.ink')

sortBy := (v, pred) => (
  partition := (v, lo, hi) => (
    pivot := pred(v.(lo))
    lsub := i => true :: {
      pred(v.(i)) < pivot -> lsub(i + 1)
      _ -> i
    }
    rsub := j => true :: {
      pred(v.(j)) > pivot -> rsub(j - 1)
      _ -> j
    }
    (sub := (i, j) => (
      i := lsub(i)
      j := rsub(j)
      true :: {
        i < j -> (
          # inlined swap!
          tmp := v.(i)
          v.(i) := v.(j)
          v.(j) := tmp

          sub(i + 1, j - 1)
        )
        _ -> j
      }
    ))(lo, hi)
  )
  (quicksort := (v, lo, hi) => true :: {
    len(v) == 0 -> v
    lo < hi -> (
      p := partition(v, lo, hi)
      quicksort(v, lo, p)
      quicksort(v, p + 1, hi)
    )
    _ -> v
  })(v, 0, len(v) - 1)
)

sort! := v => sortBy(v, x => x)

{
  sort!
  sort: v => sort!(clone(v))
}
