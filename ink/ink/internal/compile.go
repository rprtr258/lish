package internal

import (
	"fmt"
	"iter"
	"strconv"
	"strings"
)

type watWriter struct {
	strings.Builder
}

func (w *watWriter) LocalGet(s string) {
	w.WriteString("(local.get $")
	w.WriteString(s)
	w.WriteString(")")
}

func (w *watWriter) Call(s string) {
	w.WriteString("(call $")
	w.WriteString(s)
	w.WriteString(")")
}

func (w *watWriter) Param(s, typ string) {
	w.WriteString("(param $")
	w.WriteString(s)
	w.WriteString(" ")
	w.WriteString(typ)
	w.WriteString(")")
}

func (w *watWriter) Result(typ string) {
	w.WriteString("(result ")
	w.WriteString(typ)
	w.WriteString(")")
}

func compileFunc(name string, fn NodeLiteralFunction, w *watWriter) {
	w.WriteString("(func $")
	w.WriteString(name)
	for _, arg := range fn.arguments {
		switch arg := arg.(type) {
		case NodeIdentifier:
			w.Param(arg.val, "externref") // TODO: resolve type
		default:
			panic(fmt.Sprintf("unknown argument type: %T", arg))
		}
	}
	w.Result("externref") // TODO: resolve type
	compile(fn.body, w)
	w.WriteString(")")
}

func compile(n Node, w *watWriter) {
	switch n := n.(type) {
	case NodeLiteralNumber:
		w.WriteString("(f64.const ")
		w.WriteString(strconv.FormatFloat(n.val, 'f', -1, 64))
		w.WriteString(")")
	case NodeIdentifier:
		w.LocalGet(n.val)
	case NodeExprBinary:
		switch n.operator {
		case OpAdd:
			w.WriteString("(call $ink__plus ")
			compile(n.left, w)
			w.WriteString(" ")
			compile(n.right, w)
			w.WriteString(")")
		case OpLessThan:
			w.WriteString("(f64.lt ")
			compile(n.left, w)
			w.WriteString(" ")
			compile(n.right, w)
			w.WriteString(")")
		case OpMultiply:
			w.WriteString("(f64.mul ")
			compile(n.left, w)
			w.WriteString(" ")
			compile(n.right, w)
			w.WriteString(")")
		case OpSubtract:
			w.WriteString("(f64.sub ")
			compile(n.left, w)
			w.WriteString(" ")
			compile(n.right, w)
			w.WriteString(")")
		case OpDefine:
			switch lhs := n.left.(type) {
			case NodeIdentifier:
				switch rhs := n.right.(type) {
				case NodeLiteralFunction:
					compileFunc(lhs.val, rhs, w)
				default:
					panic(fmt.Sprintf("unknown rhs type: %T", rhs))
				}
			default:
				panic(fmt.Sprintf("unknown lhs type: %T", lhs))
			}
		default:
			panic(fmt.Sprintf("unknown binary operator: %s", n.operator))
		}
	case NodeFunctionCall:
		w.WriteString("(call $")
		switch fn := n.function.(type) {
		case NodeIdentifier:
			w.WriteString(fn.val)
		default:
			panic(fmt.Sprintf("unknown function type: %T", fn))
		}
		for _, arg := range n.arguments {
			w.WriteString(" ")
			compile(arg, w)
		}
		w.WriteString(")")
	case NodeLiteralObject:
		for _, entry := range n.entries {
			k, v := entry.key, entry.val
			switch k := k.(type) {
			case NodeIdentifier:
				switch v := v.(type) {
				case NodeLiteralFunction:
					compileFunc(k.val, v, w)
					w.WriteString("(export ")
					w.WriteString(strconv.Quote(k.val))
					w.WriteString(" (func $")
					w.WriteString(k.val)
					w.WriteString("))")
				case NodeIdentifier:
					w.WriteString("(export ")
					w.WriteString(strconv.Quote(k.val))
					w.WriteString(" (func $")
					w.WriteString(k.val)
					w.WriteString("))")
				default:
					panic(fmt.Sprintf("unknown value type: %T", v))
				}
			default:
				panic(fmt.Sprintf("unknown key type: %T", k))
			}
		}
	case NodeMatchExpr:
		switch condition := n.condition.(type) {
		case NodeLiteralBoolean:
			if !condition.val {
				panic("cant match on false")
			}
			if len(n.clauses) == 0 {
				panic("match on must have at least one clause")
			}

			for i, clause := range n.clauses {
				switch target := clause.target.(type) {
				case NodeIdentifierEmpty:
					w.WriteString(") (else")
					if i != len(n.clauses)-1 {
						panic("empty clause must be last")
					}
				default:
					compile(target, w)
					if i == 0 {
						w.WriteString("(if (result externref) (then ")
					} else {
						panic("not implemented")
					}
				}
				compile(clause.expression, w)
			}
			w.WriteString("))")
		default:
			panic(fmt.Sprintf("unknown match condition type: %T", n.condition))
		}
	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}

func Compile(nodes iter.Seq[Node]) string {
	var w watWriter
	w.WriteString("(module")
	for n := range nodes {
		fmt.Println(n.String())
		compile(n, &w)
	}
	w.WriteString(")")
	return w.String()
}
