package internal

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertEqual[T any](t *testing.T, want, got T) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("want/got -/+:\n" + diff)
		t.FailNow()
	}
}

func TestMain(m *testing.M) {
	// timeout
	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	fmt.Println("timeout")
	// 	os.Exit(1)
	// }()
	// L.Parse = true
	// L.Lex = true

	m.Run()
}

func TestParser_error(t *testing.T) {
	_, _, err := parseExpression(newCtx([]byte(`)`), ""))
	require.NotEqual(t, errParse{}, err)
}

func TestParser(t *testing.T) {
	f := func(
		name string,
		source string,
		parser Parser[int],
		node Node,
		check ...func(*AST, Node) bool,
	) {
		t.Run(name, func(t *testing.T) {
			c := newCtx([]byte(source), "")
			b, expr, err := parser(c)
			t.Log(c.ast.String())
			require.Equal(t, errParse{}, err)
			if b < len(source) {
				assertEqual(t, "", source[b:])
			}
			assert.IsType(t, node, c.ast.Nodes[expr])
			if len(check) > 0 {
				assert.True(t, check[0](c.ast, c.ast.Nodes[expr]))
			}
		})
	}

	t.Run("many", func(t *testing.T) {
		f(
			`no traling comma`,
			`(a,b,c)`,
			parseBlock,
			NodeExprList{},
		)
		f(
			`traling comma`,
			`(a,b,c,)`,
			parseBlock,
			NodeExprList{},
		)
	})

	t.Run("identifier", func(t *testing.T) {
		f(
			`regular`,
			`log`,
			parseIdentifier,
			NodeIdentifier{},
		)
		f(
			`predicate function`,
			`is_valid?`,
			parseIdentifier,
			NodeIdentifier{},
		)
	})

	t.Run("expression", func(t *testing.T) {
		f(
			`valid expression identifier`,
			`log`,
			parseExpression,
			NodeIdentifier{},
		)
		f(
			`valid symbols identifier`,
			`w?t_f!`,
			parseExpression,
			NodeIdentifier{},
		)
		f(
			`valid function-call`,
			`out(str)`,
			parseExpression,
			NodeFunctionCall{},
		)
		f(
			`valid literal-number`,
			`1`,
			parseExpression,
			NodeLiteralNumber{},
		)
		f(
			`valid negative number`,
			`~1`,
			parseExpression,
			NodeExprUnary{},
		)
		f(
			`valid addition`,
			`string(s) + ' '`,
			parseExpression,
			NodeExprBinary{},
		)
		f(
			"valid block, empty",
			`()`,
			parseExpression,
			NodeExprList{},
		)
	})

	t.Run("string", func(t *testing.T) {
		f(
			`regular`,
			`'
			'`,
			parseString,
			NodeLiteralString{},
		)
		f(
			`backquoted`,
			"``",
			parseString,
			NodeLiteralString{},
		)
		f(
			`with escaping`,
			`'a\n\'b'`,
			parseString,
			NodeLiteralString{},
			func(a *AST, n Node) bool {
				assertEqual(t, "a\n'b", n.(NodeLiteralString).Val)
				return true
			},
		)
		f(
			`with escaping 2`,
			`'es"c \\a"pe
me'`,
			parseString,
			NodeLiteralString{},
			func(a *AST, n Node) bool {
				assertEqual(t, `es"c \a"pe
me`, n.(NodeLiteralString).Val)
				return true
			},
		)
	})

	f(
		`valid block`,
		`(
			f(str)
			g('\n')
		)`,
		parseBlock,
		NodeExprList{},
	)

	f(
		`valid lambda, single inlined arg`,
		`str=>(
			out(str)
			out('\n')
		)`,
		parseExpression,
		NodeLiteralFunction{},
		func(_ *AST, n Node) bool {
			f := n.(NodeLiteralFunction)
			return len(f.Arguments) == 1
		},
	)
	f(
		"valid lambda, zero args",
		`() => ()`,
		parseExpression,
		NodeLiteralFunction{},
	)
	f(
		`valid lambda, two args`,
		`(a,b) => (a+b)`,
		parseExpression,
		NodeLiteralFunction{},
		func(_ *AST, n Node) bool {
			f := n.(NodeLiteralFunction)
			return len(f.Arguments) == 2
		},
	)

	f(
		"accessor",
		`this.fields`,
		parseExpression,
		NodeExprBinary{},
		func(ast *AST, n Node) bool {
			op := n.(NodeExprBinary)
			l := ast.Nodes[op.Left].(NodeIdentifier)
			r := ast.Nodes[op.Right].(NodeLiteralString)
			assert.Equal(t, "this", l.Val)
			assert.Equal(t, "fields", r.Val)
			return true
		},
	)
	f(
		"nested accessor",
		`this.fields.(len(this.fields))`,
		parseLhs,
		NodeExprBinary{},
		func(ast *AST, n Node) bool {
			op := n.(NodeExprBinary)
			// check n is (this.fields).len
			src := ast.Nodes[op.Left].(NodeExprBinary)
			assert.Equal(t, "this", ast.Nodes[src.Left].(NodeIdentifier).Val)
			assert.Equal(t, "fields", ast.Nodes[src.Right].(NodeLiteralString).Val)
			return true
		},
	)
	f(
		"sub accessor",
		`(comp.list).(2).what`,
		parseExpression,
		NodeExprBinary{},
		func(ast *AST, n Node) bool {
			op := n.(NodeExprBinary)
			complist2 := ast.Nodes[op.Left].(NodeExprBinary)
			complist := ast.Nodes[complist2.Left].(NodeExprBinary)
			comp := ast.Nodes[complist.Left].(NodeIdentifier)
			list := ast.Nodes[complist.Right].(NodeLiteralString)
			_2 := ast.Nodes[complist2.Right].(NodeLiteralNumber)
			what := ast.Nodes[op.Right].(NodeLiteralString)
			assertEqual(t, "comp", comp.Val)
			assertEqual(t, "list", list.Val)
			assertEqual(t, 2, _2.Val)
			assertEqual(t, "what", what.Val)
			return true
		},
	)
	f(
		"sub accessor",
		`(comp.list).(2).what`,
		parseLhs,
		NodeExprBinary{},
		func(ast *AST, n Node) bool {
			op := n.(NodeExprBinary)
			complist2 := ast.Nodes[op.Left].(NodeExprBinary)
			complist := ast.Nodes[complist2.Left].(NodeExprBinary)
			comp := ast.Nodes[complist.Left].(NodeIdentifier)
			list := ast.Nodes[complist.Right].(NodeLiteralString)
			_2 := ast.Nodes[complist2.Right].(NodeLiteralNumber)
			what := ast.Nodes[op.Right].(NodeLiteralString)
			assertEqual(t, "comp", comp.Val)
			assertEqual(t, "list", list.Val)
			assertEqual(t, 2, _2.Val)
			assertEqual(t, "what", what.Val)
			return true
		},
	)
	f(
		"array accessor",
		`arr.2`,
		parseLhs,
		NodeExprBinary{},
		func(ast *AST, n Node) bool {
			op := n.(NodeExprBinary)
			l := ast.Nodes[op.Left].(NodeIdentifier)
			r := ast.Nodes[op.Right].(NodeLiteralNumber)
			assertEqual(t, "arr", l.Val)
			assertEqual(t, 2, r.Val)
			return true
		},
	)

	t.Run("assignment", func(t *testing.T) {
		f(
			"lambda rhs",
			`log := (str => (out(str)
				out('\n')
			))`,
			parseAssignment,
			NodeExprBinary{},
		)
		f(
			"lambda ignoring argument rhs",
			`f := _ => 1`,
			parseAssignment,
			NodeExprBinary{},
		)
		f(
			"lambda with assignment to acessor",
			`this.setName := name => this.name := name`,
			parseAssignment,
			NodeExprBinary{},
		)
		f(
			"valid assignment into dict destructure",
			`{a, b} := load('kal')`,
			parseAssignment,
			NodeExprBinary{},
		)
		f(
			"valid assignment into acessor",
			`xs.(i) := f(item, i)`,
			parseAssignment,
			NodeExprBinary{},
		)
		f(
			"valid assignment into function result acessor",
			`xs.(len(xs)) := 1`,
			parseAssignment,
			NodeExprBinary{},
		)
		f(
			"array element",
			`arr.2 := 'second'`,
			parseAssignment,
			NodeExprBinary{},
			func(ast *AST, n Node) bool {
				op := n.(NodeExprBinary)
				r := ast.Nodes[op.Right].(NodeLiteralString)
				assertEqual(t, "second", r.Val)
				return true
			},
		)
		f(
			"two assignments with comment in between",
			`(
a := 1 # should yield a new copy
b := 1
)`,
			parseBlock,
			NodeExprList{},
		)
	})

	f(
		"comment v2",
		"`aboba` 1",
		parseExpression,
		NodeLiteralNumber{},
	)
	f(
		"valid match",
		`1 :: {
			1 -> 'hi'
			2 -> 'thing'
		}`,
		parseExpression,
		NodeExprMatch{},
	)
	f(
		"valid list",
		`[5, 4, 3, 2, 1]`,
		parseExpression,
		NodeLiteralList{},
	)
	f(
		"valid binary-op, accessor",
		`[5, 4, 3, 2, 1].2`,
		parseExpression,
		NodeExprBinary{},
	)
	f(
		"_ == anything",
		`_ = len`,
		parseExpression,
		NodeExprBinary{},
	)
	f(
		"negate expression",
		`~(1-2)`,
		parseExpression,
		NodeExprUnary{},
	)
}

//go:embed testdata/mangled.ink
var mangled string

func TestSkipSpaces(t *testing.T) {
	assertEqual(t, 5, skipSpaces([]byte("`aaa`"), false))
}

func TestParse(t *testing.T) {
	t.Run("mangled.ink", func(t *testing.T) {
		t.SkipNow()
		ast := NewAstSlice()
		nodes, err := ParseReader(ast, "testdata/mangled.ink", strings.NewReader(mangled))
		require.Nil(t, err)
		t.Log(ast.String())
		_ = nodes
		// require.Equal(t, []Node{
		// 	NodeExprBinary{Operator: 19, Left: 3, Right: 11},
		// 	NodeExprMatch{Condition: 13, Clauses: []int{15, 18}},
		// 	NodeLiteralFunction{Arguments: []int{}, Body: 20},
		// 	NodeExprBinary{Operator: 19, Left: 22, Right: 28},
		// 	NodeFunctionCall{Function: 3, Arguments: []int{22}},
		// 	NodeExprBinary{Operator: 19, Left: 31, Right: 36},
		// 	NodeExprBinary{Operator: 19, Left: 38, Right: 48},
		// 	NodeFunctionCall{Function: 3, Arguments: []int{51}},
		// 	NodeFunctionCall{Function: 3, Arguments: []int{56}},
		// 	NodeExprBinary{Operator: 19, Left: 3, Right: 63},
		// 	NodeExprBinary{Operator: 19, Left: 65, Right: 87},
		// 	NodeExprBinary{Operator: 19, Left: 89, Right: 100},
		// 	NodeFunctionCall{Function: 104, Arguments: []int{105}},
		// 	NodeFunctionCall{Function: 5, Arguments: []int{7}},
		// }, nodes)
	})

	t.Run("iife", func(t *testing.T) {
		ast := NewAstSlice()
		nodes, err := ParseReader(ast, "iife", strings.NewReader(`f:=(n,m)=>(n),(m=>f(1,m))(25)`))
		require.Nil(t, err)
		t.Log(ast.String())
		assertEqual(t, []Node{
			NodeExprBinary{Pos{"iife", 1, 1}, OpDefine, 0, 5},
			NodeFunctionCall{9, []int{11}},
			// /*  0 */ NodeIdentifierEmpty{},
			// /*  1 */ NodeLiteralBoolean{Val: false},
			// /*  2 */ NodeLiteralBoolean{Val: true},
			// /*  3 */ NodeIdentifier{Val: "f"},
			// /*  4 */ NodeIdentifier{Val: "n"},
			// /*  5 */ NodeIdentifier{Val: "m"},
			// /*  6 */ NodeExprList{Expressions: []int{4, 5}},
			// /*  7 */ NodeExprList{Expressions: []int{4}},
			// /*  8 */ NodeLiteralNumber{Val: 1},
			// /*  9 */ NodeFunctionCall{3, []int{8, 5}},
			// /* 10 */ NodeLiteralFunction{Arguments: []int{5}, Body: 9},
			// /* 11 */ NodeFunctionCall{7, []int{10}},
			// /* 12 */ NodeLiteralFunction{Arguments: []int{4, 5}, Body: 11},
			// /* 13 */ NodeExprBinary{Operator: 19, Left: 3, Right: 12},
			// /* 14 */ NodeLiteralNumber{Val: 25},
			// /* 15 */ NodeFunctionCall{13, []int{14}},
		}, nodes)
	})
}
