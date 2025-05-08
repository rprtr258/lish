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

var typeString = TypeString{}

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

var typeNull = TypeNull{}

type TypeVoid struct{}

func (TypeVoid) isType()        {}
func (TypeVoid) String() string { return "<void>" }

var typeVoid = TypeVoid{}

type TypeFunction struct {
	Args   []Type
	Result Type
}

func (TypeFunction) isType() {}
func (t TypeFunction) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	for i, arg := range t.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.String())
	}
	sb.WriteString(") -> ")
	sb.WriteString(t.Result.String())
	return sb.String()
}

type typeBinding struct {
	varname string
	ty      Type
}

type TypeContext []typeBinding

func (ctx TypeContext) String() string {
	var sb strings.Builder
	for _, decl := range ctx {
		fmt.Fprintf(&sb, "%s : %s\n", decl.varname, decl.ty.String())
	}
	return sb.String()
}

func (ctx TypeContext) Lookup(varname string) (Type, bool) {
	for _, decl := range ctx {
		if decl.varname == varname {
			return decl.ty, true
		}
	}
	return nil, false
}

func (ctx TypeContext) Append(varname string, ty Type) TypeContext {
	fmt.Println(varname, ":", ty.String())
	return append(ctx, typeBinding{varname, ty})
}

func typeIs[T Type](ty Type) bool {
	_, ok := ty.(T)
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

func typeInfer(ast *AST, n NodeID, ctx TypeContext) (Type, TypeContext) {
	switch n := ast.Nodes[n].(type) {
	case NodeLiteralBoolean:
		return typeBool, ctx
	case NodeLiteralNumber:
		return typeNumber, ctx
	case NodeLiteralString:
		return typeString, ctx
	case NodeIdentifier:
		varType, ok := ctx.Lookup(n.Val)
		if !ok {
			panicf("var %s is not defined", n.Val)
		}
		return varType, ctx
	case NodeLiteralFunction:
		// TODO: infer args types somehow or require to mark them explicitly
		ctxFn := ctx
		args := make([]Type, len(n.Arguments))
		for i, arg := range n.Arguments {
			args[i] = typeAny
			ctxFn = ctxFn.Append(ast.Nodes[arg].(NodeIdentifier).Val, typeAny)
		}
		typeResult, _ := typeInfer(ast, n.Body, ctxFn)
		return TypeFunction{args, typeResult}, ctx
	case NodeExprUnary:
		// NOTE: there is single unary operator ~ so type of unary expression is same as operand's
		return typeInfer(ast, n.Operand, ctx)
	case NodeExprBinary:
		rhs, _ := typeInfer(ast, n.Right, ctx)
		switch op := n.Operator; op {
		case OpAccessor:
			return typeAny, ctx
		case OpMultiply, OpSubtract, OpDivide: // number only operators
			lhs, _ := typeInfer(ast, n.Left, ctx)
			if !typeCheck(lhs, typeNumber) {
				typeCheckError("lhs of "+op.String(), typeNumber, lhs, ast.Nodes[n.Left].Position(ast))
			}
			if !typeCheck(rhs, typeNumber) {
				typeCheckError("rhs of "+op.String(), typeNumber, rhs, ast.Nodes[n.Right].Position(ast))
			}
			return typeNumber, ctx
		case OpDefine:
			rhs, _ := typeInfer(ast, n.Right, ctx)
			switch lhs := ast.Nodes[n.Left].(type) {
			case NodeIdentifier:
				return rhs, ctx.Append(lhs.Val, rhs)
			default:
				panicf("cant typecheck define operator with lhs %T", lhs)
			}
		default:
			panicf("cant typecheck binary operator %s", op.String())
		}
	case NodeExprList:
		if len(n.Expressions) == 0 {
			return typeNull, ctx
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
		typeFn, _ := typeInfer(ast, n.Function, ctx)
		args := make([]Type, len(n.Arguments))
		for i, arg := range n.Arguments {
			args[i], _ = typeInfer(ast, arg, ctx)
		}

		typeFunction, ok := typeFn.(TypeFunction)
		if !ok {
			typeCheckError("thing called as function", TypeFunction{}, typeFn, ast.Nodes[n.Function].Position(ast))
		}
		if len(typeFunction.Args) != len(n.Arguments) {
			panicf(
				"expected %d args of types %v but found %d of types %v",
				len(typeFunction.Args), typeFunction.Args,
				len(n.Arguments), args,
			)
		}
		for i, argExpected := range typeFunction.Args {
			argActual := args[i]
			if !typeCheck(argActual, argExpected) {
				typeCheckError(fmt.Sprintf("arg #%d", i), argExpected, argActual, ast.Nodes[n.Arguments[i]].Position(ast))
			}
		}

		return typeAny, ctx
	case NodeExprMatch:
		typeCond, _ := typeInfer(ast, n.Condition, ctx)
		typeResult := Type(typeVoid)
		for _, clause := range n.Clauses {
			if false { // TODO: implement pattern matching
				typeTarget, _ := typeInfer(ast, clause.Target, ctx)
				if !typeCheck(typeCond, typeTarget) {
					typeCheckError("target", typeTarget, typeCond, ast.Nodes[clause.Target].Position(ast))
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
