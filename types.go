package main

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
)

const (
	AtomKindBool   AtomKind = "bool"
	AtomKindInt    AtomKind = "int64"
	AtomKindFloat  AtomKind = "f64"
	AtomKindString AtomKind = "string"
	AtomKindError  AtomKind = "error"
	AtomKindHash   AtomKind = "hash"
	AtomKindStream AtomKind = "stream"
)

type Bool bool

func (s Bool) String() string   { return fmt.Sprint(bool(s)) }
func (s Bool) GoString() string { return fmt.Sprint(bool(s)) }
func (v Bool) Cmp(other Value) (int, bool) {
	va := bool(v)
	vb := bool(other.(Bool))
	res := 0
	switch {
	case !va && vb: // false < true
		res = 1
	case va && !vb: // true > false
		res = -1
	}
	return res, true
}

type Int int64

func (s Int) String() string              { return fmt.Sprint(int64(s)) }
func (s Int) GoString() string            { return fmt.Sprint(int64(s)) }
func (s Int) Cmp(other Value) (int, bool) { return cmp.Compare(s, other.(Int)), true }

type Float float64

func (s Float) String() string              { return fmt.Sprint(float64(s)) }
func (s Float) GoString() string            { return fmt.Sprint(float64(s)) }
func (s Float) Cmp(other Value) (int, bool) { return cmp.Compare(s, other.(Float)), true }

type String string

func (s String) String() string              { return string(s) }
func (s String) GoString() string            { return strconv.Quote(string(s)) }
func (s String) Cmp(other Value) (int, bool) { return cmp.Compare(s, other.(String)), true }

type Error string

func (s Error) String() string              { return "ERROR: " + strconv.Quote(string(s)) }
func (s Error) GoString() string            { return "ERROR: " + strconv.Quote(string(s)) }
func (s Error) Cmp(other Value) (int, bool) { return cmp.Compare(s, other.(Error)), true }

type Hash map[string]Atom

func (v Hash) String() string {
	items := make([]string, 0, len(v)*2)
	for k, v := range v {
		items = append(items, strconv.Quote(k), v.GoString())
	}
	return "{" + strings.Join(items, " ") + "}"
}
func (v Hash) GoString() string {
	items := make([]string, 0, len(v)*2)
	for k, v := range v {
		items = append(items, strconv.Quote(k), v.GoString())
	}
	return "{" + strings.Join(items, " ") + "}"
}
func (v Hash) Cmp(other Value) (int, bool) {
	va := v
	vb := other.(Hash)
	if len(va) != len(vb) {
		return 0, false
	}
	for k := range va {
		if !atomEq(va[k], vb[k]) {
			return 0, false
		}
	}
	return 0, true
}

type Stream chan Atom

func (s Stream) String() string              { return "#stream" }
func (v Stream) GoString() string            { return fmt.Sprintf("#stream@%v", v) }
func (s Stream) Cmp(other Value) (int, bool) { return 0, false }

func atomString[T ~string](s T) Atom {
	return Atom{AtomKindString, String(s)}
}

func atomBool[T ~bool](b T) Atom {
	return Atom{AtomKindBool, Bool(b)}
}

func atomHash(m map[string]Atom) Atom {
	return Atom{AtomKindHash, Hash(m)}
}

func atomInt[T interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}](n T) Atom {
	return Atom{AtomKindInt, Int(n)}
}

func atomFloat[T interface {
	float32 | float64
}](x T) Atom {
	return Atom{AtomKindFloat, Float(x)}
}

func lisherr(format string, args ...any) Atom {
	return Atom{AtomKindError, Error(fmt.Sprintf(format, args...))}
}

func atomCmp(a, b Atom) (int, bool) {
	if a.Kind != b.Kind {
		return 0, false
	}

	return a.Value.Cmp(b.Value)
}

func atomEq(a, b Atom) bool {
	cmp, ok := atomCmp(a, b)
	return ok && cmp == 0
}
