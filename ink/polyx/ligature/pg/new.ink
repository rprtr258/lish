# /new

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
str := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/str.ink')
f := std.format

HeadTemplate := import('head.ink')

render := cb => (
  cb(Template())
)

Template := () => f('
{{ head }}

<body>
  <form action="/new" method="POST" class="noteEditForm">
    <header>
      <a href="/" class="title">&larr; ligature</a>
      <input type="submit" value="create" class="frost card block"/>
    </header>

    <div class="card">
      <div class="frost block light">label</div>
      <input type="text" name="label" class="paper block" placeholder="new-note" required autofocus>
    </div>
  </form>
</body>
', {
  head: HeadTemplate('new | ligature')
})

{render: render}
