` main.ink serves files `

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
log := std.log
f := std.format
readFile := std.readFile
str := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/str.ink')
contains? := str.contains?
http := import('../vendor/http.ink')
mimeForPath := import('../vendor/mime.ink').forPath

Port := 9999

server := (http.new)()
NotFound := {status: 404, body: 'file not found'}
MethodNotAllowed := {status: 405, body: 'method not allowed'}

serveStatic := path => (req, end) => req.method :: {
  'GET' -> readFile('static/' + path, file => file :: {
    () -> end(NotFound)
    _ -> end({
      status: 200
      headers: {'Content-Type': mimeForPath(path)}
      body: file
    })
  })
  _ -> end(MethodNotAllowed)
}

serveGZip := (path, end) => readFile('static/indexes/' + path + '.gz', file => file :: {
  () -> end(NotFound)
  _ -> end({
    status: 200
    headers: {
      'Content-Type': mimeForPath(path)
      'Content-Encoding': 'gzip'
    }
    body: file
  })
})
serveBrotli := (path, end) => readFile('static/indexes/' + path + '.br', file => file :: {
  () -> end(NotFound)
  _ -> end({
    status: 200
    headers: {
      'Content-Type': mimeForPath(path)
      'Content-Encoding': 'br'
    }
    body: file
  })
})
serveCompressed := path => (req, end) => req.method :: {
  'GET' -> acceptEncoding := req.headers.'Accept-Encoding' :: {
    () -> serveGZip(path, end)
    ` we check brotli compatibility before serving, as it is not as
    ubiquitous as gzip `
    _ -> contains?(acceptEncoding, 'br') :: {
      true -> serveBrotli(path, end)
      _ -> serveGZip(path, end)
    }
  }
  _ -> end(MethodNotAllowed)
}

addRoute := server.addRoute

` static paths `
addRoute('/static/*staticPath', params => serveStatic(params.staticPath))
addRoute('/indexes/*indexPath', params => serveCompressed(params.indexPath))
addRoute('/favicon.ico', params => serveStatic('favicon.ico'))
addRoute('/', params => serveStatic('index.html'))

start := () => (
  end := (server.start)(Port)
  log(f('Monocle started, listening on 0.0.0.0:{{0}}', [Port]))
)

start()

