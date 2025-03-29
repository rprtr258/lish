package internal

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	ast := NewAstSlice()
	_, _, err := parseExpression(ast, []byte(`)`))
	require.NotEqual(t, errParse{}, err)
}

//go:embed testdata/mangled.ink
var mangled string

func TestParser(t *testing.T) {
	f := func(
		name string,
		source string,
		parser Parser[int],
		node Node,
		check ...func(*AST, Node) bool,
	) {
		t.Run(name, func(t *testing.T) {
			ast := NewAstSlice()
			b, expr, err := parser(ast, []byte(source))
			t.Log(ast.String())
			require.Equal(t, errParse{}, err)
			assert.Equal(t, []byte{}, b)
			assert.IsType(t, node, ast.Nodes[expr])
			if len(check) > 0 {
				assert.True(t, check[0](ast, ast.Nodes[expr]))
			}
		})
	}

	f(
		`valid identifier`,
		`log`,
		parseIdentifier,
		NodeIdentifier{},
	)
	f(
		`valid expression identifier`,
		`log`,
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
		`valid literal-string`,
		`'
'`,
		parseString,
		NodeLiteralString{},
	)
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
		"valid block, empty",
		`()`,
		parseExpression,
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
		`valid iife after function definition`,
		`f:=(n,m)=>(n)
		(m=>f(1,m))(25)`,
		parseExpression,
		NodeFunctionCall{},
		func(ast *AST, n Node) bool {
			f := n.(NodeFunctionCall)
			assert.Equal(t, []Node{
				/*  0 */ NodeIdentifierEmpty{},
				/*  1 */ NodeLiteralBoolean{Val: false},
				/*  2 */ NodeLiteralBoolean{Val: true},
				/*  3 */ NodeIdentifier{Val: "f"},
				/*  4 */ NodeIdentifier{Val: "n"},
				/*  5 */ NodeIdentifier{Val: "m"},
				/*  6 */ NodeExprList{Expressions: []int{4, 5}},
				/*  7 */ NodeExprList{Expressions: []int{4}},
				/*  8 */ NodeLiteralNumber{Val: 1},
				/*  9 */ NodeFunctionCall{3, []int{8, 5}},
				/* 10 */ NodeLiteralFunction{Arguments: []int{5}, Body: 9},
				/* 11 */ NodeFunctionCall{7, []int{10}},
				/* 12 */ NodeLiteralFunction{Arguments: []int{4, 5}, Body: 11},
				/* 13 */ NodeExprBinary{Operator: 19, Left: 3, Right: 12},
				/* 14 */ NodeLiteralNumber{Val: 25},
				/* 15 */ NodeFunctionCall{13, []int{14}},
			}, ast.Nodes)
			assert.Equal(t, NodeFunctionCall{7, []int{10}}, f)
			return len(f.Arguments) == 1
		},
	)

	f(
		"valid assignment",
		`log := (str => (out(str)
			out('\n')
		))`,
		parseAssignment,
		NodeExprBinary{},
	)
	f(
		"valid match",
		`1 :: {
			1 -> 'hi',
			2 -> 'thing'
		}`,
		parseExpression,
		NodeExprMatch{},
	)
	f(
		"valid list",
		`[5 4 3 2 1]`,
		parseExpression,
		NodeLiteralList{},
	)
	f(
		"valid binary-op, accessor",
		`[5 4 3 2 1].2`,
		parseExpression,
		NodeExprBinary{},
	)
}

func TestParse(t *testing.T) {
	ast := NewAstSlice()
	nodes := ParseReader(ast, "testdata/mangled.ink", strings.NewReader(mangled))
	t.Log(ast.String())
	require.Equal(t, []Node{
		NodeExprBinary{Operator: 19, Left: 3, Right: 11},
		NodeExprMatch{Condition: 13, Clauses: []int{15, 18}},
		NodeLiteralFunction{Arguments: []int{}, Body: 21},
		NodeExprBinary{Operator: 19, Left: 23, Right: 31},
		NodeFunctionCall{Function: 3, Arguments: []int{23}},
		NodeExprBinary{Operator: 19, Left: 34, Right: 39},
		NodeExprBinary{Operator: 19, Left: 41, Right: 53},
		NodeFunctionCall{Function: 3, Arguments: []int{56}},
		NodeFunctionCall{Function: 3, Arguments: []int{61}},
		NodeExprBinary{Operator: 19, Left: 3, Right: 68},
		NodeExprBinary{Operator: 19, Left: 70, Right: 99},
		NodeFunctionCall{Function: 117, Arguments: []int{118}},
		NodeFunctionCall{Function: 5, Arguments: []int{7}},
	}, nodes)
}
