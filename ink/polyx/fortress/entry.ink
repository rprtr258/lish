http := import('../lib/http.ink')

PORT := 7281

server := (http.new)()

addRoute := server.addRoute
# addRoute('/static', staticHandler)

(server.start)(PORT)
