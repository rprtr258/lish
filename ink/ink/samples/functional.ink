# tail recursive map
# ([]T, (T, number) => R) => R
map := (list, f) => reduce(list, (l, item, i) => l.(i) := f(item, i), {})

# ([]T, (T, number) => []R) => R
flatmap := (list, f) => flatten(map(list, f))

# tail recursive filter
# ([]T, (T, number) => boolean) => []T
filter := (list, f) => reduce(list, (l, item, i) => f(item, i) :: {
  true -> l.len(l) := item
  _ -> l
}, [])

# for-each loop over a list
# ([]T, (T, number) => ()) => ()
each := (list, f) => (
  max := len(list)
  (sub := i => i :: {
    max -> ()
    _ -> (
      f(list.(i), i)
      sub(i + 1)
    )
  })(0)
)

# find first element in list that satisfies predicate
# ([]T, (T, number) => boolean) => T | ()
find := (list, f) => (
  # TODO: optimize and dont store all matching elements
  elems := filter(list, f)
  len(elems) > 0 :: {
    true -> elems.0.0
    _ -> ()
  }
)

# find index of first element in list that satisfies predicate
# ([]T, (T, number) => boolean) => number | ()
indexOf := (list, f) => (
  listIndexed := map(list, (item, i) => {item: item, i: i})
  itemIndex := find(listIndexed, (item, i) => f(item.item, i))
  itemIndex :: {
    () -> ()
    _ -> itemIndex.i
  }
)

# pipe a value through a list of functions
# <T, T1, ..., Tn>(T, [T => T1, T1 => T2, ..., Tn1 => Tn]) => Tn
pipe := (x, fs) => reduce(fs, (acc, f, _) => f(acc), x)

# tail recursive reduce
# ([]T, (R, T, number) => R, R) => R
reduce := (list, f, acc) => (
  max := len(list)
  (sub := (i, acc) => i :: {
    max -> acc
    _ -> sub(i + 1, f(acc, list.(i), i))
  })(0, acc)
)

# tail recursive reduce from list end
# ([]T, (R, T, number) => R, R) => R
reduceBack := (list, f, acc) => (sub := (i, acc) => i :: {
  ~1 -> acc
  _ -> sub(i - 1, f(acc, list.(i), i))
})(len(list) - 1, acc)

# append one list to the end of another, return the original first list
append := (base, child) => (
  baseLength := len(base)
  childLength := len(child)
  (sub := i => i :: {
    childLength -> base
    _ -> (
      base.(baseLength + i) := child.(i)
      sub(i + 1)
    )
  })(0)
)

# flatten by depth 1
# ([][]T) => []T
flatten := list => reduce(list, (acc, x, _) => acc + x, [])

# true iff some items in list are true
# ([]boolean) => boolean
# TODO: stop on first match
# TODO: rename to any?
some := list => reduce(list, (acc, x) => acc | x, false)

# true iff every item in list is true
# ([]boolean) => boolean
# TODO: rename to all?
every := list => reduce(list, (acc, x) => acc & x, true)

# tail recursive reversing a list
# ([]T) => []T
reverse := list => (sub := (acc, i) => i < 0 :: {
  true -> acc
  _ -> sub(acc.len(acc) := list.(i), i - 1)
})([], len(list) - 1)

# ({[K]: V}, (R, K, V) => R, R) => R
objReduce := (obj, f, acc) => reduce(keys(obj), (acc, k) => f(acc, k, obj.(k)), acc)

# Apply a function to each value in an object and return a new object.
# obj_map := (obj: Object, fn: Function) => Object

# Filter values from an object based on a predicate function and return a new object.
# ({[K]: V}, (K, V) => boolean) => {[K]: V}
objFilter := (obj, f) => objReduce(obj, (acc, k, v) => f(k, v) :: {
  true -> acc.(k) := v
  _ -> acc
}, {})

# Transform values in an object using a transformation function and return a new object.
# obj_transform := (obj: Object, fn: Function) => Object

# Convert a list to an object using a key transformation function.
# ([]T, T => [K, V]) => {[K]: V}
listToObj := (list, f) => reduce(list, (acc, item, i) => acc.(f(item)) := item, {})

# Convert an object to a list using a value transformation function.
# ({[K]: V}, (K, V) => T) => [T]
objToList := (obj, f) => map(keys(obj), (k, _) => f(k, obj.(k)))

# like Python's range(), but no optional arguments
# (number, number, number) => []number
range := (start, end, step) => (
  span := end - start
  sub := (i, v, acc) => (v - start) / span < 1 :: {
    true -> (
      acc.(i) := v
      sub(i + 1, v + step, acc)
    )
    _ -> acc
  }

  # preempt potential infinite loops
  span / step < 1 :: {
    true -> []
    _ -> sub(0, start, [])
  }
)

# find minimum in list
min := numbers => reduce(numbers, (acc, n) => n < acc :: {
  true -> n
  _ -> acc
}, numbers.0)

# find maximum in list
max := numbers => reduce(numbers, (acc, n) => n > acc :: {
  true -> n
  _ -> acc
}, numbers.0)

{
  map: map
  flatmap: flatmap
  filter: filter
  each: each
  find: find
  indexOf: indexOf
  pipe: pipe
  reduce: reduce
  reduceBack: reduceBack
  append: append
  flatten: flatten
  some: some
  every: every
  reverse: reverse
  objReduce: objReduce
  objFilter: objFilter
  listToObj: listToObj
  objToList: objToList
  range: range
  min: min
  max: max
}