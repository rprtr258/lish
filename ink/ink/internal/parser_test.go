package internal

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rprtr258/fun"
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
	s := NewAst()
	_, err := parse_(s, `)`)
	require.NotEqual(t, errParse{}, err)
}

func TestParser(t *testing.T) {
	t.SkipNow() // TODO: get back
	f := func(
		name string,
		source string,
		nodeKind NodeKind,
		check ...func(*AST, Node) bool,
	) {
		t.Run(name, func(t *testing.T) {
			ast := NewAst()
			expr, err := parse_(ast, source)
			t.Log(ast.String())
			require.Equal(t, errParse{}, err)
			require.Equal(t, nodeKind, ast.Nodes[expr].Kind)
			if len(check) > 0 {
				require.True(t, check[0](ast, ast.Nodes[expr]))
			}
		})
	}

	t.Run("many", func(t *testing.T) {
		f(
			`no traling comma`,
			`(a,b,c)`,
			NodeKindExprList,
		)
		f(
			`traling comma`,
			`(a,b,c,)`,
			NodeKindExprList,
		)
	})

	t.Run("identifier", func(t *testing.T) {
		f(
			`regular`,
			`log`,
			NodeKindIdentifier,
		)
		f(
			`predicate function`,
			`is_valid?`,
			NodeKindIdentifier,
		)
	})

	t.Run("expression", func(t *testing.T) {
		f(
			`valid expression identifier`,
			`log`,
			NodeKindIdentifier,
		)
		f(
			`valid symbols identifier`,
			`w?t_f!`,
			NodeKindIdentifier,
		)
		f(
			`valid function-call`,
			`out(str)`,
			NodeKindFunctionCall,
		)
		f(
			`valid literal-number`,
			`1`,
			NodeKindLiteralNumber,
		)
		f(
			`valid negative number`,
			`~1`,
			NodeKindExprUnary,
		)
		f(
			`valid addition`,
			`string(s) + ' '`,
			NodeKindExprBinary,
		)
		f(
			"valid block, empty",
			`()`,
			NodeKindExprList,
		)
	})

	t.Run("string", func(t *testing.T) {
		f(
			`regular`,
			`'
			'`,
			NodeKindLiteralString,
		)
		// f(
		// 	`backquoted`,
		// 	"``",
		// 	NodeLiteralString{},
		// )
		f(
			`with escaping`,
			`'a\n\'b'`,
			NodeKindLiteralString,
			func(a *AST, n Node) bool {
				assertEqual(t, "a\n'b", n.Meta.(string))
				return true
			},
		)
		f(
			`with escaping 2`,
			`'es"c \\a"pe
me'`,
			NodeKindLiteralString,
			func(a *AST, n Node) bool {
				assertEqual(t, `es"c \a"pe
me`, n.Meta.(string))
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
		NodeKindExprList,
	)

	f(
		`valid lambda, single inlined arg`,
		`str=>(
			out(str)
			out('\n')
		)`,
		NodeKindLiteralFunction,
		func(_ *AST, n Node) bool {
			return len(n.Children[1:]) == 1
		},
	)
	f(
		"valid lambda, zero args",
		`() => ()`,
		NodeKindLiteralFunction,
	)
	f(
		`valid lambda, two args`,
		`(a,b) => (a+b)`,
		NodeKindLiteralFunction,
		func(_ *AST, n Node) bool {
			return len(n.Children[1:]) == 2
		},
	)

	f(
		"accessor",
		`this.fields`,
		NodeKindExprBinary,
		func(ast *AST, n Node) bool {
			require.Equal(t, "this", ast.Nodes[n.Children[0]].Meta.(string))
			require.Equal(t, "fields", ast.Nodes[n.Children[1]].Meta.(string))
			return true
		},
	)
	f(
		"nested accessor",
		`this.fields.(len(this.fields))`,
		NodeKindExprBinary,
		func(ast *AST, n Node) bool {
			require.Equal(t, NodeKindExprBinary, n.Kind)
			// check n is (this.fields).len
			src := ast.Nodes[n.Children[0]]
			require.Equal(t, NodeKindExprBinary, src.Kind)
			require.Equal(t, "this", ast.Nodes[src.Children[0]].Meta.(string))
			require.Equal(t, "fields", ast.Nodes[src.Children[1]].Meta.(string))
			// require.Equal(t, "fields", ast.Nodes[src.Right].(NodeLiteralString).Val)
			return true
		},
	)
	f(
		"sub accessor",
		`(comp.list).(2).what`,
		NodeKindExprBinary,
		func(ast *AST, n Node) bool {
			complist2 := ast.Nodes[n.Children[0]]
			require.Equal(t, NodeKindExprBinary, complist2.Kind)
			complist := ast.Nodes[ast.Nodes[complist2.Children[0]].Children[0]]
			require.Equal(t, NodeKindExprBinary, complist.Kind)
			comp := ast.Nodes[complist.Children[0]]
			require.Equal(t, NodeKindIdentifier, comp.Kind)
			// list := ast.Nodes[complist.Right].(NodeLiteralString)
			list := ast.Nodes[complist.Children[1]]
			require.Equal(t, NodeKindIdentifier, list.Kind)
			// _2 := ast.Nodes[complist2.Right].(NodeLiteralNumber)
			_2 := ast.Nodes[ast.Nodes[complist2.Children[1]].Children[0]]
			require.Equal(t, NodeKindLiteralNumber, _2.Kind)
			// what := ast.Nodes[op.Right].(NodeLiteralString)
			what := ast.Nodes[n.Children[1]]
			require.Equal(t, NodeKindIdentifier, what.Kind)
			assertEqual(t, "comp", comp.Meta.(string))
			assertEqual(t, "list", list.Meta.(string))
			assertEqual(t, 2, _2.Meta.(float64))
			assertEqual(t, "what", what.Meta.(string))
			return true
		},
	)
	f(
		"array accessor",
		`arr.2`,
		NodeKindExprBinary,
		func(ast *AST, n Node) bool {
			require.Equal(t, NodeKindExprBinary, n.Kind)
			l := ast.Nodes[n.Children[0]]
			require.Equal(t, NodeKindIdentifier, l.Kind)
			r := ast.Nodes[n.Children[1]]
			require.Equal(t, NodeKindLiteralNumber, r.Kind)
			assertEqual(t, "arr", l.Meta.(string))
			assertEqual(t, 2, r.Meta.(float64))
			return true
		},
	)

	t.Run("assignment", func(t *testing.T) {
		f(
			"lambda rhs",
			`log := (str => (out(str)
				out('\n')
			))`,
			NodeKindExprBinary,
		)
		f(
			"lambda ignoring argument rhs",
			`f := _ => 1`,
			NodeKindExprBinary,
		)
		f(
			"lambda with assignment to acessor",
			`this.setName := name => this.name := name`,
			NodeKindExprBinary,
		)
		f(
			"valid assignment into dict destructure",
			`{a, b} := load('kal')`,
			NodeKindExprBinary,
		)
		f(
			"valid assignment into acessor",
			`xs.(i) := f(item, i)`,
			NodeKindExprBinary,
		)
		f(
			"valid assignment into function result acessor",
			`xs.(len(xs)) := 1`,
			NodeKindExprBinary,
		)
		f(
			"array element",
			`arr.2 := 'second'`,
			NodeKindExprBinary,
			func(ast *AST, n Node) bool {
				require.Equal(t, NodeKindExprBinary, n.Kind)
				r := ast.Nodes[n.Children[1]]
				assertEqual(t, NodeKindLiteralString, r.Kind)
				assertEqual(t, "second", r.Meta.(string))
				return true
			},
		)
		f(
			"two assignments with comment in between",
			`(
a := 1 # should yield a new copy
b := 1
)`,
			NodeKindExprList,
		)
	})

	f(
		"comment v2",
		"`aboba` 1",
		NodeKindLiteralNumber,
	)
	f(
		"valid match",
		`1 :: {
			1 -> 'hi'
			2 -> 'thing'
		}`,
		NodeKindExprMatch,
	)
	f(
		"valid list",
		`[5, 4, 3, 2, 1]`,
		NodeKindLiteralList,
	)
	f(
		"valid binary-op, accessor",
		`[5, 4, 3, 2, 1].2`,
		NodeKindExprBinary,
	)
	f(
		"_ == anything",
		`_ = len`,
		NodeKindExprBinary,
	)
	f(
		"negate expression",
		`~(1-2)`,
		NodeKindExprUnary,
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
		ast := NewAst()
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
		ast := NewAst()
		nodes, err := ParseReader(ast, "iife", strings.NewReader(`f:=(n,m)=>(n),(m=>f(1,m))(25)`))
		require.Nil(t, err)
		t.Log(ast.String())
		assertEqual(t, []Node{
			// NodeExprBinary{Pos{"iife", 1, 1}, OpDefine, 0, 5},
			// NodeFunctionCall{9, []int{11}},
			NodeExprList(Pos{"iife", 1, 1}, []NodeID{5, 11}),
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
		}, fun.Map[Node](func(n NodeID) Node { return ast.Nodes[n] }, ast.Nodes[nodes].Children...))
	})
}
