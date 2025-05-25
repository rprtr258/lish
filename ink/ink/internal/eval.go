package internal

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/rprtr258/fun"
)

const maxPrintLen = 120

// Value represents any value in the Ink programming language.
// Each value corresponds to some primitive or object value created
// during the execution of an Ink program.
type Value interface {
	String() string
	// Equals reports whether the given value is deep-equal to the
	// receiving value. It does not compare references.
	Equals(Value) bool
}

func isErr(v Value) bool {
	_, ok := v.(ValueError)
	return ok
}

func isInteger(n ValueNumber) bool {
	// Note: this returns false for int64 outside of the float64 range,
	// but that's ok since isIntable is used to check before ops that will
	// convert values to float64's (NumberValues) anyway
	return n == ValueNumber(int64(n))
}

// Utility func to get a consistent, language spec-compliant
// string representation of numbers
func nToS(f float64) string {
	// Prefer exact integer form if possible
	if i := int64(f); f == float64(i) {
		return strconv.FormatInt(i, 10)
	}

	return strconv.FormatFloat(f, 'g', -1, 64)
}

// zero-extend a slice of bytes to given length
func zeroExtend(s []byte, max int) []byte {
	if max <= len(s) {
		return s
	}

	extended := make([]byte, max)
	copy(extended, s)
	return extended
}

type ValueError struct{ *Err }

func (v ValueError) String() string {
	return "Error(" + v.Error() + ")"
}

func (v ValueError) Equals(other Value) bool {
	e, ok := other.(ValueError)
	return ok && v.Err == e.Err
}

// ValueEmpty is the value of the empty identifier.
// it is globally unique and matches everything in equality.
type ValueEmpty struct{}

func (v ValueEmpty) String() string {
	return "_"
}

func (v ValueEmpty) Equals(other Value) bool {
	return true
}

// ValueNull is a value that only exists at the type level,
// and is represented by the empty expression list `()`.
type ValueNull struct{}

// The singleton Null value is interned into a single value
var Null = ValueNull(struct{}{})

func (ValueNull) String() string {
	return "()"
}

func (ValueNull) Equals(other Value) bool {
	if _, isEmpty := other.(ValueEmpty); isEmpty {
		return true
	}

	_, ok := other.(ValueNull)
	return ok
}

// ValueNumber represents the number type (integer and floating point)
// in the Ink language.
type ValueNumber float64

func (v ValueNumber) String() string {
	return nToS(float64(v))
}

func (v ValueNumber) Equals(other Value) bool {
	switch ov := other.(type) {
	case ValueEmpty:
		return true
	case ValueNumber:
		return v == ov
	default:
		return false
	}
}

// ValueString represents all characters and strings in Ink
type ValueString struct {
	b *[]byte
}

var stringValueReplacer = strings.NewReplacer(
	`\`, `\\`,
	`'`, `\'`,
	"\n", `\n`,
	"\r", `\r`,
	"\t", `\t`,
)

func (v ValueString) String() string {
	return "'" + stringValueReplacer.Replace(string(*v.b)) + "'"
}

func (v ValueString) Equals(other Value) bool {
	switch ov := other.(type) {
	case ValueEmpty:
		return true
	case ValueString:
		return bytes.Equal(*v.b, *ov.b)
	default:
		return false
	}
}

// ValueBoolean is either `true` or `false`
type ValueBoolean bool

func (v ValueBoolean) String() string {
	return strconv.FormatBool(bool(v))
}

func (v ValueBoolean) Equals(other Value) bool {
	switch ov := other.(type) {
	case ValueEmpty:
		return true
	case ValueBoolean:
		return v == ov
	default:
		return false
	}
}

type ValueList struct {
	xs *[]Value
}

func (v ValueList) String() string {
	n := len(*v.xs)

	var sb strings.Builder
	sb.WriteString("[")
	for i, val := range *v.xs {
		sb.WriteString(val.String())
		if i < n-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (v ValueList) Equals(other Value) bool {
	switch ov := other.(type) {
	case ValueEmpty:
		return true
	case ValueList:
		if len(*v.xs) != len(*ov.xs) {
			return false
		}

		for i, val := range *v.xs {
			otherVal := (*ov.xs)[i]
			if !val.Equals(otherVal) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// ValueComposite includes all objects
type ValueComposite map[string]Value

func (v ValueComposite) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	i := 0
	for key, val := range v {
		sb.WriteString(key)
		sb.WriteString(": ")
		sb.WriteString(val.String())
		i++
		if i < len(v) {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func (v ValueComposite) Equals(other Value) bool {
	switch ov := other.(type) {
	case ValueEmpty:
		return true
	case ValueComposite:
		if len(v) != len(ov) {
			return false
		}

		for key, val := range v {
			otherVal, ok := ov[key]
			if !ok || !val.Equals(otherVal) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// ValueFunction is the value of any variables referencing functions defined in an Ink program.
type ValueFunction struct {
	id    fnID
	defn  *NodeLiteralFunction
	scope *Scope
}

func (v ValueFunction) String() string {
	if v.scope == nil {
		return fmt.Sprintf("fn(#%d)", v.id)
	}
	// ellipsize function body at a reasonable length,
	// so as not to be too verbose in repl environments
	fstr := v.defn.String()
	if len(fstr) > maxPrintLen {
		fstr = fstr[:maxPrintLen] + ".."
	}
	return fstr
}

func (v ValueFunction) Equals(other Value) bool {
	switch ov := other.(type) {
	case ValueEmpty:
		return true
	case ValueFunction:
		// to compare structs containing slices, we really want
		// a pointer comparison, not a value comparison
		return v.defn == ov.defn
	default:
		return false
	}
}

func (n NodeFunctionCall) Eval(scope *Scope, ast *AST) Value {
	fn := ast.Nodes[n.Function].Eval(scope, ast)

	args := make([]Value, len(n.Arguments))
	for i, arg := range n.Arguments {
		args[i] = ast.Nodes[arg].Eval(scope, ast)
	}

	pos := ast.Nodes[n.Function].Position(ast)
	return evalInkFunction(nil, fn, pos, args...)
}

func (n NodeConstFunctionCall) Eval(scope *Scope, ast *AST) Value {
	args := make([]Value, len(n.Arguments))
	for i, arg := range n.Arguments {
		args[i] = ast.Nodes[arg].Eval(scope, ast)
	}

	return n.Function.exec(n.Function.ctx, n.Pos, args)
}

func operandToStringKey(scope *Scope, ast *AST, keyOperand Node) (string, *Err) {
	switch keyNode := keyOperand.(type) {
	case NodeIdentifier:
		return keyNode.Val, nil
	case NodeLiteralString:
		return keyNode.Val, nil
	case NodeLiteralNumber:
		return nToS(keyNode.Val), nil
	default:
		rightEvaluatedValue := keyOperand.Eval(scope, ast)
		if err, ok := rightEvaluatedValue.(ValueError); ok {
			return "", err.Err
		}

		switch rv := rightEvaluatedValue.(type) {
		case ValueString:
			return string(*rv.b), nil
		case ValueNumber:
			return rv.String(), nil
		default:
			return "", &Err{nil, ErrRuntime, fmt.Sprintf("cannot access invalid property name %s of a composite value", rightEvaluatedValue), keyOperand.Position(ast)}
		}
	}
}

func define(scope *Scope, ast *AST, leftNode Node, rightValue Value) Value {
	if _, isEmpty := rightValue.(ValueEmpty); isEmpty {
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot assign an empty value to %s (actually anything)", leftNode), leftNode.Position(ast)}}
	}

	switch leftSide := leftNode.(type) {
	case NodeIdentifier:
		scope.Set(leftSide.Val, rightValue)
		return rightValue
	case NodeAccessor:
		leftValue := ast.Nodes[leftSide.Arg].Eval(scope, ast)
		if isErr(leftValue) {
			return leftValue
		}

		leftKey, err := operandToStringKey(scope, ast, ast.Nodes[leftSide.Path])
		if err != nil {
			return ValueError{err}
		}

		switch left := leftValue.(type) {
		case ValueComposite:
			left[leftKey] = rightValue
			return left
		case ValueList:
			rightNum, errr := strconv.Atoi(leftKey)
			if errr != nil {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing list %s at an index, found non-integer index %s", left, leftKey), ast.Nodes[leftSide.Path].Position(ast)}}
			}

			if rightNum < 0 || rightNum > len(*left.xs) {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("out of bounds %d while accessing list %s at an index, found non-integer index %s", rightNum, left, leftKey), ast.Nodes[leftSide.Path].Position(ast)}}
			}

			if rightNum == len(*left.xs) { // append
				*left.xs = append(*left.xs, rightValue)
			} else { // set
				(*left.xs)[rightNum] = rightValue
			}
			return left
		case ValueString:
			leftIdent, isLeftIdent := ast.Nodes[leftSide.Path].(NodeIdentifier)
			if !isLeftIdent {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot set string %s at index because string is not an identifier", left), ast.Nodes[leftSide.Path].Position(ast)}}
			}

			rightString, isString := rightValue.(ValueString)
			if !isString {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot set part of string to a non-character %s", rightValue), leftNode.Position(ast)}} // TODO: put right position
			}

			rightNum, errr := strconv.Atoi(leftKey)
			if errr != nil {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing string %s at an index, found non-integer index %s", left, leftKey), ast.Nodes[leftSide.Path].Position(ast)}}
			}

			switch rn := rightNum; {
			case 0 <= rn && rn < len(*left.b):
				for i, r := range *rightString.b {
					if rn+i < len(*left.b) {
						(*left.b)[rn+i] = r
					} else {
						*left.b = append(*left.b, r)
					}
				}
				scope.Update(leftIdent.Val, left)
				return left
			case rn == len(*left.b):
				*left.b = append(*left.b, *rightString.b...)
				scope.Update(leftIdent.Val, left)
				return left
			default:
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("tried to modify string %s at out of bounds index %s", left, leftKey), ast.Nodes[leftSide.Path].Position(ast)}}
			}
		default:
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot set property of a non-composite value %s", leftValue), ast.Nodes[leftSide.Arg].Position(ast)}}
		}
	case NodeLiteralList: // list destructure: [a, b, c] = list
		rightList, isList := rightValue.(ValueList)
		if !isList {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure non-list value %s into list", rightValue), leftNode.Position(ast)}}
		} else if len(leftSide.Vals) != len(*rightList.xs) {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure list into different length: %d value into %d", len(*rightList.xs), len(leftSide.Vals)), leftNode.Position(ast)}}
		}

		xs := make([]Value, len(leftSide.Vals))
		res := ValueList{&xs}
		var k_ func(int) Value
		k_ = func(i int) Value {
			if i < len(leftSide.Vals) {
				leftSide := leftSide.Vals[i]
				v := define(scope, ast, ast.Nodes[leftSide], (*rightList.xs)[i])
				if isErr(v) {
					return v
				}
				(*res.xs)[i] = v
				return k_(i + 1)
			} else {
				return res
			}
		}
		return k_(0)
	case NodeLiteralComposite: // dict destructure: {log, format: f} = std
		rightComposite, isComposite := rightValue.(ValueComposite)
		if !isComposite {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure non-dict value %s into dict", rightValue), leftNode.Position(ast)}}
		}

		res := make(ValueComposite, len(leftSide.Entries))
		var k_ func(int) Value
		k_ = func(i int) Value {
			if i < len(leftSide.Entries) {
				entry := leftSide.Entries[i]
				key, err := operandToStringKey(scope, ast, ast.Nodes[entry.Key])
				if err != nil {
					return ValueError{&Err{err, ErrRuntime, "invalid key in dict destructure assignment", entry.Pos}}
				}

				rightSide, ok := rightComposite[key]
				if !ok {
					knownKeys := fun.Keys(rightComposite)
					return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure unknown key %s in dict, known keys are: %v", key, knownKeys), ast.Nodes[entry.Key].Position(ast)}}
				}

				v := define(scope, ast, ast.Nodes[entry.Val], rightSide)
				if isErr(v) {
					return v
				}
				res[key] = v
				return k_(i + 1)
			} else {
				return res
			}
		}
		return k_(0)
	default:
		// TODO: show node as-is, store position start and end instead of just start
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot assign value to non-identifier %s", leftNode), leftNode.Position(ast)}}
	}
}

func (n NodeDefine) Eval(scope *Scope, ast *AST) Value {
	// ty: func() Type {
	// 	switch lhs := ast.Nodes[n.Left].(type) {
	// 	case NodeIdentifier:
	// 		varType := ctx.typevar()
	// 		ctx2 := ctx.Append(lhs.Val, varType)
	// 		rhs, _ := typeInfer(ast, n.Right, ctx2)
	// 		return ctx2.unify(varType, rhs, "var "+lhs.Val, lhs.Pos), ctx2
	// 	default:
	// 		panicf("cant typecheck define operator with lhs %T", lhs)
	// 	}
	// },

	rightValue := ast.Nodes[n.Value].Eval(scope, ast)
	if err, ok := rightValue.(ValueError); ok {
		return ValueError{&Err{err.Err, ErrRuntime, "cannot evaluate right-side of assignment", ast.Nodes[n.Defined].Position(ast)}}
	}

	return define(scope, ast, ast.Nodes[n.Defined], rightValue)
}

func (n NodeAccessor) Eval(scope *Scope, ast *AST) Value {
	// ty: func() Type {
	// 	assert(len(v) == 2)
	// 	lhs, _ := typeInfer(ast, n.Left, ctx)
	// 	switch lhs := lhs.(type) {
	// 	case TypeComposite:
	// 		field, _ := typeInfer(ast, n.Right, ctx)
	// 		// TODO: open composite indexing
	// 		fieldName := string(*field.(TypeValue).value.(ValueString).b)
	// 		typeValue, ok := lhs.fields[fieldName]
	// 		if !ok {
	// 			ctx.typeCheckError("map being accessed", TypeComposite{map[string]Type{fieldName: typeAny}}, typeNull, n.Pos)
	// 		}
	// 		return typeValue, ctx
	// 	default:
	// 		return typeAny, ctx
	// 	}
	// },

	leftValue := ast.Nodes[n.Arg].Eval(scope, ast)
	if isErr(leftValue) {
		return leftValue
	}

	rightValueStr, err := operandToStringKey(scope, ast, ast.Nodes[n.Path])
	if err != nil {
		return ValueError{err}
	}

	switch left := leftValue.(type) {
	case ValueComposite:
		if v, ok := left[rightValueStr]; ok {
			return v
		}

		return Null
	case ValueList:
		rightNum, err := strconv.Atoi(rightValueStr)
		if err != nil {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing list %s at an index, found non-integer index %s", left, rightValueStr), ast.Nodes[n.Path].Position(ast)}}
		}
		if rightNum < 0 || rightNum >= len(*left.xs) {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("out of bounds %d while accessing list %s at an index, found non-integer index %s", rightNum, left, rightValueStr), ast.Nodes[n.Path].Position(ast)}}
		}

		v := (*left.xs)[rightNum]
		return v
	case ValueString:
		rightNum, err := strconv.Atoi(rightValueStr)
		if err != nil {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing string %s at an index, found non-integer index %s", left, rightValueStr), ast.Nodes[n.Path].Position(ast)}}
		}

		if rn := int(rightNum); 0 <= rn && rn < len(*left.b) {
			b := []byte{(*left.b)[rn]}
			return ValueString{&b}
		}

		return Null
	default:
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot access property %q of a non-list/composite value %v", rightValueStr, left), ast.Nodes[n.Path].Position(ast)}}
	}
}

func evalInkFunction(ctx *Context, fn Value, pos Pos, args ...Value) Value {
	// call into an Ink callback function synchronously
	switch fn := fn.(type) {
	case ValueFunction:
		// TODO: // TCO: used for evaluating expressions that may be in tail positions
		// // at the end of Nodes whose evaluation allocates another Scope
		// // like ExpressionList and FunctionLiteral's body
		//
		// // expand out recursive structure of thunks into flat for loop control structure

		vm := &VM{ctx, args, []frame{{fn.id, 0, fn.scope}}}
		return vm.Execute()
	case ValueNativeFunction:
		return fn.exec(fn.ctx, pos, args)
	default:
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("attempted to call a non-function value %s", fn), pos}}
	}
}

func (n NodeExprMatch) Eval(scope *Scope, ast *AST) Value {
	conditionVal := ast.Nodes[n.Condition].Eval(scope, ast)
	if isErr(conditionVal) {
		return conditionVal
	}

	for _, clause := range n.Clauses {
		targetVal := ast.Nodes[clause.Target].Eval(scope, ast)
		if isErr(targetVal) {
			return targetVal
		}

		if conditionVal.Equals(targetVal) {
			return ast.Nodes[clause.Expression].Eval(scope, ast)
		}
	}
	return Null
}

func (n NodeExprList) Eval(scope *Scope, ast *AST) Value {
	length := len(n.Expressions)
	if length == 0 {
		return Null
	}

	newScope := &Scope{
		parent: scope,
		vt:     ValueTable{},
	}
	for i := range length - 1 {
		if expr := ast.Nodes[n.Expressions[i]].Eval(newScope, ast); isErr(expr) {
			return expr
		}
	}
	// return values of expression lists are tail call optimized,
	// so return a maybe ThunkValue
	return ast.Nodes[n.Expressions[length-1]].Eval(newScope, ast)
}

func (n NodeIdentifierEmpty) Eval(_ *Scope, _ *AST) Value {
	return ValueEmpty{}
}

func (n NodeIdentifier) Eval(scope *Scope, ast *AST) Value {
	LogScope(scope)
	val, ok := scope.Get(n.Val)
	if !ok {
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("%s is not defined", n.Val), n.Position(ast)}}
	}
	return val
}

func (n NodeLiteralNumber) Eval(_ *Scope, _ *AST) Value {
	return ValueNumber(n.Val)
}

func (n NodeLiteralString) Eval(_ *Scope, _ *AST) Value {
	b := []byte(n.Val)
	return ValueString{&b}
}

func (n NodeLiteralBoolean) Eval(_ *Scope, _ *AST) Value {
	return ValueBoolean(n.Val)
}

func (n NodeLiteralComposite) Eval(scope *Scope, ast *AST) Value {
	obj := make(ValueComposite, len(n.Entries))
	for _, entry := range n.Entries {
		keyStr, err := operandToStringKey(scope, ast, ast.Nodes[entry.Key])
		if err != nil {
			return ValueError{err}
		}

		v := ast.Nodes[entry.Val].Eval(scope, ast)
		if isErr(v) {
			return v
		}

		obj[keyStr] = v
	}
	return obj
}

func (n NodeLiteralList) Eval(scope *Scope, ast *AST) Value {
	xs := make([]Value, len(n.Vals))
	listVal := ValueList{&xs}
	for i, valn := range n.Vals {
		v := ast.Nodes[valn].Eval(scope, ast)
		if isErr(v) {
			return v
		}
		(*listVal.xs)[i] = v
	}
	return listVal
}

func (n NodeLiteralFunction) Eval(scope *Scope, _ *AST) Value {
	return ValueFunction{
		defn:  &n,
		scope: scope,
	}
}
