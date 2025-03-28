#!/usr/bin/env ink

# an http static file server
# with support for directory indexes

{log} := import('logging.ink')
{slice} := import('std.ink')
{format: f, join: cat} := import('str.ink')
{map, each} := import('functional.ink')
{readFile} := import('io.ink')

DIR := '.'
PORT := 7800
ALLOWINDEX := true

# short non-comprehensive list of MIME types
TYPES := {
  # text formats
  html: 'text/html; charset=utf-8'
  js: 'text/javascript; charset=utf-8'
  css: 'text/css; charset=utf-8'
  txt: 'text/plain; charset=utf-8'
  md: 'text/plain; charset=utf-8'
  # serve go & ink source code as plain text
  ink: 'text/plain; charset=utf-8'
  go: 'text/plain; charset=utf-8'

  # image formats
  jpg: 'image/jpeg'
  jpeg: 'image/jpeg'
  png: 'image/png'
  gif: 'image/gif'
  svg: 'image/svg+xml'

  # other misc
  pdf: 'application/pdf'
  zip: 'application/zip'
  json: 'application/json'
}

# given a path, get the file extension
getPathEnding := path => (
  (sub := (idx, acc) => true :: {
    idx == 0 -> path
    path.(idx) == '.' -> acc
    _ -> sub(idx - 1, path.(idx) + acc)
  })(len(path) - 1, '')
)

# given a path, get the MIME type
getType := path => (
  guess := TYPES.(getPathEnding(path))
  guess :: {
    () -> 'application/octet-stream'
    _ -> guess
  }
)

# prepare standard header
hdr := attrs => ( # TODO: just base + attrs ?
  base := {
    'X-Served-By': 'ink-serve'
    'Content-Type': 'text/plain'
  }
  each(keys(attrs), k => base.(k) := attrs.(k))
  base
)

# is this path a path to a directory?
dirPath? := path => path.(len(path) - 1) == '/'

# handles requests to validated paths
handleStat := (url, path, data, getElapsed) => data :: {
  # means file didn't exist
  () -> (
    # what if the path omits the .html extension?
    hpath := path + '.html'
    evt := stat(hpath)
    evt.type :: {
      'error' -> (
        log(f('  -> {{ url }} (.html) led to error in {{ ms }}ms: {{ error }}', {
          url
          ms: getElapsed()
          error: evt.message
        }))
        {
          status: 500
          headers: hdr({})
          body: 'server error'
        }
      )
      'data' -> evt.data :: {
        {dir: false, name: _, len: _} -> handlePath(url, hpath, getElapsed)
        _ -> (
          log(f('  -> {{ url }} not found in {{ ms }}ms', {
            url
            ms: getElapsed()
          }))
          {
            status: 404
            headers: hdr({})
            body: 'not found'
          }
        )
      }
    }
  )
  {dir: true, name: _, len: _, mod: _} -> true :: {
    dirPath?(path) -> handleDir(url, path, data, getElapsed)
    _ -> (
      log(f('  -> {{ url }} returned redirect to {{ url }}/ in {{ ms }}ms', {
        url
        ms: getElapsed()
      }))
      {
        status: 301
        headers: hdr({
          'Location': url + '/'
        })
        body: ''
      }
    )
  }
  {dir: false, name: _, len: _, mod: _} -> readFile(path, data => handleFileRead(url, path, data, getElapsed))
  _ -> {
    status: 500
    headers: hdr({})
    body: 'server invariant violation 1'
  }
}

# handles requests to readFile()
handleFileRead := (url, path, data, getElapsed) => data :: {
  () -> (
    log(f('  -> {{ url }} failed read in {{ ms }}ms', {
      url
      ms: getElapsed()
    }))
    {
      status: 500
      headers: hdr({})
      body: 'server error'
    }
  )
  _ -> (
    fileType := getType(path)
    log(f('  -> {{ url }} ({{ type }}) served in {{ ms }}ms', {
      url
      type: fileType
      ms: getElapsed()
    }))
    {
      status: 200
      headers: hdr({
        'Content-Type': getType(path)
      })
      body: data
    }
  )
}

# handle a directory we stat() confirmed to exist
handleExistingDir := (url, path, getElapsed) => true :: {
  ALLOWINDEX -> handleNoIndexDir(url, path, getElapsed)
  _ -> (
    log(f('  -> {{ url }} not allowed in {{ ms }}ms', {
      url
      ms: getElapsed()
    }))
    {
      status: 403
      headers: hdr({})
      body: 'permission denied'
    }
  )
}

# handles requests to directories '/'
handleDir := (url, path, data, getElapsed) => (
  ipath := path + 'index.html'
  evt := stat(ipath)
  evt.type :: {
    'error' -> (
      log(f('  -> {{ url }} (index) led to error in {{ ms }}ms: {{ error }}', {
        url
        ms: getElapsed()
        error: evt.message
      }))
      {
        status: 500
        headers: hdr({})
        body: 'server error'
      }
    )
    'data' -> evt.data :: {
      () -> handleExistingDir(url, path, getElapsed)
      # in the off chance that /index.html is a dir, just render index
      {dir: true, name: _, len: _, mod: _} -> handleExistingDir(url, path, getElapsed)
      {dir: false, name: _, len: _, mod: _} -> handlePath(url, ipath, getElapsed)
      _ -> {
        status: 500
        headers: hdr({})
        body: 'server invariant violation 2'
      }
    }
  }
)

# helpers for rendering the directory index page
makeIndex := (path, items) => '<title>' + path +
  '</title><style>body{font-family: system-ui,sans-serif}</style><h1>index of <code>' +
  path + '</code></h1><ul>' + items + '</ul>'
makeIndexLi := (fileStat, separator) => '<li><a href="' + fileStat.name + '" title="' + fileStat.name + '">' +
  fileStat.name + separator + ' (' + string(fileStat.len) + ' B)</a></li>'

# handles requests to dir() without /index.html
handleNoIndexDir := (url, path, getElapsed) => dir(path, evt => evt.type :: {
  'error' -> (
    log(f('  -> {{ url }} dir() led to error in {{ ms }}ms: {{ error }}', {
      url
      ms: getElapsed()
      error: evt.message
    }))
    {
      status: 500
      headers: hdr({})
      body: 'server error'
    }
  )
  'data' -> (
    log(f('  -> {{ url }} (index) served in {{ ms }}ms', {
      url
      ms: getElapsed()
    }))
    {
      status: 200
      headers: hdr({
        'Content-Type': 'text/html'
      })
      body: makeIndex(
        slice(path, 2, len(path))
        cat(map(evt.data, fileStat => makeIndexLi(
          fileStat
          true :: {
            fileStat.dir -> '/'
            _ -> ''
          }
        )), '')
      )
    }
  )
})

# trim query parameters
trimQP := path => (
  max := len(path)
  (sub := (idx, acc) => idx :: {
    max -> path
    _ -> path.(idx) :: {
      '?' -> acc
      _ -> sub(idx + 1, acc + path.(idx))
    }
  })(0, '')
)

# handles requests to path with given parameters
handlePath := (url, path, getElapsed) => (
  evt := stat(path)
  evt.type :: {
    'error' -> (
      log(f('  -> {{ url }} led to error in {{ ms }}ms: {{ error }}', {
        url
        ms: getElapsed()
        error: evt.message
      }))
      {
        status: 500
        headers: hdr({})
        body: 'server error'
      }
    )
    'data' -> handleStat(url, path, evt.data, getElapsed)
  }
)

# main server handler
listen('0.0.0.0:' + string(PORT), evt => evt.type :: {
  'error' -> log('server error: ' + evt.message)
  'req' -> (
    log(f('{{ method }}: {{ url }}', evt.data))

    # set up timer
    start := time()
    # trim the elapsed-time millisecond count at 2-3 decimal digits
    getElapsed := () => slice(string(floor((time() - start) * 1000000) / 1000), 0, 5)

    # normalize path
    url := trimQP(evt.data.url)

    # respond to file request
    evt.data.method :: {
      'GET' -> handlePath(url, DIR + url, getElapsed)
      _ -> (
        # if other methods, just drop the request
        log('  -> ' + evt.data.url + ' dropped')
        {
          status: 405
          headers: hdr({})
          body: 'method not allowed'
        }
      )
    }
  )
})