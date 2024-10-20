# first program written in Ink, kept for historical reasons

log := import('std.ink').log

fn1 := n => log('Hello, World!')

fn2 := () => (
  log('Hello, World 2!')
)

out('Hello test\n')

log('Hello with \' apostrophe test')

(
  fn1()
  fn2(1, 2, false)
)
