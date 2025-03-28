# a primitive HTTP server

{log} := import('logging.ink')

close := listen('0.0.0.0:8080', evt => (
  log(evt)
  evt.type :: {
    'error' -> log('Error: ' + evt.message)
    'req' -> {
      status: 200
      headers: {'Content-Type': 'text/plain'}
      body: 'Hello, World!'
    }
  }
))

wait(5)
close()
