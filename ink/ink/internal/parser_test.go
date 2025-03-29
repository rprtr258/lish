package internal

import (
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
	L.Parse = true
	L.Lex = true

	m.Run()
}

func TestParser_error(t *testing.T) {
	ast := NewAstSlice()
	_, _, err := parseExpression(ast, []byte(`)`))
	require.NotEqual(t, errParse{}, err)
}

func TestParser(t *testing.T) {
	f := func(
		name string,
		source string,
		parser Parser[int],
		node Node,
	) {
		t.Run(name, func(t *testing.T) {
			ast := NewAstSlice()
			b, expr, err := parser(ast, []byte(source))
			require.Equal(t, errParse{}, err)
			assert.Equal(t, []byte{}, b)
			assert.IsType(t, node, ast.Nodes[expr])
			t.Log(ast.String())
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

  g('
')
)`,
		parseBlock,
		NodeExprList{},
	)
	return
	f(
		`valid lambda`,
		`str => (
  out(str)

  out('
')
)`,
		parseLambda,
		NodeLiteralFunction{},
	)
	f(
		"valid assignment",
		`log :=
(str => (
  out(str)

  out('
')
))`,
		parseAssignment,
		NodeExprBinary{},
	)
}
