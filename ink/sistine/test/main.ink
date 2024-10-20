# runMarkdownTests := import('md').run
# runReaderTests := import('reader').run

s := (import('../vendor/suite').suite)(
  'Sistine test suite'
)

# runMarkdownTests(s.mark, s.test)
# runReaderTests(s.mark, s.test)

(s.end)()

