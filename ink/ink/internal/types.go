package internal

import (
	"fmt"
	"strings"
)

type Type interface {
	isType()
	String() string
}

type TypeString struct{}

func (TypeString) isType()        {}
func (TypeString) String() string { return "string" }

type TypeNumber struct{}

func (TypeNumber) isType()        {}
func (TypeNumber) String() string { return "number" }

var typeNumber = TypeNumber{}

type TypeBool struct{}

func (TypeBool) isType()        {}
func (TypeBool) String() string { return "bool" }

var typeBool = TypeBool{}

type TypeAny struct{}

func (TypeAny) isType()        {}
func (TypeAny) String() string { return "any" }

var typeAny = TypeAny{}

type TypeError struct{}

func (TypeError) isType()        {}
func (TypeError) String() string { return "error" }

type TypeNull struct{}

func (TypeNull) isType()        {}
func (TypeNull) String() string { return "()" }

type TypeVoid struct{}

func (TypeVoid) isType()        {}
func (TypeVoid) String() string { return "<void>" }

var typeVoid = TypeVoid{}

type TypeContext struct {
	prev    *TypeContext
	varname string
	typ     Type
}

func (ctx *TypeContext) String() string {
	var sb strings.Builder
	for c := ctx; c != nil; c = c.prev {
		fmt.Fprintf(&sb, "%s: %s\n", c.varname, c.typ.String())
	}
	return sb.String()
}

func (ctx *TypeContext) Lookup(varname string) (Type, bool) {
	switch {
	case ctx == nil:
		return nil, false
	case ctx.varname == varname:
		return ctx.typ, true
	case ctx.prev != nil:
		return ctx.Lookup(varname)
	default:
		return nil, false
	}
}

func typeIs[T Type](t Type) bool {
	_, ok := t.(T)
	return ok
}

func typeUnion(a, b Type) Type {
	switch {
	case typeIs[TypeVoid](a):
		return b
	case typeIs[TypeVoid](b):
		return a
	case typeCheck(a, b):
		return b
	case typeCheck(b, a):
		return a
	default:
		return typeAny
	}
}

func panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

// typeCheck checks that a can be assigned to b
func typeCheck(a, b Type) bool {
	switch b := b.(type) {
	case TypeAny:
		return true
	case TypeNumber:
		return typeIs[TypeNumber](a)
	case TypeString:
		return typeIs[TypeString](a)
	default:
		panicf("unknown type check case: %s : %s", a, b)
	}
	panic(1) // unreachable
}

func typeCheckError(
	thing string,
	typeExpected, typeActual Type,
	pos Pos,
) {
	panicf(
		"%s: %s is expected to be of type %s but it is %s",
		pos.String(), thing, typeExpected.String(), typeActual.String(),
	)
}

func typeInfer(ast *AST, n NodeID, ctx *TypeContext) (Type, *TypeContext) {
	switch n := ast.Nodes[n].(type) {
	case NodeLiteralBoolean:
		return typeBool, ctx
	case NodeLiteralNumber:
		return typeNumber, ctx
	case NodeLiteralString:
		return TypeString{}, ctx
	case NodeIdentifier:
		varType, ok := ctx.Lookup(n.Val)
		if !ok {
			panicf("var %s is not defined", n.Val)
		}
		return varType, ctx
	case NodeLiteralFunction:
		return typeAny, ctx
	case NodeExprUnary:
		// NOTE: there is single unary operator ~ so type of unary expression is same as operand's
		return typeInfer(ast, n.Operand, ctx)
	case NodeExprBinary:
		rhs, _ := typeInfer(ast, n.Right, ctx)
		switch op := n.Operator; op {
		case OpAccessor:
			return typeAny, ctx
		case OpMultiply:
			lhs, _ := typeInfer(ast, n.Left, ctx)
			if !typeCheck(lhs, typeNumber) {
				typeCheckError("lhs of *", typeNumber, lhs, ast.Nodes[n.Left].Position(ast))
			}
			if !typeCheck(rhs, typeNumber) {
				typeCheckError("rhs of *", typeNumber, rhs, ast.Nodes[n.Right].Position(ast))
			}
			return typeNumber, ctx
		case OpSubtract:
			lhs, _ := typeInfer(ast, n.Left, ctx)
			if !typeCheck(lhs, typeNumber) {
				typeCheckError("lhs of -", typeNumber, lhs, ast.Nodes[n.Left].Position(ast))
			}
			if !typeCheck(rhs, typeNumber) {
				typeCheckError("rhs of -", typeNumber, rhs, ast.Nodes[n.Right].Position(ast))
			}
			return typeNumber, ctx
		case OpDefine:
			rhs, _ := typeInfer(ast, n.Right, ctx)
			switch lhs := ast.Nodes[n.Left].(type) {
			case NodeIdentifier:
				ctx1 := &TypeContext{ctx, lhs.Val, rhs}
				return rhs, ctx1
			default:
				panicf("cant typecheck define operator with lhs %T", lhs)
			}
		default:
			panicf("cant typecheck binary operator %s", op.String())
		}
	case NodeExprList:
		if len(n.Expressions) == 0 {
			return TypeNull{}, ctx
		}

		ctxOrig := ctx
		ln := len(n.Expressions)
		for _, n := range n.Expressions[:ln-1] {
			//  type is not used, pass context to next expression
			_, ctx = typeInfer(ast, n, ctx)
		}
		typeResult, _ := typeInfer(ast, n.Expressions[ln-1], ctx)
		return typeResult, ctxOrig
	case NodeLiteralList:
		return typeAny, ctx
	case NodeLiteralComposite:
		return typeAny, ctx
	case NodeFunctionCall:
		return typeAny, ctx
	case NodeExprMatch:
		typeCond, _ := typeInfer(ast, n.Condition, ctx)
		typeResult := Type(typeVoid)
		for _, clause := range n.Clauses {
			if false { // TODO: implement pattern matching
				typeTarget, _ := typeInfer(ast, clause.Target, ctx)
				if !typeCheck(typeCond, typeTarget) {
					typeCheckError("target", typeTarget, typeCond, ast.Nodes[clause.Target].Position(ast))
					panicf("cond is of type %T but target is of type %T", typeCond, typeTarget)
				}
			}

			typeExpr, _ := typeInfer(ast, clause.Expression, ctx)
			typeResult = typeUnion(typeResult, typeExpr)
		}
		return typeResult, ctx
	default:
		panicf("cant typecheck %T", n)
	}
	panic(1) // unreachable
}
