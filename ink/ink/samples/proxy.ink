# basic HTTP proxy

{log} := import('logging.ink')
{format: f} := import('str.ink')
{slice} := import('std.ink')

PORT := 7900

# map of route prefixes to proxied endpoints
PROXIES := {
  '/gh':     'https://api.github.com'
  '/github': 'https://api.github.com'
}

# headers provided for error responses
DefaultHeaders := {
  'Content-Type': 'text/plain; charset=utf-8'
  'X-Proxied-By': 'ink-proxy'
  'X-Served-By': 'ink-serve'
}

# handles when proxied request succeeds
handleProxyResponse := (dest, data) => (
  log(f('Proxied {{ dest }} success', {
    dest
  }))
  {
    status: data.status
    headers: data.headers.('X-Proxied-By') := 'ink-proxy'
    body: data.body
  }
)

# handles when proxied request fails
handleProxyError := (dest, data) => (
  log(f('Error in proxied request to {{ dest }}: {{ err }}', {
    dest
    err: data.message
  }))
  {
    status: 502
    headers: DefaultHeaders
    body: f('proxied service {{ dest }} was not available for {{ url }}', {
      dest
      url: data.url
    })
  }
)

# responds to all requests to the proxy
handleRequest := data => (
  prefixes := keys(PROXIES)
  max := len(prefixes) - 1

  (sub := i => (
    prefix := prefixes.(i)

    # check that proxy prefix matches exactly.
    # i.e. /gh/ should match but /ghub should not
    slice(data.url + '/', 0, len(prefix) + 1) :: {
      prefix + '/' -> (
        dest := PROXIES.(prefix) + slice(data.url, len(prefix), len(data.url))
        evt := req({
          method: data.method
          url: dest
          headers: data.headers.('X-Proxied-By') := 'ink-proxy'
          body: data.body
        })
        evt => evt.type :: {
          'error' -> handleProxyError(dest, evt.data)
          'resp' -> handleProxyResponse(dest, evt.data)
        }
      )
      _ -> i :: {
        max -> {
          status: 404
          headers: DefaultHeaders
          body: 'could not locate proxy destination for ' + data.url
        }
        _ -> sub(i + 1)
      }
    }
  ))(0)
)

listen('0.0.0.0:' + string(PORT), evt => evt.type :: {
  'error' -> log('Error starting server:' + evt.message)
  'req' -> handleRequest(evt.data)
})
