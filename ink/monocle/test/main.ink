runTokenizerTests := import('tokenizer.ink').run

s := (import('../vendor/suite.ink').suite)(
  'Monocle test suite'
)

runTokenizerTests(s.mark, s.test)

(s.end)()

