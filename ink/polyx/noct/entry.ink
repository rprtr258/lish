` noct server `

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
log := std.log
readFile := std.readFile
writeFile := std.writeFile
json := import('../vendor/json.ink')

http := import('../lib/http.ink')
cli := import('../lib/cli.ink')
pctDecode := import('../lib/percent.ink').decode

fs := import('fs.ink')
cleanPath := fs.cleanPath
describe := fs.describe
flatten := fs.flatten
ensurePDE := fs.ensureParentDirExists
sync := import('sync.ink')

Port := 7280

given := (cli.parsed)()
givenPath := (given.args.0 :: {
  () -> '.'
  _ -> given.args.0
})
RootFS := cleanPath(givenPath)

server := (http.new)()

addRoute := server.addRoute
addRoute('/desc/*descPath', params => (_, end) => (
  descPath := RootFS + '/' + cleanPath(pctDecode(params.descPath))
  describe(descPath, RootFS, desc => end({
    status: 200
    body: (json.ser)(desc)
  }))
))
addRoute('/desc/', _ => (_, end) => (
  describe(RootFS, RootFS, desc => end({
    status: 200
    body: (json.ser)(desc)
  }))
))
addRoute('/sync/*downPath', params => (req, end) => req.method :: {
  'GET' -> (
    downPath := RootFS + '/' + cleanPath(pctDecode(params.downPath))
    readFile(downPath, file => file :: {
      () -> end({
        status: 404
        body: 'file not found'
      })
      _ -> end({
        status: 200
        headers: {
          'Content-Type': 'application/octet-stream'
        }
        body: file
      })
    })
  )
  'POST' -> (
    downPath := RootFS + '/' + cleanPath(pctDecode(params.downPath))
    ensurePDE(downPath, r => r :: {
      true -> writeFile(downPath, req.body, r => r :: {
        true -> end({
          status: 201
          body: ''
        })
        _ -> end({
          status: 500
          body: 'upload failed, could not write file'
        })
      })
      _ -> end({
        status: 500
        body: 'upload failed, could not create dir'
      })
    })
  )
  _ -> end({
    status: 400
    body: 'invalid request'
  })
})

start := () => (
  close := (server.start)(Port)
  log('Noct server started with fs in ' + RootFS)
)
