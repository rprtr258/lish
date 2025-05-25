package internal

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/rprtr258/scuf"
)

type Type interface {
	isType()
	String(ctx typeContext) string
}

type TypeString struct{}

func (TypeString) isType()                   {}
func (TypeString) String(typeContext) string { return "string" }

var typeString = TypeString{}

type TypeNumber struct{}

func (TypeNumber) isType()                   {}
func (TypeNumber) String(typeContext) string { return "number" }

var typeNumber = TypeNumber{}

type TypeBoolean struct{}

func (TypeBoolean) isType()                   {}
func (TypeBoolean) String(typeContext) string { return "bool" }

var typeBool = TypeBoolean{}

type TypeAny struct{}

func (TypeAny) isType()                   {}
func (TypeAny) String(typeContext) string { return "any" }

var typeAny = TypeAny{}

type TypeError struct{}

func (TypeError) isType()                   {}
func (TypeError) String(typeContext) string { return "error" }

type TypeNull struct{}

func (TypeNull) isType()                   {}
func (TypeNull) String(typeContext) string { return "()" }

var typeNull = TypeNull{}

type TypeVoid struct{}

func (TypeVoid) isType()                   {}
func (TypeVoid) String(typeContext) string { return "<void>" }

var typeVoid = TypeVoid{}

type TypeValue struct{ value Value }

func (TypeValue) isType()                     {}
func (t TypeValue) String(typeContext) string { return t.value.String() }

type TypeVar struct{ id int }

func (TypeVar) isType()                     {}
func (t TypeVar) String(typeContext) string { return fmt.Sprintf("$%d", t.id) }

type TypeComposite struct {
	// TODO: open/closed composite: {x: number} vs {[string]: number}
	// TODO: generics?
	fields map[string]Type
}

func (TypeComposite) isType()                   {}
func (TypeComposite) String(typeContext) string { return "<void>" }

type TypeFunction struct {
	Arguments []Type
	Result    Type
}

func (TypeFunction) isType() {}
func (t TypeFunction) String(ctx typeContext) string {
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
func (t TypeUnion) String(ctx typeContext) string {
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
func (t TypeList) String(ctx typeContext) string {
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

type typeContext struct {
	varID         *int
	bindings      []typeBinding
	substitutions map[int]Type
}

func (ctx typeContext) String() string {
	var sb strings.Builder
	for _, decl := range ctx.bindings {
		fmt.Fprintf(&sb, "%s : %s\n", decl.varname, decl.ty.String(ctx))
	}
	return sb.String()
}

func (ctx typeContext) Lookup(varname string) (Type, bool) {
	for _, decl := range ctx.bindings {
		if decl.varname == varname {
			return decl.ty, true
		}
	}
	return nil, false
}

func (ctx typeContext) Append(varname string, ty Type) typeContext {
	fmt.Println(varname, ":", ty.String(ctx))
	return typeContext{
		ctx.varID,
		append(slices.Clip(ctx.bindings), typeBinding{varname, ty}),
		ctx.substitutions,
	}
}

func (ctx typeContext) typevar() Type {
	id := *ctx.varID
	*ctx.varID++
	return TypeVar{id}
}

func (ctx typeContext) substitute(ty Type) Type {
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

func (ctx typeContext) unify(
	ty1, ty2 Type,
	thing string,
	pos Pos,
) (t Type) {
	// substitute known type vars
	ty1 = ctx.substitute(ty1)
	ty2 = ctx.substitute(ty2)

	defer func() {
		if p := recover(); p != nil {
			panic(p)
		}
		scuf.
			New(os.Stdout).
			String(fmt.Sprintf("unify %s, %s |- %s", ty1.String(ctx), ty2.String(ctx), t.String(ctx)), scuf.FgBlack).
			NL()
	}()

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
	case ctx.isSubtypeOf(ty1, ty2):
		return ty1
	case ctx.isSubtypeOf(ty2, ty1):
		return ty2
	default:
		ctx.typeCheckError("unify "+thing, ty1, ty2, pos)
	}
	panic(1) // unreachable
}

func typeUnion(ctx typeContext, a, b Type) Type {
	switch {
	case typeIs[TypeVoid](a):
		return b
	case typeIs[TypeVoid](b):
		return a
	case ctx.isSubtypeOf(a, b):
		return b
	case ctx.isSubtypeOf(b, a):
		return a
	default:
		return typeAny
	}
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
//   - (???union) a | c iso a | b | c
func (ctx typeContext) isSubtypeOf(subtype, supertype Type) bool {
	subtype = ctx.substitute(subtype)
	supertype = ctx.substitute(supertype)

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
			if !ok || !ctx.isSubtypeOf(typeValueSub, typeValueSuper) {
				return false
			}
		}
		return true
	case TypeUnion:
		switch typeSub := subtype.(type) {
		case TypeUnion:
			// TODO: very stupid implementation, does not handle many cases
			for i, variant := range typeSub {
				if ctx.isSubtypeOf(variant, typeSuper[i]) {
					return true
				}
			}
			return false
		default:
			for _, variant := range typeSuper {
				if ctx.isSubtypeOf(subtype, variant) {
					return true
				}
			}
			return false
		}
	// not implemented cases
	default:
		fmt.Printf("unknown subtype case: %s : %s\n", subtype.String(ctx), supertype.String(ctx))
		return false
	}
}

func (ctx typeContext) typeCheckError(
	thing string,
	typeExpected, typeActual Type,
	pos Pos,
) {
	panicf(
		"%s: %s is expected to be of type %s but it is %s",
		pos.String(), thing, typeExpected.String(ctx), typeActual.String(ctx),
	)
}

func typeIs[T Type](ty Type) bool {
	_, ok := ty.(T)
	return ok
}

func valueIs[V Value](v Value) bool {
	_, ok := v.(V)
	return ok
}

func panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

func baseType(ty Type) Type {
	switch ty := ty.(type) {
	case TypeValue:
		switch ty.value.(type) {
		case ValueBoolean:
			return typeBool
		case ValueNumber:
			return typeNumber
		case ValueString:
			return typeString
		default:
			panicf("unknown value type %T", ty.value)
		}
	default:
		return ty
	}
	panic(1) // unreachable
}

func typeInfer(ast *AST, n NodeID, ctx typeContext) (Type, typeContext) {
	switch n := ast.Nodes[n].(type) {
	case NodeLiteralBoolean:
		return typeBool, ctx
		// return TypeValue{ValueBoolean(n.Val)}, ctx
	case NodeLiteralNumber:
		return typeNumber, ctx
		// return TypeValue{ValueNumber(n.Val)}, ctx
	case NodeLiteralString:
		return typeString, ctx
		// b := []byte(n.Val)
		// return TypeValue{ValueString{&b}}, ctx
	case NodeIdentifier:
		varType, ok := ctx.Lookup(n.Val)
		if !ok {
			panicf("var %s is not defined", n.Val)
		}
		return varType, ctx
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
				ctx.typeCheckError("key", typeString, ty, ast.Nodes[entry.Key].Position(ast))
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
	// case NodeConstFunctionCall: // TODO: get back
	// 	// NOTE: there is single unary operator ~ so type of unary expression is same as operand's
	// 	return typeInfer(ast, n.Operand, ctx)
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
			ctx.typeCheckError("thing called as function", typeExpected, typeFn, ast.Nodes[n.Function].Position(ast))
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
		typeResult := Type(nil)
		for _, clause := range n.Clauses {
			if _, ok := ast.Nodes[clause.Target].(NodeIdentifierEmpty); !ok { // TODO: implement pattern matching
				typeTarget, _ := typeInfer(ast, clause.Target, ctx)
				typeCond = ctx.unify(typeTarget, typeCond, "match target", ast.Nodes[clause.Target].Position(ast))
			}

			typeExpr, _ := typeInfer(ast, clause.Expression, ctx)
			if typeResult == nil {
				typeResult = typeExpr
			} else {
				typeResult = ctx.unify(typeResult, typeExpr, "match result", ast.Nodes[clause.Expression].Position(ast))
			}
		}
		return typeResult, ctx
	default:
		panicf("cant typecheck %T", n)
	}
	panic(1) // unreachable
}
