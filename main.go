package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/rprtr258/fun"
)

const HISTORY_FILE = ".lish_history"

// TODO: implement completer
type autocomplete struct{}

// impl Hinter for LishHelper {
//     type Hint = String;

//	    fn hint(&self, line: &str, pos: usize, ctx: &Context<'_>) -> Option<Self::Hint> {
//	        let history = ctx.history();
//	        let history_index = ctx.history_index();
//	        //println!("history: {:?}, index: {}", history.iter().collect::<Vec<&String>>(), history_index);
//	        //println!("{} at {}", line, pos);
//	        Some(history
//	            .get(history_index)
//	            .map(|s| s[line.len()..].to_owned())
//	            .unwrap_or("".to_owned())
//	        )
//	    }
//	}
func (self autocomplete) Do(line []rune, pos int) (newLine [][]rune, length int) {
	return [][]rune{[]rune("pokus"), []rune("fokus")}, 0
}

func run() error {
	editor, _ := readline.NewEx(&readline.Config{
		Prompt:       "=> ",
		HistoryFile:  HISTORY_FILE,
		AutoComplete: autocomplete{},
	})
	defer editor.Close()

	replEnv := newEnvRepl()
	cmdArgs := os.Args
	replEnv.set("*ARGV*", atomList(fun.Map[Atom](func(s string) Atom { return atomString(s) }, os.Args...)...))
	// TODO: rename to load ?
	// TODO: detect error
	rep(`(set load-file (fn (f) (eval (read (join "(progn\n" (slurp f) "\n)")))))`, replEnv)

	// execute file
	if len(cmdArgs) > 1 {
		rep(`(load-file "`+cmdArgs[1]+`")`, replEnv)
		return nil
	}

	// repl
	for {
		inputBuffer, err := editor.Readline()
		switch err {
		case nil:
			if inputBuffer == "" {
				continue
			}

			// editor.AddHistory(inputBuffer)
			if result := rep(inputBuffer, replEnv); result != "()" {
				fmt.Println(result)
			}
		case readline.ErrInterrupt:
			fmt.Println("CTRL-C")
			return nil
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
