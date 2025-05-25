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
		nnn = NodeLiteralBoolean(nn.Position(*ast), bool(v))
	case ValueNumber:
		nnn = NodeLiteralNumber(nn.Position(*ast), float64(v))
	case ValueString:
		nnn = NodeLiteralString(nn.Position(*ast), string(*v.b))
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
	switch nn := ast.Nodes[n]; nn.Kind {
	case NodeKindLiteralBoolean, NodeKindLiteralNumber, NodeKindLiteralString, NodeKindIdentifierEmpty:
		return true
	case NodeKindIdentifier:
		return false
	case NodeKindLiteralFunction, NodeKindLiteralList, NodeKindFunctionCall:
		return constFoldList(ast, n, nn.Children)
	case NodeKindExprUnary:
		isConst := constantFolding(ast, nn.Children[0])
		if !isConst {
			return false
		}
		constEval(ast, n)
		return true
	case NodeKindExprBinary:
		l := constantFolding(ast, nn.Children[0])
		r := constantFolding(ast, nn.Children[1])
		isConst := l && r && nn.Meta.(Kind) != OpDefine
		if !isConst {
			return false
		}
		constEval(ast, n)
		return true
	case NodeKindExprList:
		isConst := true
		for _, n := range nn.Children {
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
	case NodeKindLiteralComposite:
		isConst := constFoldList(ast, n, nn.Children)
		if !isConst {
			return false
		}
		constEval(ast, n)
		return true
	case NodeKindExprMatch:
		isConst := constantFolding(ast, nn.Children[0])
		for _, clause := range nn.Children[1:] {
			constantFolding(ast, clause) // TODO: Target might be side effect(assignment) ????
		}
		// TODO: if one of constant-folded expressions matches AND ALL expressions before it are also constant,
		// just substitute with its target
		if !isConst {
			return false
		}
		constEval(ast, n)
		return true
	default:
		panic(fmt.Sprintf("cant optimise %s", ast.Nodes[n]))
	}
}

func listexprSimplifier(ast *AST, n NodeID) {
	switch nn := ast.Nodes[n]; nn.Kind {
	case NodeKindLiteralBoolean, NodeKindLiteralNumber, NodeKindLiteralString, NodeKindIdentifierEmpty, NodeKindIdentifier:
	case NodeKindLiteralFunction, NodeKindExprUnary:
		listexprSimplifier(ast, nn.Children[0])
	case NodeKindExprBinary:
		listexprSimplifier(ast, nn.Children[0])
		listexprSimplifier(ast, nn.Children[1])
	case NodeKindExprList:
		for _, n := range nn.Children {
			listexprSimplifier(ast, n)
		}
		if len(nn.Children) == 1 {
			ast.Nodes[n] = ast.Nodes[nn.Children[0]]
		}
	case NodeKindLiteralList, NodeKindLiteralComposite, NodeKindFunctionCall,
		NodeKindExprMatch: // TODO: target might be side effect(assignment) ????
		for _, n := range nn.Children {
			listexprSimplifier(ast, n)
		}
	default:
		panic(fmt.Sprintf("cant optimise %T", nn))
	}
}
