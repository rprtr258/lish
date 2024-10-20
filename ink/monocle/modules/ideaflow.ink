` Module "ideaflow" indexes notes from my Ideaflow thoughtstream. Note that
Ideaflow notes have ontological graph-like properties that Monocle cannot index
faithfully, so we simply index the textual serialized representation.

The serialized thoughtstream can be obtained with the code:

const notes = this.noteRepository.getAll();
JSON.stringify(
  notes.map((note) => {
    const pmNode = noteToProsemirrorNode(note);
    const serialized = clipboardTextSerializer(pmNode.slice(0));
    return {id: note.id, content: serialized};
  }),
);
`

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
str := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/str.ink')
json := import('../vendor/json.ink')

log := std.log
slice := std.slice
cat := std.cat
map := std.map
each := std.each
filter := std.filter
append := std.append
readFile := std.readFile
split := str.split
trim := str.trim
replace := str.replace
hasPrefix? := str.hasPrefix?
trimPrefix := str.trimPrefix
deJSON := json.de

tokenizer := import('../lib/tokenizer')
tokenize := tokenizer.tokenize
tokenFrequencyMap := tokenizer.tokenFrequencyMap

Newline := char(10)

IdeaflowExportPath := '/tmp/ideaflow.json'

getDocs := withDocs => readFile(IdeaflowExportPath, file => file :: {
  () -> (
    log('[ideaflow] could not read ideaflow export file!')
    []
  )
  _ -> (
    notes := deJSON(file)
    docs := map(notes, (note, i) => (
      lines := split(note.content, Newline)
      title := (len(lines) > 1 :: {
        true -> lines.0
        _ -> ()
      })
      content := (len(lines) > 1 :: {
        true -> cat(slice(lines, 1, len(lines)), Newline)
        _ -> note.content
      })

      {
        id: 'if' + string(i)
        tokens: tokenize(note.content)
        content: content
        title: title
        href: 'https://ideaflow.app/?page=1&noteId=' + note.id
      }
    ))
    withDocs(docs)
  )
})

