# simple ping-pong request-response test over HTTP

{log} := import('logging.ink')
{format: f} := import('str.ink')

# helper for logging errors
logErr := msg => log('error: ' + msg)

# start a server
closeServer := listen('0.0.0.0:9600', evt => evt.type :: {
  'error' -> logErr(evt.message)
  'req' -> (
    log(f('Request ---> {{ data }}', evt))

    {data: dt, end} := evt
    [dt.method, dt.url, dt.body] :: {
      ['POST', '/test', 'ping'] -> end({
        status: 302 # test that it doesn't auto-follow redirects
        headers: {
          'Content-Type': 'text/plain'
          'Location': 'https://dotink.co'
        }
        body: 'pong'
      })
      _ -> end({
        status: 400
        headers: {
          'Content-Type': 'text/plain'
        }
        body: 'invalid request!'
      })
    }
  )
})

# send a request
send := () => (
  closeRequest := req({
    method: 'POST'
    url: 'http://127.0.0.1:9600/test'
    headers: {
      'Accept': 'text/html'
    }
    body: 'ping'
  }, evt => evt.type :: {
    'error' -> logErr(evt.message)
    'resp' -> (
      log(f('Response ---> {{ data }}', evt))

      dt := evt.data
      [dt.status, dt.body] :: {
        [302, 'pong'] -> (
          log('---> ping-pong, success!')
          closeServer()
        )
        _ -> logErr('communication failed!')
      }
    )
  })

  # half-second timeout on the request
  wait(0.5, closeRequest)
)

# give server time to start up before sending first request
wait(0.5, send)
