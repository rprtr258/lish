`
(type) => type
Future := (T) => {
  ready?: () => boolean
  try_await: () => T | () # non-blocking
  await: () => T # blocking
}
`

# ((T => R) => R) => Future(T)
future := f => {
  res := {value: ()}
  f(value => res.value = value)
  {
    ready?: () => res.value :: {() -> false, _ -> true}
    try_await: () => res.value
    await: () => res.value # TODO: wait until ready
  }
}

# (T) => Future(T)
instant := value => future((f) => f(value))

# (Future(T), T => R) => Future(R)
map := (v, f) => future((g) => g(f(v.await())))

# (Future(T), T => Future(R)) => Future(R)
flatmap := (v, f) => g(f(v.await()))