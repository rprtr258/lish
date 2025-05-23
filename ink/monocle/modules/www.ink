` Module "www" indexes every post from my blog thesephist.com. `

std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
str := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/str.ink')

log := std.log
f := std.format
slice := std.slice
map := std.map
each := std.each
filter := std.filter
append := std.append
flatten := std.flatten
readFile := std.readFile
index := str.index
split := str.split
trim := str.trim

tokenizer := import('../lib/tokenizer.ink')
tokenize := tokenizer.tokenize
tokenFrequencyMap := tokenizer.tokenFrequencyMap

Newline := char(10)

ContentDir := env().HOME + '/src/www/content/posts'

getDocs := withDocs => dir(ContentDir, evt => evt.type :: {
  'error' -> (
    log('[www] could not read the posts directory')
    withDocs([])
  )
  'data' -> (
    entries := evt.data

    ` Post : [
      name: string
      content: string
    ]`
    posts := []

    ifAllRead := () => len(posts) :: {
      len(entries) -> (
        docs := map(posts, (post, i) => (
          log('[www] tokenizing post ' + post.name)
          {
            id: 'www' + string(i)
            tokens: tokenize(post.content)
            content: post.content
            ` NOTE: leaning on a lot of implicit assumptions about
            the way thesephist/www posts are formatted in Markdown
            front matter. `
            title: slice(
              post.content
              index(post.content, '"') + 1
              index(post.content, 'date: ') - 2
            )
            ` link is generated from the name of the Markdown file `
            href: f('https://thesephist.com/posts/{{ 0 }}/'
              [slice(post.name, 0, len(post.name) - 3)])
          }
        ))
        withDocs(docs)
      )
    }

    each(entries, entry => (
      readFile(ContentDir + '/' + entry.name, file => file :: {
        () -> (
          log('[www] could not read post' + entry.name)
          posts.len(posts) := {
            name: entry.name
            content: ''
          }
          ifAllRead()
        )
        _ -> (
          posts.len(posts) := {
            name: entry.name
            content: file
          }
          ifAllRead()
        )
      })
    ))
  )
})

