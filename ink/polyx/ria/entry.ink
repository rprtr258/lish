http := import('../lib/http.ink')

PORT := 7284

server := (http.new)()

addRoute := server.addRoute
# addRoute('/static', staticHandler)

(server.start)(PORT)
