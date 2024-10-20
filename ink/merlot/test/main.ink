s := (import('../vendor/mod.ink').suite)(
  'Merlot test suite'
)

import('md.ink')(s.mark, s.test)
import('reader.ink')(s.mark, s.test)
import('util.ink')(s.mark, s.test)

(s.end)()
