#!/usr/bin/env ink

cli := import('../vendor/cli.ink')

# sistine commands
help := import('help.ink')

given := (cli.parsed)()
given.verb :: {
  'build' -> import('build.ink')()
  'help' -> help()
  _ -> help()
}
