#!/usr/bin/env ink

cli := import('../vendor/cli')

` sistine commands `
build := import('build').main
help := import('help').main

given := (cli.parsed)()
given.verb :: {
  'build' -> build()
  'help' -> help()
  _ -> build()
}

