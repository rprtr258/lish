# mime type

split := import('../vendor/str').split

MimeTypes := {
  'html': 'text/html'
  'css':  'text/css'
  'ink':  'text/plain'

  'js':   'application/javascript'
  'json': 'application/json'

  'jpg': 'image/jpeg'
  'png': 'image/png'
  'svg': 'image/svg+xml'
}

forPath := path => (
  parts := split(path, '.')
  ending := parts.(len(parts) - 1)

  guess := MimeTypes.(ending) :: {
  	() -> 'application/octet-stream'
  	_ -> guess
  }
)
