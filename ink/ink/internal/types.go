package internal

import (
	"fmt"
	"strings"
)

type Type interface {
	isType()
	String(ctx TypeContext) string
}

type TypeString struct{}

func (TypeString) isType()                   {}
func (TypeString) String(TypeContext) string { return "string" }

var typeString = TypeString{}

type TypeNumber struct{}

func (TypeNumber) isType()                   {}
func (TypeNumber) String(TypeContext) string { return "number" }

var typeNumber = TypeNumber{}

type TypeBoolean struct{}

func (TypeBoolean) isType()                   {}
func (TypeBoolean) String(TypeContext) string { return "bool" }

var typeBool = TypeBoolean{}

type TypeAny struct{}

func (TypeAny) isType()                   {}
func (TypeAny) String(TypeContext) string { return "any" }

var typeAny = TypeAny{}

type TypeError struct{}

func (TypeError) isType()                   {}
func (TypeError) String(TypeContext) string { return "error" }

type TypeNull struct{}

func (TypeNull) isType()                   {}
func (TypeNull) String(TypeContext) string { return "()" }

var typeNull = TypeNull{}

type TypeVoid struct{}

func (TypeVoid) isType()                   {}
func (TypeVoid) String(TypeContext) string { return "<void>" }

var typeVoid = TypeVoid{}

type TypeValue struct{ value Value }

func (TypeValue) isType()                     {}
func (t TypeValue) String(TypeContext) string { return t.value.String() }

type TypeVar struct{ id int }

func (TypeVar) isType()                     {}
func (t TypeVar) String(TypeContext) string { return fmt.Sprintf("$%d", t.id) }

type TypeComposite struct {
	// TODO: open/closed composite: {x: number} vs {[string]: number}
	// TODO: generics?
	fields map[string]Type
}

func (TypeComposite) isType()                   {}
func (TypeComposite) String(TypeContext) string { return "<void>" }

type TypeFunction struct {
	Arguments []Type
	Result    Type
}

func (TypeFunction) isType() {}
func (t TypeFunction) String(ctx TypeContext) string {
	var sb strings.Builder
	if len(t.Arguments) == 1 {
		sb.WriteString(ctx.substitute(t.Arguments[0]).String(ctx))
	} else {
		sb.WriteString("(")
		for i, arg := range t.Arguments {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(ctx.substitute(arg).String(ctx))
		}
		sb.WriteString(")")
	}
	sb.WriteString(" -> ")
	sb.WriteString(ctx.substitute(t.Result).String(ctx))
	return sb.String()
}

type TypeUnion []Type // TODO: make sure they are mutually disjoint

func (TypeUnion) isType() {}
func (t TypeUnion) String(ctx TypeContext) string {
	var sb strings.Builder
	for i, variant := range t {
		if i > 0 {
			sb.WriteString(" | ")
		}
		sb.WriteString(ctx.substitute(variant).String(ctx))
	}
	return sb.String()
}

type TypeList struct {
	elems []Type
	rest  Type // might be nil
}

func (TypeList) isType() {}
func (t TypeList) String(ctx TypeContext) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, elem := range t.elems {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ctx.substitute(elem).String(ctx))
	}
	if t.rest != nil {
		if len(t.elems) > 0 {
			sb.WriteString(", ...")
		}
		sb.WriteString(ctx.substitute(t.rest).String(ctx))
	}
	sb.WriteString("]")
	return sb.String()
}

type typeBinding struct {
	varname string
	ty      Type
}

type TypeContext struct {
	varID         *int
	bindings      []typeBinding
	substitutions map[int]Type
}

func (ctx TypeContext) String() string {
	var sb strings.Builder
	for _, decl := range ctx.bindings {
		fmt.Fprintf(&sb, "%s : %s\n", decl.varname, decl.ty.String(ctx))
	}
	return sb.String()
}

func (ctx TypeContext) Lookup(varname string) (Type, bool) {
	for _, decl := range ctx.bindings {
		if decl.varname == varname {
			return decl.ty, true
		}
	}
	return nil, false
}

func (ctx TypeContext) Append(varname string, ty Type) TypeContext {
	fmt.Println(varname, ":", ty.String(ctx))
	return TypeContext{
		ctx.varID,
		append(ctx.bindings, typeBinding{varname, ty}),
		ctx.substitutions,
	}
}

func (ctx TypeContext) typevar() Type {
	id := *ctx.varID
	*ctx.varID++
	return TypeVar{id}
}

func (ctx TypeContext) substitute(ty Type) Type {
	switch ty := ty.(type) {
	case TypeVar:
		typeRes, ok := ctx.substitutions[ty.id]
		if !ok {
			return ty
		}
		return typeRes
	case TypeComposite:
		panic("not implemented")
	case TypeFunction:
		args := make([]Type, len(ty.Arguments))
		for i, arg := range ty.Arguments {
			args[i] = ctx.substitute(arg)
		}
		res := ctx.substitute(ty.Result)
		return TypeFunction{args, res}
	default:
		return ty
	}
}

func (ctx TypeContext) unify(
	ty1, ty2 Type,
	thing string,
	pos Pos,
) Type {
	// substitute known type vars
	ty1 = ctx.substitute(ty1)
	ty2 = ctx.substitute(ty2)

	switch {
	case typeIs[TypeVar](ty1):
		id := ty1.(TypeVar).id
		ctx.substitutions[id] = ty2
		fmt.Printf("$%d is resolved to %s\n", id, ty2.String(ctx))
		return ty2
	case typeIs[TypeVar](ty2):
		id := ty2.(TypeVar).id
		ctx.substitutions[id] = ty1
		fmt.Printf("$%d is resolved to %s\n", id, ty1.String(ctx))
		return ty1
	case typeIs[TypeFunction](ty1) && typeIs[TypeFunction](ty1):
		ty1 := ty1.(TypeFunction)
		ty2 := ty2.(TypeFunction)
		if len(ty1.Arguments) != len(ty2.Arguments) {
			panicf(
				"expected %d args of types %v but found %d of types %v",
				len(ty1.Arguments), ty1.Arguments,
				len(ty2.Arguments), ty2.Arguments,
			)
		}
		args := make([]Type, len(ty1.Arguments))
		for i := range len(ty1.Arguments) {
			args[i] = ctx.unify(ty1.Arguments[i], ty2.Arguments[i], fmt.Sprintf("%s-th arg #%d", thing, i), pos)
		}
		res := ctx.unify(ty1.Result, ty2.Result, fmt.Sprintf("%s-th result", thing), pos)
		return TypeFunction{args, res}
	case isSubtypeOf(ctx, ty1, ty2):
		return ty1
	case isSubtypeOf(ctx, ty2, ty1):
		return ty2
	default:
		typeCheckError(ctx, "unify "+thing, ty1, ty2, pos)
	}
	panic(1) // unreachable
}

func typeIs[T Type](ty Type) bool {
	_, ok := ty.(T)
	return ok
}

func valueIs[V Value](v Value) bool {
	_, ok := v.(V)
	return ok
}

func typeUnion(ctx TypeContext, a, b Type) Type {
	switch {
	case typeIs[TypeVoid](a):
		return b
	case typeIs[TypeVoid](b):
		return a
	case isSubtypeOf(ctx, a, b):
		return b
	case isSubtypeOf(ctx, b, a):
		return a
	default:
		return typeAny
	}
}

func panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

// isSubtypeOf checks that subtype can be assigned to supertype (iso - is subtype of):
//   - (same type) e.g. string iso string
//   - (any type is top) any type iso any
//   - (?transitivity) if A iso B and B iso C, then A iso C
//   - (literal subtyping) literal type iso its primitive type, e.g. "aboba" type iso string
//   - (row polymorphism) {x: number, y: number} iso {x: number, y: number, z: number}
//   - (closed enum extension) 1|2 iso 1|2|3 and 1|2|... and 1|2|3|...
//   - (open enum extension) 1|2|3|... iso 1|2|...
//   - (?sum type extension) A+B iso A+B+C
//   - (intersection) A&B iso A
//   - (func/args covariance) A iso B, then B -> C iso A -> C
//   - (func/result covariance) A iso B, then C -> A iso C -> B
func isSubtypeOf(ctx TypeContext, subtype, supertype Type) bool {
	switch typeSuper := supertype.(type) {
	// any type is top
	case TypeAny:
		return true
	// void type is bottom, even void is not iso itself // TODO: ?
	case TypeVoid:
		return false
	// same types
	case TypeNull:
		return typeIs[TypeNull](subtype)
	case TypeBoolean:
		return typeIs[TypeBoolean](subtype) ||
			typeIs[TypeValue](subtype) && valueIs[ValueBoolean](subtype.(TypeValue).value)
	case TypeNumber:
		return typeIs[TypeNumber](subtype) ||
			typeIs[TypeValue](subtype) && valueIs[ValueNumber](subtype.(TypeValue).value)
	case TypeString:
		return typeIs[TypeString](subtype) ||
			typeIs[TypeValue](subtype) && valueIs[ValueString](subtype.(TypeValue).value)
	// value type is supertype only of its value type
	case TypeValue:
		return typeIs[TypeValue](subtype) &&
			// TODO: assert same value type explicitly
			fmt.Sprintf("%T", typeSuper.value) == fmt.Sprintf("%T", subtype.(TypeValue).value)
	// row polymorphism
	case TypeComposite:
		typeSub, ok := subtype.(TypeComposite)
		if !ok {
			return false
		}
		for k, typeValueSuper := range typeSuper.fields {
			typeValueSub, ok := typeSub.fields[k]
			if !ok || !isSubtypeOf(ctx, typeValueSub, typeValueSuper) {
				return false
			}
		}
		return true
	case TypeUnion:
		for _, variant := range typeSuper {
			if isSubtypeOf(ctx, subtype, variant) {
				return true
			}
		}
		return false
	// not implemented cases
	default:
		fmt.Printf("unknown subtype case: %s : %s\n", subtype.String(ctx), supertype.String(ctx))
		return false
	}
}

func typeCheckError(
	ctx TypeContext,
	thing string,
	typeExpected, typeActual Type,
	pos Pos,
) {
	panicf(
		"%s: %s is expected to be of type %s but it is %s",
		pos.String(), thing, typeExpected.String(ctx), typeActual.String(ctx),
	)
}

func typeInfer(ast *AST, n NodeID, ctx TypeContext) (Type, TypeContext) {
	switch n := ast.Nodes[n].(type) {
	case NodeLiteralBoolean:
		return TypeValue{ValueBoolean(n.Val)}, ctx
	case NodeLiteralNumber:
		return TypeValue{ValueNumber(n.Val)}, ctx
	case NodeLiteralString:
		b := []byte(n.Val)
		return TypeValue{ValueString{&b}}, ctx
	case NodeIdentifier:
		varType, ok := ctx.Lookup(n.Val)
		if !ok {
			panicf("var %s is not defined", n.Val)
		}
		return varType, ctx
	case NodeExprUnary:
		// NOTE: there is single unary operator ~ so type of unary expression is same as operand's
		return typeInfer(ast, n.Operand, ctx)
	case NodeExprBinary:
		switch op := n.Operator; op {
		case OpDefine:
			switch lhs := ast.Nodes[n.Left].(type) {
			case NodeIdentifier:
				varType := ctx.typevar()
				ctx2 := ctx.Append(lhs.Val, varType)
				rhs, _ := typeInfer(ast, n.Right, ctx2)
				return ctx2.unify(varType, rhs, "var "+lhs.Val, lhs.Pos), ctx2
			default:
				panicf("cant typecheck define operator with lhs %T", lhs)
			}
		case OpAccessor:
			lhs, _ := typeInfer(ast, n.Left, ctx)
			switch lhs := lhs.(type) {
			case TypeComposite:
				field, _ := typeInfer(ast, n.Right, ctx)
				// TODO: open composite indexing
				fieldName := string(*field.(TypeValue).value.(ValueString).b)
				typeValue, ok := lhs.fields[fieldName]
				if !ok {
					typeCheckError(ctx, "map being accessed", TypeComposite{map[string]Type{fieldName: typeAny}}, typeNull, n.Pos)
				}
				return typeValue, ctx
			default:
				return typeAny, ctx
			}
		case OpMultiply, OpSubtract, OpDivide, OpModulus: // number only operators
			lhs, _ := typeInfer(ast, n.Left, ctx)
			rhs, _ := typeInfer(ast, n.Right, ctx)
			_ = ctx.unify(lhs, typeNumber, "lhs of "+op.String(), ast.Nodes[n.Left].Position(ast))
			_ = ctx.unify(rhs, typeNumber, "rhs of "+op.String(), ast.Nodes[n.Right].Position(ast))
			return typeNumber, ctx
		case OpAdd, OpLessThan, OpGreaterThan: // T = number | string, check T op T
			lhs, _ := typeInfer(ast, n.Left, ctx)
			rhs, _ := typeInfer(ast, n.Right, ctx)
			summandType := TypeUnion{typeString, typeNumber}
			lhsType := ctx.unify(lhs, summandType, "lhs of "+op.String(), ast.Nodes[n.Left].Position(ast))
			rhsType := ctx.unify(rhs, summandType, "rhs of "+op.String(), ast.Nodes[n.Right].Position(ast))
			return typeUnion(ctx, lhsType, rhsType), ctx
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
		elems := make([]Type, len(n.Vals))
		for i, val := range n.Vals {
			elems[i], _ = typeInfer(ast, val, ctx)
		}
		return TypeList{elems, nil}, ctx
	case NodeLiteralComposite:
		fields := make(map[string]Type, len(n.Entries))
		for _, entry := range n.Entries {
			// TODO: if cant infer key literal, union them and add as open composite part
			switch key := ast.Nodes[entry.Key].(type) {
			case NodeIdentifier:
				value, _ := typeInfer(ast, entry.Val, ctx)
				fields[key.Val] = value
			case NodeLiteralString:
				value, _ := typeInfer(ast, entry.Val, ctx)
				fields[key.Val] = value
			default:
				ty, _ := typeInfer(ast, entry.Key, ctx)
				typeCheckError(ctx, "key", typeString, ty, ast.Nodes[entry.Key].Position(ast))
			}
		}
		return TypeComposite{fields}, ctx
	case NodeLiteralFunction:
		ctxFn := ctx
		args := make([]Type, len(n.Arguments))
		for i, arg := range n.Arguments {
			args[i] = ctx.typevar()
			ctxFn = ctxFn.Append(ast.Nodes[arg].(NodeIdentifier).Val, args[i])
		}
		typeResult, _ := typeInfer(ast, n.Body, ctxFn)
		return TypeFunction{args, typeResult}, ctx
	case NodeFunctionCall:
		typeFn, _ := typeInfer(ast, n.Function, ctx)
		args := make([]Type, len(n.Arguments))
		for i, arg := range n.Arguments {
			args[i], _ = typeInfer(ast, arg, ctx)
		}

		typeExpected := TypeFunction{args, ctx.typevar()}
		typeFunction, ok := ctx.unify(typeFn, typeExpected, "fn call", n.Position(ast)).(TypeFunction)
		if !ok {
			// TODO: return type might be deducted
			typeCheckError(ctx, "thing called as function", typeExpected, typeFn, ast.Nodes[n.Function].Position(ast))
		}
		if len(typeFunction.Arguments) != len(n.Arguments) {
			panicf(
				"expected %d args of types %v but found %d of types %v",
				len(typeFunction.Arguments), typeFunction.Arguments,
				len(n.Arguments), args,
			)
		}
		for i, argExpected := range typeFunction.Arguments {
			argActual := args[i]
			pos := ast.Nodes[n.Arguments[i]].Position(ast)
			_ = ctx.unify(argExpected, argActual, fmt.Sprintf("arg #%d", i), pos)
		}

		return typeFunction.Result, ctx
	case NodeExprMatch:
		typeCond, _ := typeInfer(ast, n.Condition, ctx)
		typeResult := Type(typeVoid)
		for _, clause := range n.Clauses {
			if false { // TODO: implement pattern matching
				typeTarget, _ := typeInfer(ast, clause.Target, ctx)
				if !isSubtypeOf(ctx, typeCond, typeTarget) {
					typeCheckError(ctx, "target", typeTarget, typeCond, ast.Nodes[clause.Target].Position(ast))
				}
			}

			typeExpr, _ := typeInfer(ast, clause.Expression, ctx)
			typeResult = typeUnion(ctx, typeResult, typeExpr)
		}
		return typeResult, ctx
	default:
		panicf("cant typecheck %T", n)
	}
	panic(1) // unreachable
}
