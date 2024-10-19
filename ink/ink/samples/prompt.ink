# scan() / in() based prompt demo

log := load('logging').log
scan := load('std').scan

ask := (question, cb) => (
  log(question)
  scan(cb)
)

ask('What\'s your name?', name =>
  log('Great to meet you, ' + name + '!')
)
