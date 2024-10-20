# ink standard test suite tools

# borrow from std
log := import('logging').log
each := import('functional').each
f := import('str').format

# suite constructor
label => (
  # suite data store
  s := {
    all: 0
    passed: 0
    msgs: []
  }

  # mark sections of a test suite with human labels
  mark := label => s.msgs.(len(s.msgs)) := '- ' + label

  # signal end of test suite, print out results
  end := () => (
    log(f('suite: {{}}', label))
    each(s.msgs, m => log('  ' + m))
    s.passed :: {
      s.all -> log(f('ALL {{ passed }} / {{ all }} PASSED', s))
      _ -> (
        log(f('PARTIAL: {{ passed }} / {{ all }} PASSED', s))
        exit(1)
      )
    }
  )

  # perform a new test case
  indent := '        '
  test := (label, result, expected) => (
    s.all := s.all + 1
    result :: {
      expected -> s.passed := s.passed + 1
      _ -> s.msgs.(len(s.msgs)) := f('  * {{ label }}
  {{ indent }}got {{ result }}
  {{ indent }}exp {{ expected }}', {
        label: label
        result: result
        expected: expected
        indent: indent
      })
    }
  )

  # expose API functions
  {
    mark: mark
    test: test
    end: end
  }
)
