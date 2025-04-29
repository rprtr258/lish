{each, flatmap, filter} := import('samples/functional.ink')

# TODO: break - ends loop entirely
# TODO: continue - starts next iteration
# TODO: prune/skip - does not recurse on current value, continues
for_dfs := (init, cond, next, f) => (sub := n => true :: {
  cond(n) -> (
    f(n)
    each(next(n), (m, _) => sub(m))
  )
})(init)

for_bfs := (init, cond, next, f) => (sub := ns => len(ns) :: {
  0 -> ()
  _ -> (
    ns := filter(ns, (n, _) => cond(n))
    each(ns, (n, _) => f(n))
    sub(flatmap(ns, n => next(n)))
  )
})([init])

for_dfs('', s => len(s) < 3, s => [s+'a', s+'b', s+'c'], s => out(s+' '))
out('\n')
for_bfs('', s => len(s) < 3, s => [s+'a', s+'b', s+'c'], s => out(s+' '))
out('\n')
