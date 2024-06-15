package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	for name, tc := range map[string]struct {
		input string
		res   Atom
	}{
		"parse_nothing":           {"", atomNil},
		"parse_nothing_space":     {" ", atomNil},
		"num":                     {"1", atomList(atomInt(1))},
		"num_spaces":              {"   7   ", atomList(atomInt(7))},
		"negative_num":            {"-12", atomList(atomInt(-12))},
		"r#true":                  {"true", atomList(atomBool(true))},
		"r#false":                 {"false", atomList(atomBool(false))},
		"plus":                    {"+", atomList(atomSymbol("+"))},
		"minus":                   {"-", atomList(atomSymbol("-"))},
		"dash_abc":                {"-abc", atomList(atomSymbol("-abc"))},
		"dash_arrow":              {"->>", atomList(atomSymbol("->>"))},
		"abc":                     {"abc", atomList(atomSymbol("abc"))},
		"abc_spaces":              {"   abc   ", atomList(atomSymbol("abc"))},
		"abc5":                    {"abc5", atomList(atomSymbol("abc5"))},
		"abc_dash_def":            {"abc-def", atomList(atomSymbol("abc-def"))},
		"nil":                     {"()", atomList()},
		"nil_spaces":              {"(   )", atomList()},
		"set":                     {"(set a 2)", atomList(atomSymbol("set"), atomSymbol("a"), atomInt(2))},
		"list_nil":                {"(())", atomList(atomList())},
		"list_nil_2":              {"(()())", atomList(atomList(), atomList())},
		"list_list":               {"((3 4))", atomList(atomList(atomInt(3), atomInt(4)))},
		"list_inner":              {"(+ 1 (+ 3 4))", atomList(atomSymbol("+"), atomInt(1), atomList(atomSymbol("+"), atomInt(3), atomInt(4)))},
		"list_inner_spaces":       {"  ( +   1   (+   2 3   )   )  ", atomList(atomSymbol("+"), atomInt(1), atomList(atomSymbol("+"), atomInt(2), atomInt(3)))},
		"plus_expr":               {"(+ 1 2)", atomList(atomSymbol("+"), atomInt(1), atomInt(2))},
		"star_expr":               {"(* 1 2)", atomList(atomSymbol("*"), atomInt(1), atomInt(2))},
		"pow_expr":                {"(** 1 2)", atomList(atomSymbol("**"), atomInt(1), atomInt(2))},
		"star_negnum_expr":        {"(* -1 2)", atomList(atomSymbol("*"), atomInt(-1), atomInt(2))},
		"string_spaces":           {`   "abc"   `, atomList(atomString("abc"))},
		"quote_list":              {"'(a b c)", atomList(atomSymbol("quote"), atomList(atomSymbol("a"), atomSymbol("b"), atomSymbol("c")))},
		"quote_symbol":            {"'a", atomList(atomSymbol("quote"), atomSymbol("a"))},
		"unquote_symbol":          {"`(,a b)", atomList(atomSymbol("quasiquote"), atomList(atomList(atomSymbol("unquote"), atomSymbol("a")), atomSymbol("b")))},
		"comment":                 {"123 ; such number", atomList(atomInt(123))},
		"string_arg_l":            {`(load-file "compose.lish"`, atomList(atomSymbol("load-file"), atomString("compose.lish"))},
		"string_arg_r":            {`load-file "compose.lish")`, atomList(atomSymbol("load-file"), atomString("compose.lish"))},
		"right_outer_list_simple": {"(+ 1 2", atomList(atomSymbol("+"), atomInt(1), atomInt(2))},
		"outer_list_simple":       {`echo 92`, atomList(atomSymbol("echo"), atomInt(92))},
		"outer_plus":              {"+ 1 2", atomList(atomSymbol("+"), atomInt(1), atomInt(2))},
		"right_outer_twice":       {"(+ 1 2 (+ 3 4", atomList(atomSymbol("+"), atomInt(1), atomInt(2), atomList(atomSymbol("+"), atomInt(3), atomInt(4)))},
		"left_outer_twice":        {"+-curried 1) 3)", atomList(atomList(atomSymbol("+-curried"), atomInt(1)), atomInt(3))},
		"outer_left_outer":        {"+-curried 1) 3", atomList(atomList(atomSymbol("+-curried"), atomInt(1)), atomInt(3))},
		"outer_right_outer":       {"+ 1 2 (+ 3 4", atomList(atomSymbol("+"), atomInt(1), atomInt(2), atomList(atomSymbol("+"), atomInt(3), atomInt(4)))},
		"dict": {`{"a" 1 "b" "2"`, atomList(atomHash(map[string]Atom{
			"a": atomInt(1),
			"b": atomString("2"),
		}))},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.res, read(tc.input))
		})
	}
}
