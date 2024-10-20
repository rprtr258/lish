runMarkdownTests := import('md').run
runReaderTests := import('reader').run
runUtilTests := import('util').run

s := (import('../vendor/suite').suite)(
  'Merlot test suite'
)

runMarkdownTests(s.mark, s.test)
runReaderTests(s.mark, s.test)
runUtilTests(s.mark, s.test)

(s.end)()
