package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/rprtr258/fun"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var L Logger

func init() {
	log.Logger = log.
		Output(zerolog.ConsoleWriter{
			Out: os.Stderr,
			FormatLevel: func(i any) string {
				s, _ := i.(string)
				bg := fun.Switch(s, "44").
					Case("42", zerolog.LevelInfoValue).
					Case("43", zerolog.LevelWarnValue).
					Case("41", zerolog.LevelErrorValue).
					End()

				return "\x1b[30;" + bg + "m " + strings.ToUpper(s) + " \x1b[0m"
			},
			PartsExclude: []string{zerolog.TimestampFieldName},
		})
}

type Logger struct {
	DumpFrame                 bool
	Lex, Parse, Dump, DumpAST bool
	// If FatalError is true, an error will halt the interpreter
	FatalError bool
}

func LogError(err *Err) {
	level := fun.IF(L.FatalError, zerolog.FatalLevel, zerolog.WarnLevel)
	for ee := err; ee != nil; ee = ee.Parent {
		defer log.WithLevel(level).
			Stringer("at", err.Pos).
			Stringer("kind", err.Reason).
			Msg(ee.Message)
	}
}

func LogScope(scope *Scope) {
	if !L.Dump {
		return
	}

	log.Debug().Stringer("scope", scope).Msg("frame dump")
}

func LogToken(tok Token) {
	if !L.Lex {
		return
	}

	e := log.Debug().
		Stringer("at", tok.Pos).
		Stringer("kind", tok.Kind)
	if tok.Str != "" {
		e = e.Str("str", tok.Str)
	}
	if tok.Num != 0 {
		e = e.Float64("f64", tok.Num)
	}
	e.Send()
}

func LogNode(node Node) {
	if !L.Parse {
		return
	}

	log.Debug().
		Stringer("at", node.Position()).
		Stringer("node", node).
		Send()
}

func LogAST(s *AST) {
	if !L.DumpAST {
		return
	}

	fmt.Print(s.String())
}
