#someday_maybe #project

https://github.com/rprtr258/lisp-sh

## Features

# processes pipes
something like graph declaration:
```
(pipe
    ("sorter" ("sort" "-u")
    "filter" ("rg" "500")
    "cutter" ("cut" "-d\" \"" "-f2-")) ; processes
    (("file" "logs.txt") ("stdin" "sorter")
    ("stdout" "sorter") ("stdin" "filter")
    ("stdout" "filter") ("stdin" "cutter")
    ("stdout" "cutter") "inherit"
    ("stderr" "sorter") ("file" "errors.txt")
    ("stderr" "filter") ("file" "errors.txt")
    ("stderr" "cutter") ("file" "errors.txt"))) ; pipes (aka edges)
```
vim is not a function, should be called like
```
(pipe
    (vim (vim $FILE$))
    ((stdin vim) (inherit stdin)
    (stdout vim) (inherit stdout)))
```
fd is either stdin/stdout/stderr of a process, or file, or /dev/null, or string
files on the left must exist, on the right can be created
https://askubuntu.com/questions/678915/whats-the-difference-between-and-in-bash
do a command to pipe stderr through automatically, sequential piping, more?
https://en.wikipedia.org/wiki/Plumber_(program)
https://stackoverflow.com/questions/70205473/pipe-between-two-child-processes
https://gist.github.com/mplewis/5279108

# terminate/pause/return child processes
https://stackoverflow.com/questions/49210815/how-do-i-send-a-signal-to-a-child-subprocess
https://github.com/rust-lang/rust/blob/b7c6e8f1805cd8a4b0a1c1f22f17a89e9e2cea23/src/libstd/sys/unix/process/process_unix.rs#L367

# common lisp stl
Depends on: namespaces
[CLISP - an ANSI Common Lisp](https://sourceforge.net/p/clisp/clisp/ci/tip/tree/src/init.lisp)

# default arguments
```lisp
; TODO: make default argument value, args to (&key (start 0) (end 1) (step 1))
; (range &start &end &step)
```
[argparse - Parser for command-line options, arguments and sub-commands - Python 3.10.2 documentation](https://docs.python.org/3/library/argparse.html)

# execute external programs and parse their output (s-exp/binary(asn?)) to add functions without modifying shell
[Canonical S-expressions - Wikipedia](https://en.wikipedia.org/wiki/Canonical_S-expressions)

# glob

# iterators / seqs / streams
Column: relation, transducer

# keyword args
```lisp
; equivalent
git commit :a :m ; (lish)
git commit -am   ; (bash)
; equivalent
ps :a :u :x   ; (lish)
ps ::aux      ; (lish)
ps -aux       ; (bash)
ps -a -u -x   ; (bash)
```

# lish make / goal oriented programming?
https://en.wikipedia.org/wiki/Make_(software)
https://www.gnu.org/software/make/manual/html_node/index.html
shell function returns `Either<Success, Error<status_code, stderr, stdout>>`
https://en.wikipedia.org/wiki/Icon_(programming_language)
[A short introduction to make](https://matt.might.net/articles/intro-to-make/)
[The Sentinel File Pattern](https://logicgrimoire.wordpress.com/2015/05/05/the-sentinel-file-pattern-3/) 
[makesure - make с человеческим лицом](https://habr.com/en/post/599137/)
https://ninja-build.org/manual.html

# metadata
Column: trace,untrace
```lisp
; TODO: rewrite with meta and set
(defmacro defun-trace (f args body)
  `(defun ,f ,args
    (let
      (res ,body)
      (echo "(" ',f " " (str ,@args) ") = " res)
      res)))
```

# namespaces
Column: common lisp stl

# / - paths joining operator
[PEP 428 -- The pathlib module -- object-oriented filesystem paths](https://www.python.org/dev/peps/pep-0428/)
```lisp
; эквивалентны
rm -rf "/" ;?
rm -rf $/
; эквивалентны
rm -rf "/home/abc/def" ; should strings be paths?
rm -rf $/home/abc/def
rm -rf (/ "home" "abc" "def")
rm -rf (/ 'home 'abc 'def) ;; symbol is elided to path?
```

# realize hashmap/list to filesystem 
[250bpm](https://250bpm.com/blog:153/index.html)

# relation programming
Depends on: iterators
```lisp
(->> /etc/passwd
  (select)                             ; select "/etc/passwd"
  (filter #(in (get % 1) to-delete)) ; where t.1 in to-delete
  (filter #(= (get % 1) (get % 2)))  ; where t.1 = t.2
  (map #(get % 1))                   ; select t.1
```
parallel map
[Relational shell programming](https://matt.might.net/articles/sql-in-the-shell/)
bench
[Command-line Tools can be 235x Faster than your Hadoop Cluster](https://adamdrake.com/command-line-tools-can-be-235x-faster-than-your-hadoop-cluster.html)

# rewrite some shell scripts in lish, .bashrc
rofi (other?) scripts
gtd scripts
[rprtr258/various-scripts](https://github.com/rprtr258/various-scripts/blob/master/min_books.py)

# tests
[mal/impls/tests at master · kanaka/mal](https://github.com/kanaka/mal/tree/master/impls/tests)

# trace/untrace
Depends on: metadata
[Is defun or setf preferred for creating function definitions in common lisp and why?](https://stackoverflow.com/a/42058202/15283421)

# transducers
Depends on: iterators

# документация для встроенных и определенных функций, макросов и специальных форм

# TODOS

# error reporting, not panicing

# quote-to-string
![](data/static/old/someday_maybe/programming_projects/lisp_shell/quotes.png)
```lisp
ssh a@a #'(ssh b@b #'(ssh c@c #'(ssh d@d #'(rm foo ; where #' is "quote-to-string" reader macro
(ssh a@a #'(ssh b@b #'(ssh c@c #'(ssh d@d #'(rm foo))))) ; is this
(as-> #'(rm foo) cmd
  #'(ssh d@d cmd)
  #'(ssh c@c cmd)
  #'(ssh b@b cmd)
  (ssh a@a cmd)) ; is something like this
```

# show fun/macro arguments in status bar
![](data/static/old/someday_maybe/programming_projects/lisp_shell/tui.png) 
https://www.google.com/search?q=ncurses%20rust
```lisp
;; args values suggestion
(defvariants :git-commands '(commit push status))
(def git (:git-command cmd) ...)
git <Tab> ; shows git-commands
git com<Tab> ; autocompletes to "git commit"
```


[Settling into Unix](https://matt.might.net/articles/settling-into-unix/)
[Shell programming with bash: by example, by counter-example](https://matt.might.net/articles/bash-by-example/)
[babashka/babashka](https://github.com/babashka/babashka)
[The Xonsh Shell - Python-powered shell](https://xon.sh/)
[The Janet Programming Language](https://janet-lang.org/)
[http://www.oilshell.org/](http://www.oilshell.org/)
[https://github.com/chr15m/flk](https://github.com/chr15m/flk)
[Presenting the Eshell](http://howardism.org/Technical/Emacs/eshell-present.html)
[Eshell](https://wikemacs.org/wiki/Eshell)
[A Description of Scsh](https://scsh.net/about/what.html)
[GitHub - metasoarous/oz: Data visualizations in Clojure and ClojureScript using Vega and Vega-lite](https://github.com/metasoarous/oz)
[Coroutines, exceptions, time-traveling search, generators and threads: Continuations by example](https://matt.might.net/articles/programming-with-continuations--exceptions-backtracking-search-threads-generators-coroutines/)

https://habr.com/en/post/131921/
https://habr.com/en/post/122029/
https://tkatchev.bitbucket.io/tab/docs.html
http://www.newlisp.org/newLISP_in_21_minutes.html
https://github.com/ngs-lang/ngs
https://build-your-own-x.vercel.app/#build-your-own-shell
https://github.com/janestreet/shexp
https://users.cs.northwestern.edu/~jesse/pubs/caml-shcaml/
