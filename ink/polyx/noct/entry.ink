` noct server `

std := load('../vendor/std')
log := std.log
readFile := std.readFile
writeFile := std.writeFile
json := load('../vendor/json')

http := load('../lib/http')
cli := load('../lib/cli')
pctDecode := load('../lib/percent').decode

fs := load('fs')
cleanPath := fs.cleanPath
describe := fs.describe
flatten := fs.flatten
ensurePDE := fs.ensureParentDirExists
sync := load('sync')

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
