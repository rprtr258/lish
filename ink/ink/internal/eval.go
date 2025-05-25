package internal

import (
	"bytes"
	"fmt"
	"maps"
	"strconv"
	"strings"

	"github.com/rprtr258/fun"
	"github.com/rs/zerolog/log"
)

const _debugvm = false

// const _debugvm = true

const _asserts = true

func assert(b bool, kvs ...any) {
	if !_asserts {
		return
	}

	if b {
		return
	}

	e := log.Fatal()
	for i := 0; i < len(kvs); i += 2 {
		e.Any(kvs[i].(string), kvs[i+1])
	}
	e.Msg("assert failed")
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

// ValueComposite includes all objects and list values
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
	defn        *Node // NodeLiteralFunction
	parentFrame *StackFrame
}

func (v ValueFunction) String() string {
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

// ValueFunctionCallThunk is an internal representation of a lazy
// function evaluation used to implement tail call optimization.
type ValueFunctionCallThunk struct {
	vt       map[string]Value
	function ValueFunction
}

func (v ValueFunctionCallThunk) String() string {
	return fmt.Sprintf("Thunk of (%s)", v.function)
}

func (v ValueFunctionCallThunk) Equals(other Value) bool {
	if _, isEmpty := other.(ValueEmpty); isEmpty {
		return true
	}

	if ov, ok := other.(ValueFunctionCallThunk); ok {
		// to compare structs containing slices, we really want
		// a pointer comparison, not a value comparison
		return &v.vt == &ov.vt && &v.function == &ov.function
	}

	return false
}

// unwrapThunk expands out a recursive structure of thunks
//
//	into a flat for loop control structure
func unwrapThunk(thunk ValueFunctionCallThunk, ast *AST) (v Value) {
	isThunk := true
	for isThunk {
		frame := &StackFrame{
			parent: thunk.function.parentFrame,
			vt:     thunk.vt,
		}
		v = ast.Nodes[thunk.function.defn.Children[0]].Eval(frame, true, ast)
		if isErr(v) {
			return v
		}
		thunk, isThunk = v.(ValueFunctionCallThunk)
	}

	return
}

// call into an Ink callback function synchronously
func evalInkFunction(ast *AST, fn Value, pos Pos, allowThunk bool, args ...Value) Value {
	// call into an Ink callback function synchronously
	switch fn := fn.(type) {
	case ValueFunction:
		// TODO: check args count matches
		argValueTable := map[string]Value{}
		for i, argNode := range fn.defn.Children[1:] {
			if i < len(args) {
				if identNode := ast.Nodes[argNode]; identNode.Kind == NodeKindIdentifier {
					argValueTable[identNode.Meta.(string)] = args[i]
				}
			}
		}

		// TCO: used for evaluating expressions that may be in tail positions
		// at the end of Nodes whose evaluation allocates another StackFrame
		// like ExpressionList and FunctionLiteral's body
		returnThunk := ValueFunctionCallThunk{
			vt:       argValueTable,
			function: fn,
		}

		if allowThunk {
			return returnThunk
		}
		return unwrapThunk(returnThunk, ast)
	case NativeFunctionValue:
		return fn.exec(fn.ctx, pos, args)
	default:
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("attempted to call a non-function value %s", fn), pos}}
	}
}

func operandToStringKey(right Node, frame *StackFrame, ast *AST) (string, *Err) {
	switch right.Kind {
	case NodeKindIdentifier:
		return right.Meta.(string), nil

	case NodeKindLiteralString:
		return right.Meta.(string), nil

	case NodeKindLiteralNumber:
		return nToS(right.Meta.(float64)), nil

	default:
		rightEvaluatedValue := right.Eval(frame, false, ast)
		if isErr(rightEvaluatedValue) {
			return "", rightEvaluatedValue.(ValueError).Err
		}

		switch rv := rightEvaluatedValue.(type) {
		case ValueString:
			return string(*rv.b), nil
		case ValueNumber:
			return nToS(float64(rv)), nil
		default:
			return "", &Err{nil, ErrRuntime, fmt.Sprintf("cannot access invalid property name %s of a composite value [%s]", rightEvaluatedValue, right.Position(ast)), Pos{}}
		}
	}
}

func define(frame *StackFrame, ast *AST, leftSide Node, rightValue Value) Value {
	if _, isEmpty := rightValue.(ValueEmpty); isEmpty {
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot assign an empty value to %s (actually anything)", leftSide), leftSide.Position(ast)}}
	}

	switch leftSide.Kind {
	case NodeKindIdentifier:
		frame.Set(leftSide.Meta.(string), rightValue)
		return rightValue
	case NodeKindExprBinary:
		if leftSide.Meta.(Kind) != OpAccessor {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot assign value to %s", leftSide), leftSide.Position(ast)}}
		}

		leftValue := ast.Nodes[leftSide.Children[0]].Eval(frame, false, ast)
		if isErr(leftValue) {
			return leftValue
		}

		leftKey, err := operandToStringKey(ast.Nodes[leftSide.Children[1]], frame, ast)
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
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing list %s at an index, found non-integer index %s", left, leftKey), ast.Nodes[leftSide.Children[1]].Position(ast)}}
			}

			if rightNum < 0 || rightNum > len(*left.xs) {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("out of bounds %d while accessing list %s at an index, found non-integer index %s", rightNum, left, leftKey), ast.Nodes[leftSide.Children[1]].Position(ast)}}
			}

			if rightNum == len(*left.xs) { // append
				*left.xs = append(*left.xs, rightValue)
			} else { // set
				(*left.xs)[rightNum] = rightValue
			}
			return left
		case ValueString:
			leftIdent := ast.Nodes[leftSide.Children[0]]
			if leftIdent.Kind != NodeKindIdentifier {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot set string %s at index because string is not an identifier", left), ast.Nodes[leftSide.Children[1]].Position(ast)}}
			}

			rightString, isString := rightValue.(ValueString)
			if !isString {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot set part of string to a non-character %s", rightValue), leftSide.Position(ast)}} // TODO: put right position
			}

			rightNum, errr := strconv.Atoi(leftKey)
			if errr != nil {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing string %s at an index, found non-integer index %s", left, leftKey), ast.Nodes[leftSide.Children[1]].Position(ast)}}
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
				frame.Update(leftIdent.Meta.(string), left)
				return left
			case rn == len(*left.b):
				*left.b = append(*left.b, *rightString.b...)
				frame.Update(leftIdent.Meta.(string), left)
				return left
			default:
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("tried to modify string %s at out of bounds index %s", left, leftKey), ast.Nodes[leftSide.Children[1]].Position(ast)}}
			}
		default:
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot set property of a non-composite value %s", leftValue), ast.Nodes[leftSide.Children[0]].Position(ast)}}
		}
	case NodeKindLiteralList: // list destructure: [a, b, c] = list
		rightList, isList := rightValue.(ValueList)
		if !isList {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure non-list value %s into list", rightValue), leftSide.Position(ast)}}
		} else if len(leftSide.Children) != len(*rightList.xs) {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure list into different length: %d value into %d", len(*rightList.xs), len(leftSide.Children)), leftSide.Position(ast)}}
		}

		xs := make([]Value, len(leftSide.Children))
		res := ValueList{&xs}
		var k_ func(int) Value
		k_ = func(i int) Value {
			if i < len(leftSide.Children) {
				leftSide := leftSide.Children[i]
				v := define(frame, ast, ast.Nodes[leftSide], (*rightList.xs)[i])
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
	case NodeKindLiteralComposite: // dict destructure: {log, format: f} = std
		rightComposite, isComposite := rightValue.(ValueComposite)
		if !isComposite {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure non-dict value %s into dict", rightValue), leftSide.Position(ast)}}
		}

		res := make(ValueComposite, len(leftSide.Children)/2)
		var k_ func(int) Value
		k_ = func(i int) Value {
			if i*2 < len(leftSide.Children) {
				keyN, val := ast.Nodes[leftSide.Children[2*i]], leftSide.Children[2*i+1]
				key, err := operandToStringKey(keyN, frame, ast)
				if err != nil {
					return ValueError{&Err{err, ErrRuntime, "invalid key in dict destructure assignment", keyN.Pos}}
				}

				rightSide, ok := rightComposite[key]
				if !ok {
					knownKeys := fun.Keys(rightComposite)
					return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure unknown key %s in dict, known keys are: %v", key, knownKeys), keyN.Position(ast)}}
				}

				v := define(frame, ast, ast.Nodes[val], rightSide)
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
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot assign value to non-identifier %s", leftSide), leftSide.Position(ast)}}
	}
}

func (n Node) Eval(frame *StackFrame, allowThunk bool, ast *AST) Value {
	switch n.Kind {
	case NodeKindIdentifier:
		// LogScope(scope)
		val, ok := frame.Get(n.Meta.(string))
		if !ok {
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("%s is not defined", n.Meta.(string)), n.Position(ast)}}
		}
		return val
	case NodeKindIdentifierEmpty:
		return ValueEmpty{}
	case NodeKindLiteralNumber:
		return ValueNumber(n.Meta.(float64))
	case NodeKindLiteralString:
		b := []byte(n.Meta.(string))
		return ValueString{&b}
	case NodeKindLiteralBoolean:
		return ValueBoolean(n.Meta.(bool))
	case NodeKindLiteralComposite:
		obj := make(ValueComposite, len(n.Children)/2)
		for i := 0; i < len(n.Children); i += 2 {
			keyStr, err := operandToStringKey(ast.Nodes[n.Children[i]], frame, ast)
			if err != nil {
				return ValueError{err}
			}

			v := ast.Nodes[n.Children[i+1]].Eval(frame, false, ast)
			if isErr(v) {
				return v
			}

			obj[keyStr] = v
		}
		return obj
	case NodeKindLiteralList:
		xs := make([]Value, len(n.Children))
		listVal := ValueList{&xs}
		for i, valn := range n.Children {
			v := ast.Nodes[valn].Eval(frame, false, ast)
			if isErr(v) {
				return v
			}
			(*listVal.xs)[i] = v
		}
		return listVal
	case NodeKindLiteralFunction:
		return ValueFunction{
			defn:        &n,
			parentFrame: frame,
		}
	case NodeKindExprMatch:
		conditionVal := ast.Nodes[n.Children[0]].Eval(frame, false, ast)
		if isErr(conditionVal) {
			return conditionVal
		}

		for i := 0; i < len(n.Children[1:]); i += 2 {
			target, expression := n.Children[1+i], n.Children[1+i+1]
			targetVal := ast.Nodes[target].Eval(frame, false, ast)
			if isErr(targetVal) {
				return targetVal
			}

			if conditionVal.Equals(targetVal) {
				return ast.Nodes[expression].Eval(frame, false, ast)
			}
		}
		return Null
	case NodeKindExprList:
		length := len(n.Children)
		if length == 0 {
			return Null
		}

		callFrame := &StackFrame{
			parent: frame,
			vt:     map[string]Value{},
		}

		for _, expr := range n.Children[:length-1] {
			if expr := ast.Nodes[expr].Eval(callFrame, false, ast); isErr(expr) {
				return expr
			}
		}
		// return values of expression lists are tail call optimized,
		// so return a maybe ThunkValue
		return ast.Nodes[n.Children[length-1]].Eval(callFrame, allowThunk, ast)
	case NodeKindExprUnary:
		switch n.Meta.(Kind) {
		case OpNegation:
			operand := ast.Nodes[n.Children[0]].Eval(frame, false, ast)
			if isErr(operand) {
				return operand
			}

			switch o := operand.(type) {
			case ValueNumber:
				return -o
			case ValueBoolean:
				return ValueBoolean(!o)
			default:
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot negate non-boolean and non-number value %s", o), ast.Nodes[n.Children[0]].Position(ast)}}
			}
		default:
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("unrecognized unary operator %s", n), n.Position(ast)}}
		}
	case NodeKindExprBinary:
		left := ast.Nodes[n.Children[0]]
		right := ast.Nodes[n.Children[1]]
		return func() Value {
			switch n.Meta.(Kind) {
			case OpDefine:
				rightValue := right.Eval(frame, false, ast)
				if err, ok := rightValue.(ValueError); ok {
					return ValueError{&Err{err.Err, ErrRuntime, "cannot evaluate right-side of assignment", ast.Nodes[n.Children[0]].Position(ast)}}
				}

				return define(frame, ast, left, rightValue)
			case OpAccessor:
				leftValue := left.Eval(frame, false, ast)
				if isErr(leftValue) {
					return leftValue
				}

				rightValueStr, err := operandToStringKey(right, frame, ast)
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
						return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing list %s at an index, found non-integer index %s", left, rightValueStr), ast.Nodes[n.Children[1]].Position(ast)}}
					}
					if rightNum < 0 || rightNum >= len(*left.xs) {
						return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("out of bounds %d while accessing list %s at an index, found non-integer index %s", rightNum, left, rightValueStr), ast.Nodes[n.Children[1]].Position(ast)}}
					}

					v := (*left.xs)[rightNum]
					return v
				case ValueString:
					rightNum, err := strconv.Atoi(rightValueStr)
					if err != nil {
						return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing string %s at an index, found non-integer index %s", left, rightValueStr), ast.Nodes[n.Children[1]].Position(ast)}}
					}

					if rn := int(rightNum); 0 <= rn && rn < len(*left.b) {
						b := []byte{(*left.b)[rn]}
						return ValueString{&b}
					}

					return Null
				default:
					return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot access property %q of a non-list/composite value %v", rightValueStr, left), ast.Nodes[n.Children[1]].Position(ast)}}
				}
			}

			leftValue := left.Eval(frame, false, ast)
			if isErr(leftValue) {
				return leftValue
			}

			switch n.Meta.(Kind) {
			case OpAdd:
				rightValue := right.Eval(frame, false, ast)
				if isErr(rightValue) {
					return rightValue
				}

				switch left := leftValue.(type) {
				case ValueNumber:
					if right, ok := rightValue.(ValueNumber); ok {
						return left + right
					}
				case ValueString:
					if right, ok := rightValue.(ValueString); ok {
						// In this context, strings are immutable. i.e. concatenating
						// strings should produce a completely new string whose modifications
						// won't be observable by the original strings.
						base := make([]byte, 0, len(*left.b)+len(*right.b))
						base = append(base, *left.b...)
						base = append(base, *right.b...)
						return ValueString{&base}
					}
				// TODO: remove, same as |
				case ValueBoolean:
					if right, ok := rightValue.(ValueBoolean); ok {
						return ValueBoolean(left || right)
					}
				case ValueComposite: // dict + dict
					if right, ok := rightValue.(ValueComposite); ok {
						res := make(ValueComposite, len(left)+len(right))
						maps.Copy(res, left)
						maps.Copy(res, right)
						return res
					}
				case ValueList: // list + list
					if right, ok := rightValue.(ValueList); ok {
						xs := make([]Value, len(*left.xs)+len(*right.xs))
						for i := range len(*left.xs) {
							xs[i] = (*left.xs)[i]
						}
						for i := range len(*right.xs) {
							xs[i+len(*left.xs)] = (*right.xs)[i]
						}
						return ValueList{&xs}
					}
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support addition", leftValue, rightValue), n.Position(ast)}}
			case OpSubtract:
				rightValue := right.Eval(frame, false, ast)
				if isErr(rightValue) {
					return rightValue
				}

				switch left := leftValue.(type) {
				case ValueNumber:
					if right, ok := rightValue.(ValueNumber); ok {
						return left - right
					}
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support subtraction", leftValue, rightValue), n.Position(ast)}}
			case OpMultiply:
				rightValue := right.Eval(frame, false, ast)
				if isErr(rightValue) {
					return rightValue
				}

				switch left := leftValue.(type) {
				case ValueNumber:
					if right, ok := rightValue.(ValueNumber); ok {
						return left * right
					}
				// TODO: remove, same as &
				case ValueBoolean:
					if right, ok := rightValue.(ValueBoolean); ok {
						return ValueBoolean(left && right)
					}
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support multiplication", leftValue, rightValue), n.Position(ast)}}
			case OpDivide:
				rightValue := right.Eval(frame, false, ast)
				if isErr(rightValue) {
					return rightValue
				}

				if leftNum, isNum := leftValue.(ValueNumber); isNum {
					if right, ok := rightValue.(ValueNumber); ok {
						if right == 0 {
							return ValueError{&Err{nil, ErrRuntime, "division by zero error", ast.Nodes[n.Children[1]].Position(ast)}}
						}

						return leftNum / right
					}
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support division", leftValue, rightValue), n.Position(ast)}}
			case OpModulus:
				rightValue := right.Eval(frame, false, ast)
				if isErr(rightValue) {
					return rightValue
				}

				if leftNum, isNum := leftValue.(ValueNumber); isNum {
					if right, ok := rightValue.(ValueNumber); ok {
						if right == 0 {
							return ValueError{&Err{nil, ErrRuntime, "division by zero error in modulus", ast.Nodes[n.Children[1]].Position(ast)}}
						}

						if isInteger(right) {
							return ValueNumber(int(leftNum) % int(right))
						}

						return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take modulus of non-integer value %s", right.String()), ast.Nodes[n.Children[0]].Position(ast)}}
					}
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support modulus", leftValue, rightValue), n.Position(ast)}}
			case OpLogicalAnd:
				// TODO: do not evaluate `right` here
				fail := func() Value {
					rightValue := right.Eval(frame, false, ast)
					if isErr(rightValue) {
						return rightValue
					}

					return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical &", leftValue, rightValue), n.Position(ast)}}
				}

				switch left := leftValue.(type) {
				case ValueNumber:
					rightValue := right.Eval(frame, false, ast)
					if isErr(rightValue) {
						return rightValue
					}

					if right, ok := rightValue.(ValueNumber); ok {
						if isInteger(left) && isInteger(right) {
							return ValueNumber(int64(left) & int64(right))
						}

						return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take logical & of non-integer values %s, %s", right.String(), left.String()), n.Position(ast)}}
					}

					return fail()
				case ValueString:
					rightValue := right.Eval(frame, false, ast)
					if isErr(rightValue) {
						return rightValue
					}

					if right, ok := rightValue.(ValueString); ok {
						max := max(len(*left.b), len(*right.b))

						a, b := zeroExtend(*left.b, max), zeroExtend(*right.b, max)
						c := make([]byte, max)
						for i := range c {
							c[i] = a[i] & b[i]
						}
						return ValueString{&c}
					}

					return fail()
				case ValueBoolean:
					if !left { // false & x = false
						return ValueBoolean(false)
					}

					rightValue := right.Eval(frame, false, ast)
					if isErr(rightValue) {
						return rightValue
					}

					right, ok := rightValue.(ValueBoolean)
					if !ok {
						return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise & of %T and %T", left, right), n.Position(ast)}}
					}

					return ValueBoolean(right)
				}

				return fail()
			case OpLogicalOr:
				// TODO: do not evaluate `right` here
				fail := func() Value {
					rightValue := right.Eval(frame, false, ast)
					if isErr(rightValue) {
						return rightValue
					}

					return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical |", leftValue, rightValue), n.Position(ast)}}
				}

				switch left := leftValue.(type) {
				case ValueNumber:
					rightValue := right.Eval(frame, false, ast)
					if isErr(rightValue) {
						return rightValue
					}

					if right, ok := rightValue.(ValueNumber); ok {
						if !isInteger(left) {
							return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise | of non-integer values %s, %s", right.String(), left.String()), n.Position(ast)}}
						}

						return ValueNumber(int64(left) | int64(right))
					}
					return fail()
				case ValueString:
					rightValue := right.Eval(frame, false, ast)
					if isErr(rightValue) {
						return rightValue
					}

					if right, ok := rightValue.(ValueString); ok {
						max := max(len(*left.b), len(*right.b))

						a, b := zeroExtend(*left.b, max), zeroExtend(*right.b, max)
						c := make([]byte, max)
						for i := range c {
							c[i] = a[i] | b[i]
						}
						return ValueString{&c}
					}

					return fail()
				case ValueBoolean:
					if left { // true | x = true
						return ValueBoolean(true)
					}

					rightValue := right.Eval(frame, false, ast)
					if isErr(rightValue) {
						return rightValue
					}

					right, ok := rightValue.(ValueBoolean)
					if !ok {
						return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise | of %T and %T", left, right), n.Position(ast)}}
					}

					return ValueBoolean(right)
				}

				return fail()
			case OpLogicalXor:
				rightValue := right.Eval(frame, false, ast)
				if isErr(rightValue) {
					return rightValue
				}

				switch left := leftValue.(type) {
				case ValueNumber:
					if right, ok := rightValue.(ValueNumber); ok {
						if isInteger(left) && isInteger(right) {
							return ValueNumber(int64(left) ^ int64(right))
						}

						return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take logical ^ of non-integer values %s, %s", right.String(), left.String()), n.Position(ast)}}
					}
				case ValueString:
					if right, ok := rightValue.(ValueString); ok {
						max := max(len(*left.b), len(*right.b))

						a, b := zeroExtend(*left.b, max), zeroExtend(*right.b, max)
						c := make([]byte, max)
						for i := range c {
							c[i] = a[i] ^ b[i]
						}
						return ValueString{&c}
					}
				case ValueBoolean:
					if right, ok := rightValue.(ValueBoolean); ok {
						return ValueBoolean(left != right)
					}
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical ^", leftValue, rightValue), n.Position(ast)}}
			case OpGreaterThan:
				rightValue := right.Eval(frame, false, ast)
				if isErr(rightValue) {
					return rightValue
				}

				switch left := leftValue.(type) {
				case ValueNumber:
					if right, ok := rightValue.(ValueNumber); ok {
						return ValueBoolean(left > right)
					}
				case ValueString:
					if right, ok := rightValue.(ValueString); ok {
						return ValueBoolean(bytes.Compare(*left.b, *right.b) > 0)
					}
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf(">: values %s and %s do not support comparison", leftValue, rightValue), n.Position(ast)}}
			case OpLessThan:
				rightValue := right.Eval(frame, false, ast)
				if isErr(rightValue) {
					return rightValue
				}

				switch left := leftValue.(type) {
				case ValueNumber:
					if right, ok := rightValue.(ValueNumber); ok {
						return ValueBoolean(left < right)
					}
				case ValueString:
					if right, ok := rightValue.(ValueString); ok {
						return ValueBoolean(bytes.Compare(*left.b, *right.b) < 0)
					}
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("<: values %s and %s do not support comparison", leftValue, rightValue), n.Position(ast)}}
			case OpEqual:
				rightValue := right.Eval(frame, false, ast)
				if isErr(rightValue) {
					return rightValue
				}

				return ValueBoolean(leftValue.Equals(rightValue))
			default:
				return ValueError{&Err{nil, ErrAssert, fmt.Sprintf("unknown binary operator %s", n.String()), Pos{}}}
			}
		}()
	case NodeKindFunctionCall:
		fn := ast.Nodes[n.Children[0]].Eval(frame, false, ast)
		if isErr(fn) {
			return fn
		}

		argResults := make([]Value, len(n.Children[1:]))
		for i, arg := range n.Children[1:] {
			argResults[i] = ast.Nodes[arg].Eval(frame, false, ast)
			if isErr(argResults[i]) {
				return argResults[i]
			}
		}
		return evalInkFunction(ast, fn, n.Position(ast), allowThunk, argResults...)
	default:
		panic(fmt.Sprint("unreachable", n.Kind))
	}
}
