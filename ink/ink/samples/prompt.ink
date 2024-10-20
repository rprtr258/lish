# scan() / in() based prompt demo

log := import('logging.ink').log
scan := import('std.ink').scan

ask := (question, cb) => (
  log(question)
  scan(cb)
)

ask('What\'s your name?', name =>
  log('Great to meet you, ' + name + '!')
)
