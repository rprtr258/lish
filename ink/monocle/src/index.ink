std := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')
str := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/str.ink')
json := import('../lib/rejson.ink')

log := std.log
f := std.format
scan := std.scan
map := std.map
each := std.each
append := std.append
readFile := std.readFile
writeFile := std.writeFile
split := str.split
deJSON := json.de
serJSON := json.ser

tokenizer := import('../lib/tokenizer.ink')
tokenize := tokenizer.tokenize
tokenFrequencyMap := tokenizer.tokenFrequencyMap

indexer := import('../lib/indexer.ink')
indexDocs := indexer.indexDocs

searcher := import('../lib/searcher.ink')
findDocs := searcher.findDocs

` modules `
Modules := {
  www: import('../modules/www.ink')
  entr: import('../modules/entr.ink')
  mira: import('../modules/mira.ink')
  tweets: import('../modules/tweets.ink')
  pocket: import('../modules/pocket.ink')
  lifelog: import('../modules/lifelog.ink')
  ligature: import('../modules/ligature.ink')
  ideaflow: import('../modules/ideaflow.ink')
}
ModuleState := {
  loadedModules: 0
}

` Doc : {
  id: string
  tokens: Map<string, number>
  content: string

  title: string | ()
  href: string | ()
} `

Docs := []

lazyGetDocs := (moduleKey, getDocs, withDocs) => readFile(f('./static/indexes/{{ 0 }}.json', [moduleKey]), file => file :: {
  () -> (
    log(f('[{{ 0 }}] re-generating index', [moduleKey]))
    getDocs(docs => writeFile(f('./static/indexes/{{ 0 }}.json', [moduleKey]), serJSON(docs), res => res :: {
      true -> withDocs(docs)
      _ -> (
        log('[main] failed to persist generated index for ' + moduleKey)
        withDocs(docs)
      )
    }))
  )
  _ -> withDocs(deJSON(file))
})

each(keys(Modules), moduleKey => (
  module := Modules.(moduleKey)
  getDocs := module.getDocs
  lazyGetDocs(moduleKey, getDocs, docs => (
    each(docs, doc => Docs.(doc.id) := (doc.module := moduleKey))

    ModuleState.loadedModules := ModuleState.loadedModules + 1
    ModuleState.loadedModules :: {
      len(Modules) -> (
        next := () => (
          log('[main] indexing docs...')
          Index := indexDocs(Docs)

          log('[main] persisting index...')
          writeFile('./static/indexes/index.json', serJSON(Index), res => res :: {
            true -> main(Index)
            _ -> (
              log('[main] failed to persist index!')
              main(Index)
            )
          })
        )

        log('[main] persisting docs...')
        writeFile('./static/indexes/docs.json', serJSON(Docs), res => res :: {
          true -> next()
          _ -> (
            log('[main] failed to persist docs!')
            next()
          )
        })

      )
    }
  ))
))

main := () => log('done.')

