` Module "tweets" indexes and makes searchable all of my tweets, from a Twitter
archive download. `

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

tokenizer := import('../lib/tokenizer.ink')
tokenize := tokenizer.tokenize
tokenFrequencyMap := tokenizer.tokenFrequencyMap

Newline := char(10)

TweetsFilePath := '/tmp/tweets.json'

getDocs := withDocs => readFile(TweetsFilePath, file => file :: {
  () -> (
    log('[tweets] could not read tweet archive export file!')
    []
  )
  _ -> (
    tweets := deJSON(file)
    docs := map(tweets, (tweet, i) => (
      i % 100 :: {
        0 -> log(string(i) + ' tweets tokenized...')
      }
      {
        id: 'tw' + string(i)
        tokens: tokenize(tweet.content)
        content: tweet.content
        href: 'https://twitter.com/thesephist/status/' + tweet.id
      }
    ))
    withDocs(docs)
  )
})

