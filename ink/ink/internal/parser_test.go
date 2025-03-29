package internal

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// timeout
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("timeout")
		os.Exit(1)
	}()
	L.Parse = true

	m.Run()
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
			assert.Equal(t, node, ast.Nodes[expr])
		})
	}

	f(
		`valid literal "log"`,
		`log`,
		parseIdentifier,
		NodeIdentifier{Val: "log"},
	)
	f(
		`valid expression literal`,
		`log`,
		parseExpression,
		NodeIdentifier{Val: "log"},
	)
	return
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
