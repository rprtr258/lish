package internal

import "testing"

func TestEval(t *testing.T) {
	ast := NewAst()
	n, err := parse_(ast, `2+2`)
	assertEqual(t, errParse{}, err)
	scope := NewEngine().CreateContext()
	k0 := func(v Value) ValueThunk {
		return func() Value {
			return v
		}
	}
	thunk := ast.Nodes[n].Eval(scope.Scope, ast, k0)
	thunk, ok := thunk().(ValueThunk)
	assertEqual(t, true, ok)
	th := thunk()
	_, ok = th.(ValueThunk)
	assertEqual(t, false, ok)
	assertEqual(t, Value(ValueNumber(4)), th)
}
