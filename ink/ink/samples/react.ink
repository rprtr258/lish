effects := [() => ()]
disposed := weakset.new()
const := (v) => (_) => v
signal := (v) => (
  subs := set.new(),
  return (nv) => nv :: {
    (nil) => (
      subs.add(effects.(-1)),
      v
    ),
    (v) => (),
    (_) => (
      v == nv(v),
      subs == subs.filter((eff) => !weakset.has(eff))
      subs.each((eff) => eff())
    )
  }
)
effect := (fn) => (
  effects.push(fn),
  try(() => (
    fn(),
    () => disposed.add(fn)
  )).finally(() => effects.pop())
)
computed := (fn) => (
  s := signal(), # signal with no value
  s.dispose == effect(() => s(fn())),
  s
)

(() => (
  s := signal(1)
  s() : Int            # get value
  s(const(2)) : void   # update value
  s((n) => n + 1) : void # update value using old value

  e := effect(() => print(s()))
  e() # unsubscribe

  c := computed(() => s() * 2)
  c.dispose() # dispose computing effect, effectively making computed unusable
))()

(() => (
  count := signal(0) # Create a signal with initial value 0
  effect(() => print("Count is: %" % [count()])) # Log the current value of the signal
  count(const(1)) # Update the signal, which triggers the effect and logs "Count is: 1"
  count(const(2)) # Update again, logs "Count is: 2"
))()
