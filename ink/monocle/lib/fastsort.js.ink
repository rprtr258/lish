` fast version of quicksort used in the UI code to sort >1000 docs quickly in
the ranker, using JavaScript's native array sort. `

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')

` hacky check to see if we're in a browser. Because the browser environment's
load() returns window / globalThis, std.location will not be null; otherwise,
in the native environment std.location will be (). `
available? := ~(std.location = ())

` assumes pred() itself is computationally cheap `
fastSortBy := (v, pred) => (
  ` coerce v into a JavaScript array. `
  v := slice(v, 0, len(v))

  preds := {}
  each(v, doc => preds.(doc.id) := pred(doc))
  bind(v, 'sort')((a, b) => preds.(a.id) - preds.(b.id))
)

