{format: f} := import('str.ink')
{critical: fatal} := import('logging.ink')
{each, range: rawRange} := import('iter.ink')

range := n => rawRange(0, n, 1)

len(args()) :: {
  3 -> ()
  _ -> fatal(f('Usage: {{}} "TEXT"', args().0))
}

clrs := [31, 33, 32, 36, 34, 35]
s := args().2
half := floor(len(s) / 2)
each(range(100), colshift => (
  each(range(floor(half * (1 + cos(colshift / 4)))), _ => out(' '))
  j := {value: 0}
  each(range(len(s)), i => (
    out(f('[{{0}}m', [clrs.((j.value + colshift) % len(clrs))])+s.(i))
    j.value := j.value + (s.(i) :: {' ' -> 1, _ -> 0})
  ))
  out('\n')
))