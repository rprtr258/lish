# test: std.clone

{log, clone} := import('./runtime/std')

x := {
  key: 'value'
  k2: 2.31
  ork: [1, 2, 3]
}
log(x)
y := clone(x)
log(y)

x.key := 'v2'
x.ork.len(x.ork) := 9
log(x)
log(y)
