`` Response := {
``   type: 'error'
``   message: string
`` } | {
``   type: 'resp'
``   data: {
``     status: number
``     headers: {[string]: string}
``     body: string
``   }
`` }

`` ({
``   method: string
``   url: string
``   headers: {[string]: string}
``   body: string
`` }, Response => T) => T
`` TODO: ({
``   method: string
``   url: string
``   headers: {[string]: string}
``   body: string
`` }) => Future(Response)
`` req := req

`` (string, {[string]: string}, Response => T) => T
get := (url, headers, cb) => req({
  method: 'GET'
  url: url
  headers: headers
}, cb)

`` (string, {[string]: string}, Response => T) => T
options := (url, headers, cb) => req({
  method: 'OPTIONS'
  url: url
  headers: headers
}, cb)

`` (string, string, {[string]: string}, Response => T) => T
post := (url, body, headers, cb) => req({
  method: 'POST'
  url: url
  body: body
  headers: headers
}, cb)

`` (string, string, {[string]: string}, Response => T) => T
put := (url, body, headers, cb) => req({
method: 'PUT'
url: url
body: body
headers: headers
}, cb)

`` (string, {[string]: string}, Response => T) => T
delete := (url, headers, cb) => req({
method: 'DELETE'
url: url
headers: headers
}, cb)