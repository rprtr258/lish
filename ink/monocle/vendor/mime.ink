` mime type `

str := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/str.ink')
split := str.split

MimeTypes := {
  'html': 'text/html'
  'css': 'text/css'
  'js': 'application/javascript'
  'json': 'application/json'
  'ink': 'text/plain'

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

{
  forPath: forPath
}