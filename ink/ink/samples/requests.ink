`
Response := {
  type: 'error'
  message: string
} | {
  type: 'resp'
  data: {
    status: number
    headers: {[string]: string}
    body: string
  }
}

({
  method: string
  url: string
  headers: {[string]: string}
  body: string
}, Response => T) => T
TODO: ({
  method: string
  url: string
  headers: {[string]: string}
  body: string
}) => Future(Response)
`
# req := req

methodNoBody := method => (url, headers, cb) => req({
  method: method
  url: url
  headers: headers
}, cb)

methodBody := method => (url, body, headers, cb) => req({
  method: method
  url: url
  headers: headers
  body: body
}, cb)

{
  get:     methodNoBody('GET')     # (string, {[string]: string}, Response => T) => T
  options: methodNoBody('OPTIONS') # (string, {[string]: string}, Response => T) => T
  delete:  methodNoBody('DELETE')  # (string, {[string]: string}, Response => T) => T
  post:    methodBody('POST')      # (string, string, {[string]: string}, Response => T) => T
  put:     methodBody('PUT')       # (string, string, {[string]: string}, Response => T) => T
}
