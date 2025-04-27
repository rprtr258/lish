# ink standard test suite tools

{log} := import('logging.ink')
{each} := import('functional.ink')
{format: f} := import('str.ink')

(label, suite) => (
  # suite data store
  s := {
    all: 0
    passed: 0
  }

  # perform a new test case
  indent := '        '
  # mark sections of a test suite with human labels
  mark := (label, test_fn) => (
    log(f('suite: {{}}', label))
    test := (label, result, expected) => (
      s.all = s.all + 1
      true :: {
        result == expected -> s.passed = s.passed + 1
        _ -> log(f('  * {{ label }}
    {{ indent }}got {{ result }}
    {{ indent }}exp {{ expected }}', {label, result, expected, indent}))
      }
    )
    test_fn(test)
  )

  suite(mark)

  # print out results
  true :: {
    s.passed == s.all -> log(f('ALL {{ passed }} / {{ all }} PASSED', s))
    _ -> (
      log(f('PARTIAL: {{ passed }} / {{ all }} PASSED', s))
      exit(1)
    )
  }
)