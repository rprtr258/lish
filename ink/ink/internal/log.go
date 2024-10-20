package internal

import (
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
	DumpFrame        bool
	Lex, Parse, Dump bool
	// If FatalError is true, an error will halt the interpreter
	FatalError bool
}

func LogError(err *Err) {
	e := log.Warn()
	if L.FatalError {
		e = log.Fatal()
	}

	e.Stringer("kind", err.reason).Msg(err.message)
}

func LogErr(ctx *Context, err *Err) {
	msg := err.message
	if ctx.File != "" {
		msg += " in " + ctx.File
	}

	LogError(&Err{err.reason, msg})
}

func LogScope(scope *Scope) {
	if L.Dump {
		log.Debug().Stringer("scope", scope).Msg("frame dump")
	}
}

func LogToken(tok Token) {
	if L.Lex {
		log.Debug().Stringer("token", tok).Msg("lex")
	}
}

func LogNode(node Node) {
	if L.Parse {
		log.Debug().Stringer("node", node).Msg("parse")
	}
}
