runTokenizerTests := import('tokenizer').run

s := (import('../vendor/suite').suite)(
	'Monocle test suite'
)

runTokenizerTests(s.mark, s.test)

(s.end)()

