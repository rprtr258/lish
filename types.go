package main

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"

	"github.com/rprtr258/fun"
)

type Lambda struct {
	eval    func(ast Atom, env Env) Atom
	ast     Atom
	env     Env
	params  []string
	isMacro bool
	// meta Atom
}

type AtomKind int

const (
	AtomKindBool   AtomKind = iota // bool
	AtomKindInt                    // int64
	AtomKindFloat                  // f64
	AtomKindString                 // string
	AtomKindSymbol                 // string
	AtomKindError                  // string
	AtomKindHash                   // map[string]Atom
	AtomKindFunc                   // func([]Atom) Atom
	AtomKindLambda                 // Lambda
	AtomKindList                   // List, Nil if empty
)

func (a AtomKind) String() string {
	switch a {
	case AtomKindBool:
		return "bool"
	case AtomKindInt:
		return "int"
	case AtomKindFloat:
		return "float"
	case AtomKindString:
		return "string"
	case AtomKindSymbol:
		return "symbol"
	case AtomKindError:
		return "error"
	case AtomKindHash:
		return "hash"
	case AtomKindFunc:
		return "func"
	case AtomKindLambda:
		return "lambda"
	case AtomKindList:
		return "list"
	default:
		panic("unknown atom kind")
	}
}

type Atom struct {
	Kind  AtomKind
	Value any
}

var atomNil = Atom{AtomKindList, []Atom(nil)}

func atomSymbol(s string) Atom {
	return Atom{AtomKindSymbol, s}
}

func atomString(s string) Atom {
	return Atom{AtomKindString, s}
}

func atomBool(b bool) Atom {
	return Atom{AtomKindBool, b}
}

func atomHash(m map[string]Atom) Atom {
	return Atom{AtomKindHash, m}
}

func atomInt[T interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}](n T) Atom {
	return Atom{AtomKindInt, int64(n)}
}

func atomFloat[T interface {
	float32 | float64
}](x T) Atom {
	return Atom{AtomKindFloat, float64(x)}
}

func atomList(list ...Atom) Atom {
	if len(list) == 0 {
		return atomNil
	}

	return Atom{AtomKindList, list}
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

func atomFunc(fn func(...Atom) Atom, validators ...funcValidator) Atom {
	return Atom{AtomKindFunc, func(args []Atom) Atom {
		for _, v := range validators {
			if msg, ok := v(args); !ok {
				return lisherr("%s, but got %s", msg, strings.Join(fun.Map[string](Atom.String, args...), " "))
			}
		}

		return fn(args...)
	}}
}

func atomFuncNil(body func(...Atom), validators ...funcValidator) Atom {
	return atomFunc(func(args ...Atom) Atom {
		body(args...)
		return atomNil
	}, validators...)
}

func lisherr(format string, args ...any) Atom {
	return Atom{AtomKindError, fmt.Sprintf(format, args...)}
}

func (a Atom) String() string {
	switch a.Kind {
	case AtomKindBool:
		return fmt.Sprint(a.Value.(bool))
	case AtomKindInt:
		return fmt.Sprint(a.Value.(int64))
	case AtomKindFloat:
		return fmt.Sprint(a.Value.(float64))
	case AtomKindSymbol:
		return a.Value.(string)
	case AtomKindFunc:
		return "#fn"
	case AtomKindList:
		v := a.Value.([]Atom)
		if len(v) == 0 {
			return "()"
		}
		return "(" + strings.Join(fun.Map[string](Atom.String, v...), " ") + ")"
	case AtomKindHash:
		hash := a.Value.(map[string]Atom)
		items := make([]string, 0, len(hash)*2)
		for k, v := range hash {
			items = append(items, strconv.Quote(k), v.GoString())
		}
		return "{" + strings.Join(items, " ") + "}"
	case AtomKindLambda:
		v := a.Value.(Lambda)
		return fmt.Sprintf(
			"(%s %#v %s)",
			fun.IF(v.isMacro, "defmacro", "fn"),
			v.params,
			v.ast,
		)
	case AtomKindError:
		return "ERROR: " + a.Value.(string)
	case AtomKindString:
		return a.Value.(string)
	default:
		panic("unknown atom kind")
	}
}

func (a Atom) GoString() string {
	switch a.Kind {
	case AtomKindBool:
		return fmt.Sprint(a.Value.(bool))
	case AtomKindInt:
		return fmt.Sprint(a.Value.(int64))
	case AtomKindFloat:
		return fmt.Sprint(a.Value.(float64))
	case AtomKindSymbol:
		return a.Value.(string)
	case AtomKindFunc:
		return fmt.Sprintf("fn@%v", a.Value)
	case AtomKindList:
		v := a.Value.([]Atom)
		if len(v) == 0 {
			return "()"
		}
		return fmt.Sprint(len(v)) + "(" + strings.Join(fun.Map[string](Atom.GoString, v...), " ") + ")"
	case AtomKindHash:
		hash := a.Value.(map[string]Atom)
		items := make([]string, 0, len(hash)*2)
		for k, v := range hash {
			items = append(items, strconv.Quote(k), v.GoString())
		}
		return "{" + strings.Join(items, " ") + "}"
	case AtomKindLambda:
		v := a.Value.(Lambda)
		return fmt.Sprintf(
			"(%s %#v %#v)",
			fun.IF(v.isMacro, "defmacro", "fn"),
			v.params,
			v.ast,
		)
	case AtomKindError:
		return "ERROR: " + a.Value.(string)
	case AtomKindString:
		return fmt.Sprintf("%q", a.Value.(string))
	default:
		panic("unknown atom kind")
	}
}

func (a Atom) IsMacro() bool {
	return a.Kind == AtomKindLambda && a.Value.(Lambda).isMacro
}

func atomCmp(a, b Atom) (int, bool) {
	if a.Kind != b.Kind {
		return 0, false
	}

	switch a.Kind {
	case AtomKindBool:
		va := a.Value.(bool)
		vb := b.Value.(bool)
		res := 0
		switch {
		case !va && vb: // false < true
			res = 1
		case va && !vb: // true > false
			res = -1
		}
		return res, true
	case AtomKindInt:
		return cmp.Compare(a.Value.(int64), b.Value.(int64)), true
	case AtomKindString, AtomKindSymbol, AtomKindError:
		return cmp.Compare(a.Value.(string), b.Value.(string)), true
	case AtomKindList:
		va := a.Value.([]Atom)
		vb := b.Value.([]Atom)
		for i := 0; i < len(va); i++ {
			if c, ok := atomCmp(va[i], vb[i]); !ok {
				return 0, false
			} else if c != 0 {
				return c, true
			}
		}
		return 0, true
	default:
		return 0, false
	}
}

func atomEq(a, b Atom) bool {
	if a.Kind != b.Kind {
		return false
	}

	switch a.Kind {
	case AtomKindBool, AtomKindInt, AtomKindString, AtomKindSymbol, AtomKindError:
		return a.Value == b.Value
	case AtomKindList:
		va := a.Value.([]Atom)
		vb := b.Value.([]Atom)
		if len(va) != len(vb) {
			return false
		}
		for i := 0; i < len(va); i++ {
			if !atomEq(va[i], vb[i]) {
				return false
			}
		}
		return true
	case AtomKindHash:
		va := a.Value.(map[string]Atom)
		vb := b.Value.(map[string]Atom)
		if len(va) != len(vb) {
			return false
		}
		for k := range va {
			if !atomEq(va[k], vb[k]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
