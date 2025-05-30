package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/rs/zerolog/log"
	ink "github.com/thesephist/ink/internal"
)

const cliVersion = "0.1.9"

func usage() {
	fmt.Printf(`Ink is a minimal, powerful, functional programming language.
	ink v%s

By default, ink interprets from stdin.
	ink < main.ink
Run Ink programs from source files by passing it to the interpreter.
	ink main.ink
Start an interactive repl.
	ink
Run from the command line with -eval.
	ink -eval "f := () => out('hi'), f()"

`, cliVersion)
	flag.PrintDefaults()
}

func repl(ctx *ink.Context) {
	// add repl-specific builtins
	ctx.LoadFunc("clear", func(_ *ink.Context, _ ink.Pos, _ []ink.Value) ink.Value {
		fmt.Printf("\x1b[2J\x1b[H")
		return ink.Null
	})
	ctx.LoadFunc("dump", func(ctx *ink.Context, _ ink.Pos, _ []ink.Value) ink.Value {
		fmt.Println(ctx.Frame.String())
		return ink.Null
	})

	// run interactively in a repl
	reader := bufio.NewReader(os.Stdin)
	for {
		const greenArrow = "\x1b[32;1m>\x1b[0m "
		fmt.Print(greenArrow)

		text, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal().Err(err).Stringer("kind", ink.ErrSystem).Msg("unexpected end of input")
		}

		// we don't really care if expressions fail to eval
		// at the top level, user will see regardless, so drop err
		nodes, errP := ink.ParseReader(ctx.Engine.AST, "stdin", strings.NewReader(text))
		if errP != nil {
			ink.LogError(errP)
			continue
		}
		val := ctx.Eval(nodes) // TODO: fix not saving scope
		if err, ok := val.(ink.ValueError); ok {
			ink.LogError(err.Err)
		} else {
			fmt.Println(val.String())
		}
	}
}

const _pprof = false

func main() {
	if _pprof {
		f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	flag.Usage = usage

	// cli arguments
	verbose := flag.Bool("verbose", false, "Log all interpreter debug information")
	debugLexer := flag.Bool("debug-lex", false, "Log lexer output")
	debugParser := flag.Bool("debug-parse", false, "Log parser output")
	dump := flag.Bool("dump", false, "Dump global frame after eval")
	dumpAST := flag.Bool("dumpast", false, "Dump AST")

	version := flag.Bool("version", false, "Print version string and exit")
	help := flag.Bool("help", false, "Print help message and exit")

	eval := flag.String("eval", "", "Evaluate argument as an Ink program")

	flag.Parse()

	// if asked for version, disregard everything else
	switch {
	case *version:
		fmt.Println(cliVersion)
	case *help:
		flag.Usage()
	default:
		// collect all other non-parsed arguments from the CLI as files to be run
		args := flag.Args()

		ink.L = ink.Logger{
			DumpFrame:  *dump,
			Lex:        *verbose || *debugLexer,
			Parse:      *verbose || *debugParser,
			Dump:       *verbose || *dump,
			DumpAST:    *verbose || *dumpAST,
			FatalError: true,
		}

		stdin, _ := os.Stdin.Stat()
		eng := ink.NewEngine()
		ctx := eng.CreateContext()
		var node ink.NodeID
		switch {
		case *eval != "":
			var errP *ink.Err
			node, errP = ink.ParseReader(eng.AST, "eval", strings.NewReader(*eval))
			if errP != nil {
				ink.LogError(errP)
				return
			}
		case len(args) > 0:
			filePath := args[0]
			var err *ink.Err
			if node, err = ctx.ExecPath(filePath); err != nil {
				ink.LogError(err)
				log.Fatal().
					Stringer("kind", ink.ErrRuntime).
					Str("filepath", filePath).
					Msg("failed to execute file")
			}
		case stdin.Mode()&os.ModeCharDevice == 0:
			var errP *ink.Err
			node, errP = ink.ParseReader(eng.AST, "stdin", os.Stdin)
			if errP != nil {
				ink.LogError(errP)
				return
			}
		default:
			// if no files given and no stdin, default to repl
			ink.L.FatalError = false
			repl(ctx)
			eng.Listeners.Wait()
			return
		}

		// just run
		if err, ok := ctx.Eval(node).(ink.ValueError); ok {
			log.Fatal().Err(err.Err).Stringer("kind", ink.ErrRuntime).Msg("failed to execute")
		}
		eng.Listeners.Wait()
	}
}
