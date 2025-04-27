package internal

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func constEval(ast *AST, n NodeID) {
	nn := ast.Nodes[n]
	v := Value(Null) // TODO: nn.Eval(&Scope{}, ast)
	if isErr(v) {
		panic(v.String())
	}

	log.Printf("optimised %s to %v", nn.String(), v)

	var nnn Node
	switch v := v.(type) {
	case ValueBoolean:
		nnn = NodeLiteralBoolean{nn.Position(ast), bool(v)}
	case ValueNumber:
		nnn = NodeLiteralNumber{nn.Position(ast), float64(v)}
	case ValueString:
		nnn = NodeLiteralString{nn.Position(ast), string(*v.b)}
	case ValueNull, ValueComposite, ValueFunction:
		return
		// TODO: null, composite, list
	default:
		panic(fmt.Sprintf("unknown value %T", v))
	}
	ast.Nodes[n] = nnn
}

func constFoldList(ast *AST, n NodeID, nodes []NodeID) bool {
	isConst := true
	for _, n := range nodes {
		if !constantFolding(ast, n) {
			isConst = false
		}
	}
	if !isConst {
		return false
	}
	constEval(ast, n)
	return true
}

func constantFolding(ast *AST, n NodeID) bool {
	switch nn := ast.Nodes[n].(type) {
	case NodeLiteralBoolean, NodeLiteralNumber, NodeLiteralString, NodeIdentifierEmpty:
		return true
	case NodeIdentifier:
		return false
	case NodeLiteralFunction:
		return constFoldList(ast, n, append(nn.Arguments, nn.Body))
	case NodeExprUnary:
		isConst := constantFolding(ast, nn.Operand)
		if !isConst {
			return false
		}
		constEval(ast, n)
		return true
	case NodeExprBinary:
		l := constantFolding(ast, nn.Left)
		r := constantFolding(ast, nn.Right)
		isConst := l && r && nn.Operator != OpDefine
		if !isConst {
			return false
		}
		constEval(ast, n)
		return true
	case NodeExprList:
		isConst := true
		for _, n := range nn.Expressions {
			_ = n
			// TODO: get back
			// if !constantFolding(ast, n) {
			// 	isConst = false
			// }
		}
		if !isConst {
			return false
		}
		return true
	case NodeLiteralList:
		return constFoldList(ast, n, nn.Vals)
	case NodeLiteralComposite:
		isConst := true
		for _, n := range nn.Entries {
			if !constantFolding(ast, n.Key) {
				isConst = false
			}
			if !constantFolding(ast, n.Val) {
				isConst = false
			}
		}
		if !isConst {
			return false
		}
		constEval(ast, n)
		return true
	case NodeFunctionCall:
		return constFoldList(ast, n, append(nn.Arguments, nn.Function))
	case NodeExprMatch:
		isConst := constantFolding(ast, nn.Condition)
		for _, clause := range nn.Clauses {
			constantFolding(ast, clause.Expression) // TODO: might be side effect(assignment) ????
			constantFolding(ast, clause.Target)
		}
		// TODO: if one of constant-folded expressions matches AND ALL expressions before it are also constant,
		// just substitute with its target
		if !isConst {
			return false
		}
		constEval(ast, n)
		return true
	default:
		panic(fmt.Sprintf("cant optimise %T", nn))
	}
}

func listexprSimplifier(ast *AST, n NodeID) {
	switch nn := ast.Nodes[n].(type) {
	case NodeLiteralBoolean, NodeLiteralNumber, NodeLiteralString, NodeIdentifierEmpty, NodeIdentifier:
	case NodeLiteralFunction:
		listexprSimplifier(ast, nn.Body)
	case NodeExprUnary:
		listexprSimplifier(ast, nn.Operand)
	case NodeExprBinary:
		listexprSimplifier(ast, nn.Left)
		listexprSimplifier(ast, nn.Right)
	case NodeExprList:
		for _, n := range nn.Expressions {
			listexprSimplifier(ast, n)
		}
		if len(nn.Expressions) == 1 {
			ast.Nodes[n] = ast.Nodes[nn.Expressions[0]]
		}
	case NodeLiteralList:
		for _, n := range nn.Vals {
			listexprSimplifier(ast, n)
		}
	case NodeLiteralComposite:
		for _, n := range nn.Entries {
			listexprSimplifier(ast, n.Key)
			listexprSimplifier(ast, n.Val)
		}
	case NodeFunctionCall:
		for _, a := range nn.Arguments {
			listexprSimplifier(ast, a)
		}
		listexprSimplifier(ast, nn.Function)
	case NodeExprMatch:
		listexprSimplifier(ast, nn.Condition)
		for _, clause := range nn.Clauses {
			listexprSimplifier(ast, clause.Expression) // TODO: might be side effect(assignment) ????
			listexprSimplifier(ast, clause.Target)
		}
	default:
		panic(fmt.Sprintf("cant optimise %T", nn))
	}
}
