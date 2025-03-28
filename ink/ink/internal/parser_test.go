package internal

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// timeout
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("timeout")
		os.Exit(1)
	}()

	m.Run()
}

func TestParser(t *testing.T) {
	f := func(
		name string,
		source string,
		parser Parser[int],
	) {
		t.Run(name, func(t *testing.T) {
			ast := NewAstSlice()
			b, expr, err := parser(ast, []byte(source))
			require.Equal(t, errParse{}, err)
			require.Equal(t, 3, expr)
			require.Equal(t, []byte{}, b)
			require.Equal(t, []Node{}, ast.Nodes)
		})
	}

	f(
		`valid literal "log"`,
		`log`,
		parseIdentifier,
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
	)
}
