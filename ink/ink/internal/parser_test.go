package internal

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rprtr258/fun"
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

func parse_(ast *AST, source string) (NodeID, errParse) {
	node, err := ParseReader(ast, "", bytes.NewReader([]byte(source)))
	if err != nil {
		return -1, errParse{err}
	}
	return node, errParse{}
}

func TestParser_error(t *testing.T) {
	s := NewAstSlice()
	_, err := parse_(s, `)`)
	require.NotEqual(t, errParse{}, err)
}

func TestParser(t *testing.T) {
	t.SkipNow() // TODO: get back
	f := func(
		name string,
		source string,
		node Node,
		check ...func(*AST, Node) bool,
	) {
		t.Run(name, func(t *testing.T) {
			ast := NewAstSlice()
			expr, err := parse_(ast, source)
			t.Log(ast.String())
			require.Equal(t, errParse{}, err)
			assert.IsType(t, node, ast.Nodes[expr])
			if len(check) > 0 {
				assert.True(t, check[0](ast, ast.Nodes[expr]))
			}
		})
	}

	t.Run("many", func(t *testing.T) {
		f(
			`no traling comma`,
			`(a,b,c)`,
			NodeExprList{},
		)
		f(
			`traling comma`,
			`(a,b,c,)`,
			NodeExprList{},
		)
	})

	t.Run("identifier", func(t *testing.T) {
		f(
			`regular`,
			`log`,
			NodeIdentifier{},
		)
		f(
			`predicate function`,
			`is_valid?`,
			NodeIdentifier{},
		)
	})

	t.Run("expression", func(t *testing.T) {
		f(
			`valid expression identifier`,
			`log`,
			NodeIdentifier{},
		)
		f(
			`valid symbols identifier`,
			`w?t_f!`,
			NodeIdentifier{},
		)
		f(
			`valid function-call`,
			`out(str)`,
			NodeFunctionCall{},
		)
		f(
			`valid literal-number`,
			`1`,
			NodeLiteralNumber{},
		)
		f(
			`valid negative number`,
			`~1`,
			NodeExprUnary{},
		)
		f(
			`valid addition`,
			`string(s) + ' '`,
			NodeExprBinary{},
		)
		f(
			"valid block, empty",
			`()`,
			NodeExprList{},
		)
	})

	t.Run("string", func(t *testing.T) {
		f(
			`regular`,
			`'
			'`,
			NodeLiteralString{},
		)
		// f(
		// 	`backquoted`,
		// 	"``",
		// 	NodeLiteralString{},
		// )
		f(
			`with escaping`,
			`'a\n\'b'`,
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
		NodeExprList{},
	)

	f(
		`valid lambda, single inlined arg`,
		`str=>(
			out(str)
			out('\n')
		)`,
		NodeLiteralFunction{},
		func(_ *AST, n Node) bool {
			f := n.(NodeLiteralFunction)
			return len(f.Arguments) == 1
		},
	)
	f(
		"valid lambda, zero args",
		`() => ()`,
		NodeLiteralFunction{},
	)
	f(
		`valid lambda, two args`,
		`(a,b) => (a+b)`,
		NodeLiteralFunction{},
		func(_ *AST, n Node) bool {
			f := n.(NodeLiteralFunction)
			return len(f.Arguments) == 2
		},
	)

	f(
		"accessor",
		`this.fields`,
		NodeExprBinary{},
		func(ast *AST, n Node) bool {
			op := n.(NodeExprBinary)
			l := ast.Nodes[op.Left].(NodeIdentifier)
			r := ast.Nodes[op.Right].(NodeIdentifier)
			assert.Equal(t, "this", l.Val)
			assert.Equal(t, "fields", r.Val)
			return true
		},
	)
	f(
		"nested accessor",
		`this.fields.(len(this.fields))`,
		NodeExprBinary{},
		func(ast *AST, n Node) bool {
			op := n.(NodeExprBinary)
			// check n is (this.fields).len
			src := ast.Nodes[op.Left].(NodeExprBinary)
			assert.Equal(t, "this", ast.Nodes[src.Left].(NodeIdentifier).Val)
			assert.Equal(t, "fields", ast.Nodes[src.Right].(NodeIdentifier).Val)
			// assert.Equal(t, "fields", ast.Nodes[src.Right].(NodeLiteralString).Val)
			return true
		},
	)
	f(
		"sub accessor",
		`(comp.list).(2).what`,
		NodeExprBinary{},
		func(ast *AST, n Node) bool {
			op := n.(NodeExprBinary)
			complist2 := ast.Nodes[op.Left].(NodeExprBinary)
			complist := ast.Nodes[ast.Nodes[complist2.Left].(NodeExprList).Expressions[0]].(NodeExprBinary)
			comp := ast.Nodes[complist.Left].(NodeIdentifier)
			// list := ast.Nodes[complist.Right].(NodeLiteralString)
			list := ast.Nodes[complist.Right].(NodeIdentifier)
			// _2 := ast.Nodes[complist2.Right].(NodeLiteralNumber)
			_2 := ast.Nodes[ast.Nodes[complist2.Right].(NodeExprList).Expressions[0]].(NodeLiteralNumber)
			// what := ast.Nodes[op.Right].(NodeLiteralString)
			what := ast.Nodes[op.Right].(NodeIdentifier)
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
			NodeExprBinary{},
		)
		f(
			"lambda ignoring argument rhs",
			`f := _ => 1`,
			NodeExprBinary{},
		)
		f(
			"lambda with assignment to acessor",
			`this.setName := name => this.name := name`,
			NodeExprBinary{},
		)
		f(
			"valid assignment into dict destructure",
			`{a, b} := load('kal')`,
			NodeExprBinary{},
		)
		f(
			"valid assignment into acessor",
			`xs.(i) := f(item, i)`,
			NodeExprBinary{},
		)
		f(
			"valid assignment into function result acessor",
			`xs.(len(xs)) := 1`,
			NodeExprBinary{},
		)
		f(
			"array element",
			`arr.2 := 'second'`,
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
			NodeExprList{},
		)
	})

	f(
		"comment v2",
		"`aboba` 1",
		NodeLiteralNumber{},
	)
	f(
		"valid match",
		`1 :: {
			1 -> 'hi'
			2 -> 'thing'
		}`,
		NodeExprMatch{},
	)
	f(
		"valid list",
		`[5, 4, 3, 2, 1]`,
		NodeLiteralList{},
	)
	f(
		"valid binary-op, accessor",
		`[5, 4, 3, 2, 1].2`,
		NodeExprBinary{},
	)
	f(
		"_ == anything",
		`_ = len`,
		NodeExprBinary{},
	)
	f(
		"negate expression",
		`~(1-2)`,
		NodeExprUnary{},
	)
}

//go:embed testdata/mangled.ink
var mangled string

func TestSkipSpaces(t *testing.T) {
	// assertEqual(t, 5, skipSpaces([]byte("`aaa`"), false))
}

func TestParse(t *testing.T) {
	t.SkipNow()
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
			// NodeExprBinary{Pos{"iife", 1, 1}, OpDefine, 0, 5},
			// NodeFunctionCall{9, []int{11}},
			NodeExprList{Pos{"iife", 1, 1}, []NodeID{5, 11}},
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
		}, fun.Map[Node](func(n NodeID) Node { return ast.Nodes[n] }, ast.Nodes[nodes].(NodeExprList).Expressions...))
	})
}
