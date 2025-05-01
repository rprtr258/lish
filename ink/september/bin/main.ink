#!/usr/bin/env ink

{log, format: f, cat, each, readFile} := import('https://gist.githubusercontent.com/rprtr258/e208d8a04f3c9a22b79445d4e632fe98/raw/std.ink')

{parsed} := import('../vendor/cli.ink')

# september subcommands
highlight := import('../src/highlight.ink')
translate := import('../src/translate.ink')

PreamblePath := './runtime/ink.js'

commands := {
  # syntax-highlight input Ink programs from the token stream
  # and print them to stdout
  print: args => (
    files := args
    each(files, path => (
      readFile(path, data => out(highlight(data)))
    ))
  )
  # translate translates input Ink programs to JavaScript and
  # print them to stdout
  translate: args => (
    js := []
    files := args
    each(files, (path, i) => readFile(path, data => (
      js.(i) := translate(data)
      len(files) :: {
        len(js) -> log(cat(js, '\n'))
      }
    )))
  )
  'translate-full': args => readFile(PreamblePath, preamble => (
    js := [preamble]
    files := args
    each(files, (path, i) => readFile(path, data => (
      js.(i + 1) := translate(data)
      len(files) + 1 :: {
        len(js) -> log(cat(js, '\n'))
      }
    )))
  ))
  run: args => readFile(PreamblePath, preamble => (
    js := [preamble]
    files := args
    each(files, (path, i) => readFile(path, data => (
      js.(i + 1) := translate(data)
      len(files) + 1 :: {
        len(js) -> exec(
          'node'
          ['--']
          cat(js, ';' + '\n')
          evt => out(evt.data)
        )
      }
    )))
  ))
  # start an interactive REPL backed by Node.js, if installed.
  # might end up being the default behavior
  repl: args => log('command "repl" not implemented!')
}

{verb, args} := parsed()
verb :: {
  () -> log('September supports: \n\t' + cat(keys(commands), '\n\t'))
  _ -> commands.(verb) :: {
    () -> (
      log(f('command "{{ verb }}" not recognized', {verb}))
      log('September supports: \n\t' + cat(keys(commands), '\n\t'))
    )
    cmd -> cmd(args)
  }
}
