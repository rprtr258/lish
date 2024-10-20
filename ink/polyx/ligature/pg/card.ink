# note card component

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
f := std.format

note => (
  note.firstLine :: {
    '' -> note.firstLine := '(empty)'
  }
  f('
  <li>
    <a href="/note/{{ label }}" class="noteCard card" data-mod="{{ mod }}">
      <div class="paper block">
        {{ firstLine }}
      </div>
      <div class="noteMeta frost block light">
        <div>{{ label }}</div>
        <div class="modDate"></div>
      </div>
    </a>
  </li>
  ', note)
)
