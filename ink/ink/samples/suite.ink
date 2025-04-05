# ink standard test suite tools

# borrow from std
{log} := import('logging.ink')
{each} := import('functional.ink')
{format: f} := import('str.ink')

# suite constructor
(label, suite) => (
  # suite data store
  s := {
    all: 0
    passed: 0
    msgs: []
  }

  # mark sections of a test suite with human labels
  mark := label => s.msgs.len(s.msgs) := '- ' + label

  # perform a new test case
  indent := '        '
  test := (label, result, expected) => (
    s.all := s.all + 1
    true :: {
      result = expected -> s.passed := s.passed + 1
      _ -> s.msgs.(len(s.msgs)) := f('  * {{ label }}
  {{ indent }}got {{ result }}
  {{ indent }}exp {{ expected }}', {label, result, expected, indent})
    }
  )

  suite({mark, test})

  # print out results
  log(f('suite: {{}}', label))
  each(s.msgs, m => log('  ' + m))
  true :: {
    s.passed = s.all -> log(f('ALL {{ passed }} / {{ all }} PASSED', s))
    _ -> (
      log(f('PARTIAL: {{ passed }} / {{ all }} PASSED', s))
      exit(1)
    )
  }
)