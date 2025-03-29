package internal

import (
	"fmt"
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

func compileFunc(w *watWriter, ast *AST, name string, fn NodeLiteralFunction) {
	w.WriteString("(func $")
	w.WriteString(name)
	for _, arg := range fn.Arguments {
		switch arg := ast.Nodes[arg].(type) {
		case NodeIdentifier:
			w.Param(arg.Val, "externref") // TODO: resolve type
		default:
			panic(fmt.Sprintf("unknown argument type: %T", arg))
		}
	}
	w.Result("externref") // TODO: resolve type
	compile(ast.Nodes[fn.Body], ast, w)
	w.WriteString(")")
}

func compile(n Node, ast *AST, w *watWriter) {
	switch n := n.(type) {
	case NodeLiteralNumber:
		w.WriteString("(f64.const ")
		w.WriteString(strconv.FormatFloat(n.Val, 'f', -1, 64))
		w.WriteString(")")
	case NodeIdentifier:
		w.LocalGet(n.Val)
	case NodeExprBinary:
		switch n.Operator {
		case OpAdd:
			w.WriteString("(call $ink__plus ")
			compile(ast.Nodes[n.Left], ast, w)
			w.WriteString(" ")
			compile(ast.Nodes[n.Right], ast, w)
			w.WriteString(")")
		case OpLessThan:
			w.WriteString("(f64.lt ")
			compile(ast.Nodes[n.Left], ast, w)
			w.WriteString(" ")
			compile(ast.Nodes[n.Right], ast, w)
			w.WriteString(")")
		case OpMultiply:
			w.WriteString("(f64.mul ")
			compile(ast.Nodes[n.Left], ast, w)
			w.WriteString(" ")
			compile(ast.Nodes[n.Right], ast, w)
			w.WriteString(")")
		case OpSubtract:
			w.WriteString("(f64.sub ")
			compile(ast.Nodes[n.Left], ast, w)
			w.WriteString(" ")
			compile(ast.Nodes[n.Right], ast, w)
			w.WriteString(")")
		case OpDefine:
			switch lhs := ast.Nodes[n.Left].(type) {
			case NodeIdentifier:
				switch rhs := ast.Nodes[n.Right].(type) {
				case NodeLiteralFunction:
					compileFunc(w, ast, lhs.Val, rhs)
				default:
					panic(fmt.Sprintf("unknown rhs type: %T", rhs))
				}
			default:
				panic(fmt.Sprintf("unknown lhs type: %T", lhs))
			}
		default:
			panic(fmt.Sprintf("unknown binary operator: %s", n.Operator))
		}
	case NodeFunctionCall:
		w.WriteString("(call $")
		switch fn := ast.Nodes[n.Function].(type) {
		case NodeIdentifier:
			w.WriteString(fn.Val)
		default:
			panic(fmt.Sprintf("unknown function type: %T", fn))
		}
		for _, arg := range n.Arguments {
			w.WriteString(" ")
			compile(ast.Nodes[arg], ast, w)
		}
		w.WriteString(")")
	case NodeLiteralObject:
		for _, entry := range n.Entries {
			k, v := entry.Key, entry.Val
			switch k := ast.Nodes[k].(type) {
			case NodeIdentifier:
				switch v := ast.Nodes[v].(type) {
				case NodeLiteralFunction:
					compileFunc(w, ast, k.Val, v)
					w.WriteString("(export ")
					w.WriteString(strconv.Quote(k.Val))
					w.WriteString(" (func $")
					w.WriteString(k.Val)
					w.WriteString("))")
				case NodeIdentifier:
					w.WriteString("(export ")
					w.WriteString(strconv.Quote(k.Val))
					w.WriteString(" (func $")
					w.WriteString(k.Val)
					w.WriteString("))")
				default:
					panic(fmt.Sprintf("unknown value type: %T", v))
				}
			default:
				panic(fmt.Sprintf("unknown key type: %T", k))
			}
		}
	case NodeExprMatch:
		switch condition := ast.Nodes[n.Condition].(type) {
		case NodeLiteralBoolean:
			if !condition.Val {
				panic("cant match on false")
			}
			if len(n.Clauses) == 0 {
				panic("match on must have at least one clause")
			}

			for i, clause := range n.Clauses {
				clauseNode := ast.Nodes[clause].(NodeMatchClause)
				switch target := ast.Nodes[clauseNode.Target].(type) {
				case NodeIdentifierEmpty:
					w.WriteString(") (else")
					if i != len(n.Clauses)-1 {
						panic("empty clause must be last")
					}
				default:
					compile(target, ast, w)
					if i == 0 {
						w.WriteString("(if (result externref) (then ")
					} else {
						panic("not implemented")
					}
				}
				compile(ast.Nodes[clauseNode.Expression], ast, w)
			}
			w.WriteString("))")
		default:
			panic(fmt.Sprintf("unknown match condition type: %T", n.Condition))
		}
	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}

func Compile(nodes []Node, ast *AST) string {
	var w watWriter
	w.WriteString("(module")
	for _, n := range nodes {
		fmt.Println(n.String())
		compile(n, ast, &w)
	}
	w.WriteString(")")
	return w.String()
}
