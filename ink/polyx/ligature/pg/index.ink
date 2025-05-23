` /index.html `

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
f := std.format
each := std.each
map := std.map
cat := std.cat
readFile := std.readFile
str := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/str.ink')
split := str.split

quicksort := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/quicksort.ink')
sortBy := quicksort.sortBy

HeadTemplate := import('head.ink')
NoteCard := import('card.ink')

render := (dbPath, cb) => dir(dbPath, evt => evt.type :: {
  'error' -> cb('error finding notes')
  _ -> (
    ` notes are sorted by date last modified (reverse chron) `
    sortBy(evt.data, fstat => ~(fstat.mod))

    notes := map(evt.data, fileInfo => {
      label: split(fileInfo.name, '.').0
      mod: fileInfo.mod
      firstLine: '...?'
    })

    s := {
      count: 0
      total: len(notes)
    }
    len(notes) :: {
      0 -> cb(Template(notes))
    }
    each(notes, n => readFile(dbPath + '/' + n.label + '.md', file => (
      file :: {
        () -> n.firstLine := 'error reading...'
        ` hand-rolled efficient code to only trim file to
          the first line (first newline character) `
        _ -> n.firstLine := (sub := (acc, i) => file.(i) :: {
          () -> acc
          char(10) -> acc
          _ -> (
            acc + file.(i)
            sub(acc + file.(i), i + 1)
          )
        })('', 0)
      }

      s.count := s.count + 1
      s.count :: {
        s.total -> cb(Template(notes))
      }
    )))
  )
})

Template := notes => f('
{{ head }}

<body>
  <header>
    <a href="/" class="title">ligature</a>
    <form action="/find" method="GET" class="searchBar card">
      <input type="text" name="q" placeholder="search..." class="searchInput paper block" autofocus/>
      <input type="submit" value="find" class="frost block"/>
    </form>
    <a href="/new" class="newButton frost card block">new</a>
  </header>

  <ul class="noteList">
    {{ noteCards }}
  </ul>
  <script src="/static/js/ligature.js"></script>
</body>
', {
  head: HeadTemplate('ligature')
  noteCards: cat(map(notes, NoteCard), '')
})

{render: render}
