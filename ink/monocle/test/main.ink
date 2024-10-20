s := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/6b87250c1cc7f5962f40c9f85f656bf3fb9c55c6/suite.ink')(
  'Monocle test suite'
)

import('tokenizer.ink')(s.mark, s.test)

(s.end)()

