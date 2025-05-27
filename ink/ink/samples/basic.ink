# first program written in Ink, kept for historical reasons

log := (s) => out(s + '\n')

fn1 := () => log('Hello, World!')

fn2 := () => (
  log('Hello, World 2!')
)

out('Hello test\n')

log('Hello with \' apostrophe test')

(
  fn1()
  fn2()
)
