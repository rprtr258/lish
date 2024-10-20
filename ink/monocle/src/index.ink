std := import('../vendor/std')
str := import('../vendor/str')
json := import('../lib/rejson')

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

tokenizer := import('../lib/tokenizer')
tokenize := tokenizer.tokenize
tokenFrequencyMap := tokenizer.tokenFrequencyMap

indexer := import('../lib/indexer')
indexDocs := indexer.indexDocs

searcher := import('../lib/searcher')
findDocs := searcher.findDocs

` modules `
Modules := {
	www: import('../modules/www')
	entr: import('../modules/entr')
	mira: import('../modules/mira')
	tweets: import('../modules/tweets')
	pocket: import('../modules/pocket')
	lifelog: import('../modules/lifelog')
	ligature: import('../modules/ligature')
	ideaflow: import('../modules/ideaflow')
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

