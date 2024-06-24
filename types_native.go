package main

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/rprtr258/fun"
)

type AtomKind string

const (
	AtomKindSymbol AtomKind = "symbol"
	AtomKindFunc   AtomKind = "fn"
	AtomKindLambda AtomKind = "lambda"
	AtomKindList   AtomKind = "list"
)

type Symbol string

func (s Symbol) String() string   { return string(s) }
func (s Symbol) GoString() string { return string(s) }
func (s Symbol) Cmp(other Value) (int, bool) {
	return cmp.Compare(s, other.(Symbol)), true
}

type Func func([]Atom) Atom

func (s Func) String() string        { return "#fn" }
func (s Func) GoString() string      { return fmt.Sprintf("fn@%v", s) }
func (s Func) Cmp(Value) (int, bool) { return 0, false }

type Lambda struct {
	eval    func(ast Atom, env Env) Atom
	ast     Atom
	env     Env
	params  []Symbol
	isMacro bool
	// meta Atom
}

func (v Lambda) String() string {
	return fmt.Sprintf(
		"(%s %#v %s)",
		fun.IF(v.isMacro, "defmacro", "fn"),
		v.params,
		v.ast,
	)
}
func (v Lambda) GoString() string {
	return fmt.Sprintf(
		"(%s %#v %#v)",
		fun.IF(v.isMacro, "defmacro", "fn"),
		v.params,
		v.ast,
	)
}
func (v Lambda) Cmp(Value) (int, bool) { return 0, false }

type List []Atom // List, Nil if empty
func (v List) String() string {
	if len(v) == 0 {
		return "()"
	}
	return "(" + strings.Join(fun.Map[string](Atom.String, v...), " ") + ")"
}
func (v List) GoString() string {
	if len(v) == 0 {
		return "()"
	}
	return fmt.Sprint(len(v)) + "(" + strings.Join(fun.Map[string](Atom.GoString, v...), " ") + ")"
}
func (va List) Cmp(other Value) (int, bool) {
	vb := other.(List)
	for i := 0; i < len(va); i++ {
		if c, ok := atomCmp(va[i], vb[i]); !ok {
			return 0, false
		} else if c != 0 {
			return c, true
		}
	}
	return 0, true
}

type Value interface {
	fmt.Stringer
	fmt.GoStringer
	// compare, given two atoms of same kind
	// return -1 if less, 0 if equal, 1 if greater
	// false if not comparable
	Cmp(Value) (int, bool)
}

type Atom struct {
	Kind  AtomKind
	Value Value
}

var atomNil = Atom{AtomKindList, List(nil)}

func (a Atom) String() string   { return a.Value.String() }
func (a Atom) GoString() string { return a.Value.GoString() }

func atomSymbol(s string) Atom {
	return Atom{AtomKindSymbol, Symbol(s)}
}

func atomList(list ...Atom) Atom {
	if len(list) == 0 {
		return atomNil
	}

	return Atom{AtomKindList, List(list)}
}

type funcValidator = func([]Atom) (string, bool)

func validateMinArgs(n int) funcValidator {
	return func(args []Atom) (string, bool) {
		if len(args) < n {
			return fmt.Sprintf("Expected at least %d arguments", n), false
		}
		return "", true
	}
}

func validateExactArgs(n int) funcValidator {
	return func(args []Atom) (string, bool) {
		if len(args) != n {
			return fmt.Sprintf("Expected exactly %d arguments", n), false
		}
		return "", true
	}
}

func validateArgsOfKind(kind AtomKind) funcValidator {
	return func(args []Atom) (string, bool) {
		for i, arg := range args {
			if arg.Kind != kind {
				return fmt.Sprintf(
					"Expected all arguments to be %s, but %d-th argument is %s",
					kind, i, arg,
				), false
			}
		}
		return "", true
	}
}

func atomLambda(lambda Lambda) Atom {
	return Atom{AtomKindLambda, lambda}
}

func atomFunc(
	fn func(...Atom) Atom,
	validators ...funcValidator,
) Atom {
	return Atom{AtomKindFunc, Func(func(args []Atom) Atom {
		for _, v := range validators {
			if msg, ok := v(args); !ok {
				return lisherr("%s, but got %s", msg, strings.Join(fun.Map[string](Atom.String, args...), " "))
			}
		}

		return fn(args...)
	})}
}

func atomFuncNil(
	body func(...Atom),
	validators ...funcValidator,
) Atom {
	return atomFunc(func(args ...Atom) Atom {
		body(args...)
		return atomNil
	}, validators...)
}
