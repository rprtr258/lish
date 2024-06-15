package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	for name, tc := range map[string]struct {
		ast Atom
		res string
	}{
		"print_true":       {atomBool(true), "true"},
		"print_false":      {atomBool(false), "false"},
		"print_float":      {atomFloat(3.14), "3.14"},
		"print_int":        {atomInt(92), "92"},
		"print_empty_list": {atomNil, "()"},
		"print_list":       {atomList(atomInt(1), atomInt(2)), "(1 2)"},
		"print_symbol":     {atomSymbol("abc"), "abc"},
		"test_print_nice":  {atomString("\n"), "\n"},
		"test_print_dict": {atomHash(map[string]Atom{
			"a": atomInt(1),
			"b": atomString("2"),
		}), `{"a" 1 "b" "2"}`},
		"print_func": {atomFunc(func(x ...Atom) Atom { return x[0] }), "#fn"},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.ast.String(), tc.res)
		})
	}
}

func TestPrintDebug(t *testing.T) {
	for name, tc := range map[string]struct {
		ast Atom
		res string
	}{
		"print_nil":                  {atomNil, "()"},
		"print_string":               {atomString("abc"), `"abc"`},
		"print_string_with_slash":    {atomString(`\`), `"\\"`},
		"print_string_with_2slashes": {atomString(`\\`), `"\\\\"`},
		"print_string_with_newline":  {atomString("\n"), `"\n"`},
		"test_print_debug_dict": {atomHash(map[string]Atom{
			"a": atomInt(1),
			"b": atomString("2"),
		}), `{"a" 1 "b" "2"}`},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.ast.GoString(), tc.res)
		})
	}
}
