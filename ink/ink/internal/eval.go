package internal

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/rprtr258/fun"
	"github.com/rs/zerolog/log"
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
type ValueString []byte

var stringValueReplacer = strings.NewReplacer(
	`\`, `\\`,
	`'`, `\'`,
	"\n", `\n`,
	"\r", `\r`,
	"\t", `\t`,
)

func (v ValueString) String() string {
	return "'" + stringValueReplacer.Replace(string(v)) + "'"
}

func (v ValueString) Equals(other Value) bool {
	switch ov := other.(type) {
	case ValueEmpty:
		return true
	case ValueString:
		return bytes.Equal(v, ov)
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

// ValueComposite includes all objects and list values
type ValueComposite map[string]Value

func (v ValueComposite) isList() bool {
	for i := 0; i < len(v); i++ {
		if _, ok := v[strconv.Itoa(i)]; !ok {
			return false
		}
	}
	return true
}

func (v ValueComposite) String() string {
	var sb strings.Builder
	if v.isList() {
		n := len(v)

		sb.WriteString("[")
		for i := 0; i < n; i++ {
			val := v[strconv.Itoa(i)]
			sb.WriteString(val.String())
			if i < n-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("]")
	} else {
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
	}
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
	defn  *NodeLiteralFunction
	scope *Scope
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
	vt       ValueTable
	function ValueFunction
}

func (v ValueFunctionCallThunk) String() string {
	var sb strings.Builder
	for k, v := range v.vt {
		fmt.Fprintf(&sb, "%s=%s, ", k, v.String())
	}
	return fmt.Sprintf("Thunk[%s](%s)", sb.String(), v.function)
}

func (v ValueFunctionCallThunk) Equals(other Value) bool {
	switch ov := other.(type) {
	case ValueEmpty:
		return true
	case ValueFunctionCallThunk:
		// to compare structs containing slices, we really want
		// a pointer comparison, not a value comparison
		return &v.vt == &ov.vt && &v.function == &ov.function
	default:
		return false
	}
}

// unwrapThunk expands out a recursive structure of thunks into a flat for loop control structure
func unwrapThunk(thunk ValueFunctionCallThunk) (Value, *Err) {
	for {
		v, err := thunk.function.defn.body.Eval(&Scope{
			parent: thunk.function.scope,
			vt:     thunk.vt,
		}, true)

		var isThunk bool
		thunk, isThunk = v.(ValueFunctionCallThunk)
		if err != nil || !isThunk {
			return v, err
		}
	}
}

func (n NodeExprUnary) Eval(scope *Scope, _ bool) (Value, *Err) {
	switch n.operator {
	case OpNegation:
		operand, err := n.operand.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		switch o := operand.(type) {
		case ValueNumber:
			return -o, nil
		case ValueBoolean:
			return ValueBoolean(!o), nil
		default:
			return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot negate non-boolean and non-number value %s", o), n.operand.Position()}
		}
	default:
		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("unrecognized unary operator %s", n), n.Position()}
	}
}

func operandToStringKey(scope *Scope, keyOperand Node) (string, *Err) {
	switch keyNode := keyOperand.(type) {
	case NodeIdentifier:
		return keyNode.val, nil
	case NodeLiteralString:
		return keyNode.val, nil
	case NodeLiteralNumber:
		return nToS(keyNode.val), nil
	default:
		rightEvaluatedValue, err := keyOperand.Eval(scope, false)
		if err != nil {
			return "", err
		}

		switch rv := rightEvaluatedValue.(type) {
		case ValueString:
			return string(rv), nil
		case ValueNumber:
			return rv.String(), nil
		default:
			return "", &Err{nil, ErrRuntime, fmt.Sprintf("cannot access invalid property name %s of a composite value", rightEvaluatedValue), keyOperand.Position()}
		}
	}
}

func define(scope *Scope, leftNode Node, rightValue Value) (Value, *Err) {
	if _, isEmpty := rightValue.(ValueEmpty); isEmpty {
		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot assign an empty value to %s (actually anything)", leftNode), leftNode.Position()}
	}

	switch leftSide := leftNode.(type) {
	case NodeIdentifier:
		scope.Set(leftSide.val, rightValue)
		return rightValue, nil
	case NodeExprBinary:
		if leftSide.operator != OpAccessor {
			return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot assign value to %s", leftSide), leftNode.Position()}
		}

		leftValue, err := leftSide.left.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		leftKey, err := operandToStringKey(scope, leftSide.right)
		if err != nil {
			return nil, err
		}

		switch left := leftValue.(type) {
		case ValueComposite:
			left[leftKey] = rightValue
			return left, nil
		case ValueString:
			leftIdent, isLeftIdent := leftSide.left.(NodeIdentifier)
			if !isLeftIdent {
				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot set string %s at index because string is not an identifier", left), leftSide.right.Position()}
			}

			rightString, isString := rightValue.(ValueString)
			if !isString {
				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot set part of string to a non-character %s", rightValue), leftNode.Position()} // TODO: put right position
			}

			rightNum, errr := strconv.ParseInt(leftKey, 10, 64)
			if errr != nil {
				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("while accessing string %s at an index, found non-integer index %s", left, leftKey), leftSide.right.Position()}
			}

			switch rn := int(rightNum); {
			case 0 <= rn && rn < len(left):
				for i, r := range rightString {
					if rn+i < len(left) {
						left[rn+i] = r
					} else {
						left = append(left, r)
					}
				}
				scope.Update(leftIdent.val, left)
				return left, nil
			case rn == len(left):
				left = append(left, rightString...)
				scope.Update(leftIdent.val, left)
				return left, nil
			default:
				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("tried to modify string %s at out of bounds index %s", left, leftKey), leftSide.right.Position()}
			}
		default:
			return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot set property of a non-composite value %s", leftValue), leftSide.left.Position()}
		}
	case NodeLiteralList: // list destructure: [a, b, c] = list
		rightComposite, isComposite := rightValue.(ValueComposite)
		if !isComposite || !rightComposite.isList() {
			return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure non-list value %s into list", rightValue), leftNode.Position()}
		} else if len(leftSide.vals) != len(rightComposite) {
			return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure list into different length: %d value into %d", len(rightComposite), len(leftSide.vals)), leftNode.Position()}
		}

		res := make(ValueComposite, len(leftSide.vals))
		for i, leftSide := range leftSide.vals {
			k := strconv.Itoa(i)
			v, err := define(scope, leftSide, rightComposite[k])
			if err != nil {
				return nil, err
			}
			res[k] = v
		}
		return res, nil
	case NodeLiteralObject: // dict destructure: {log: log, format: f} = std
		rightComposite, isComposite := rightValue.(ValueComposite)
		if !isComposite {
			return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure non-list value %s into list", rightValue), leftNode.Position()}
		}

		res := make(ValueComposite, len(leftSide.entries))
		for _, entry := range leftSide.entries {
			k, err := operandToStringKey(scope, entry.key)
			if err != nil {
				return nil, &Err{err, ErrRuntime, "invalid key in dict destructure assignment", entry.Pos}
			}

			rightSide, ok := rightComposite[k]
			if !ok {
				knownKeys := fun.Keys(rightComposite)
				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure unknown key %s in dict, known keys are: %v", k, knownKeys), entry.key.Position()}
			}

			v, err := define(scope, entry.val, rightSide)
			if err != nil {
				return nil, err
			}
			res[k] = v
		}
		return res, nil
	default:
		// TODO: show node as-is, store position start and end instead of just start
		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot assign value to non-identifier %s", leftNode), leftNode.Position()}
	}
}

func (n NodeExprBinary) Eval(scope *Scope, _ bool) (Value, *Err) {
	switch n.operator {
	case OpDefine:
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, &Err{err, ErrRuntime, "cannot evaluate right-side of assignment", n.left.Position()}
		}

		return define(scope, n.left, rightValue)
	case OpAccessor:
		leftValue, err := n.left.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		rightValueStr, err := operandToStringKey(scope, n.right)
		if err != nil {
			return nil, err
		}

		switch left := leftValue.(type) {
		case ValueComposite:
			if v, ok := left[rightValueStr]; ok {
				return v, nil
			}

			return Null, nil
		case ValueString:
			rightNum, err := strconv.ParseInt(rightValueStr, 10, 64)
			if err != nil {
				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("while accessing string %s at an index, found non-integer index %s", left, rightValueStr), n.right.Position()}
			}

			if rn := int(rightNum); 0 <= rn && rn < len(left) {
				return ValueString([]byte{left[rn]}), nil
			}

			return Null, nil
		default:
			return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot access property %s of a non-composite value %s", n.right, left), n.right.Position()}
		}
	}

	leftValue, err := n.left.Eval(scope, false)
	if err != nil {
		return nil, err
	}

	switch n.operator {
	case OpAdd:
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		switch left := leftValue.(type) {
		case ValueNumber:
			if right, ok := rightValue.(ValueNumber); ok {
				return left + right, nil
			}
		case ValueString:
			if right, ok := rightValue.(ValueString); ok {
				// In this context, strings are immutable. i.e. concatenating
				// strings should produce a completely new string whose modifications
				// won't be observable by the original strings.
				base := make([]byte, 0, len(left)+len(right))
				base = append(base, left...)
				base = append(base, right...)
				return ValueString(base), nil
			}
		// TODO: remove, same as |
		case ValueBoolean:
			if right, ok := rightValue.(ValueBoolean); ok {
				return ValueBoolean(left || right), nil
			}
		case ValueComposite:
			if right, ok := rightValue.(ValueComposite); ok {
				leftIsList := left.isList()
				rightIsList := right.isList()
				if leftIsList && rightIsList { // list + list
					res := make(ValueComposite, len(left)+len(right))
					for i := 0; i < len(left); i++ {
						k := strconv.Itoa(i)
						res[k] = left[k]
					}
					for i := 0; i < len(right); i++ {
						res[strconv.Itoa(i+len(left))] = right[strconv.Itoa(i)]
					}
					return ValueComposite(res), nil
				} else if !leftIsList && !rightIsList { // dict + dict
					res := make(ValueComposite, len(left)+len(right))
					for k, v := range left {
						res[k] = v
					}
					for k, v := range right {
						res[k] = v
					}
					return ValueComposite(res), nil
				}
			}
		}

		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support addition", leftValue, rightValue), n.Position()}
	case OpSubtract:
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		switch left := leftValue.(type) {
		case ValueNumber:
			if right, ok := rightValue.(ValueNumber); ok {
				return left - right, nil
			}
		}

		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support subtraction", leftValue, rightValue), n.Position()}
	case OpMultiply:
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		switch left := leftValue.(type) {
		case ValueNumber:
			if right, ok := rightValue.(ValueNumber); ok {
				return left * right, nil
			}
		// TODO: remove, same as &
		case ValueBoolean:
			if right, ok := rightValue.(ValueBoolean); ok {
				return ValueBoolean(left && right), nil
			}
		}

		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support multiplication", leftValue, rightValue), n.Position()}
	case OpDivide:
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		if leftNum, isNum := leftValue.(ValueNumber); isNum {
			if right, ok := rightValue.(ValueNumber); ok {
				if right == 0 {
					return nil, &Err{nil, ErrRuntime, fmt.Sprintf("division by zero error"), n.right.Position()}
				}

				return leftNum / right, nil
			}
		}

		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support division", leftValue, rightValue), n.Position()}
	case OpModulus:
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		if leftNum, isNum := leftValue.(ValueNumber); isNum {
			if right, ok := rightValue.(ValueNumber); ok {
				if right == 0 {
					return nil, &Err{nil, ErrRuntime, fmt.Sprintf("division by zero error in modulus"), n.right.Position()}
				}

				if isInteger(right) {
					return ValueNumber(int(leftNum) % int(right)), nil
				}

				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot take modulus of non-integer value %s", right.String()), n.left.Position()}
			}
		}

		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support modulus", leftValue, rightValue), n.Position()}
	case OpLogicalAnd:
		switch left := leftValue.(type) {
		case ValueNumber:
			rightValue, err := n.right.Eval(scope, false)
			if err != nil {
				return nil, err
			}

			if right, ok := rightValue.(ValueNumber); ok {
				if isInteger(left) && isInteger(right) {
					return ValueNumber(int64(left) & int64(right)), nil
				}

				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot take logical & of non-integer values %s, %s", right.String(), left.String()), n.Position()}
			}
		case ValueString:
			rightValue, err := n.right.Eval(scope, false)
			if err != nil {
				return nil, err
			}

			if right, ok := rightValue.(ValueString); ok {
				max := max(len(left), len(right))

				a, b := zeroExtend(left, max), zeroExtend(right, max)
				c := make([]byte, max)
				for i := range c {
					c[i] = a[i] & b[i]
				}
				return ValueString(c), nil
			}
		case ValueBoolean:
			if !left { // false & x = false
				return ValueBoolean(false), nil
			}

			rightValue, err := n.right.Eval(scope, false)
			if err != nil {
				return nil, err
			}

			right, ok := rightValue.(ValueBoolean)
			if !ok {
				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise & of %T and %T", left, right), n.Position()}
			}

			return ValueBoolean(right), nil
		}

		// TODO: do not evaluate `right` here
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical &", leftValue, rightValue), n.Position()}
	case OpLogicalOr:
		switch left := leftValue.(type) {
		case ValueNumber:
			rightValue, err := n.right.Eval(scope, false)
			if err != nil {
				return nil, err
			}

			if right, ok := rightValue.(ValueNumber); ok {
				if !isInteger(left) || !isInteger(left) {
					return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise | of non-integer values %s, %s", right.String(), left.String()), n.Position()}
				}

				return ValueNumber(int64(left) | int64(right)), nil
			}
		case ValueString:
			rightValue, err := n.right.Eval(scope, false)
			if err != nil {
				return nil, err
			}

			if right, ok := rightValue.(ValueString); ok {
				max := max(len(left), len(right))

				a, b := zeroExtend(left, max), zeroExtend(right, max)
				c := make([]byte, max)
				for i := range c {
					c[i] = a[i] | b[i]
				}
				return ValueString(c), nil
			}
		case ValueBoolean:
			if left { // true | x = true
				return ValueBoolean(true), nil
			}

			rightValue, err := n.right.Eval(scope, false)
			if err != nil {
				return nil, err
			}

			right, ok := rightValue.(ValueBoolean)
			if !ok {
				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise | of %T and %T", left, right), n.Position()}
			}

			return ValueBoolean(right), nil
		}

		// TODO: do not evaluate `right` here
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical |", leftValue, rightValue), n.Position()}
	case OpLogicalXor:
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		switch left := leftValue.(type) {
		case ValueNumber:
			if right, ok := rightValue.(ValueNumber); ok {
				if isInteger(left) && isInteger(right) {
					return ValueNumber(int64(left) ^ int64(right)), nil
				}

				return nil, &Err{nil, ErrRuntime, fmt.Sprintf("cannot take logical ^ of non-integer values %s, %s", right.String(), left.String()), n.Position()}
			}
		case ValueString:
			if right, ok := rightValue.(ValueString); ok {
				max := max(len(left), len(right))

				a, b := zeroExtend(left, max), zeroExtend(right, max)
				c := make([]byte, max)
				for i := range c {
					c[i] = a[i] ^ b[i]
				}
				return ValueString(c), nil
			}
		case ValueBoolean:
			if right, ok := rightValue.(ValueBoolean); ok {
				return ValueBoolean(left != right), nil
			}
		}

		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical ^", leftValue, rightValue), n.Position()}
	case OpGreaterThan:
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		switch left := leftValue.(type) {
		case ValueNumber:
			if right, ok := rightValue.(ValueNumber); ok {
				return ValueBoolean(left > right), nil
			}
		case ValueString:
			if right, ok := rightValue.(ValueString); ok {
				return ValueBoolean(bytes.Compare(left, right) > 0), nil
			}
		}

		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support comparison", leftValue, rightValue), n.Position()}
	case OpLessThan:
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		switch left := leftValue.(type) {
		case ValueNumber:
			if right, ok := rightValue.(ValueNumber); ok {
				return ValueBoolean(left < right), nil
			}
		case ValueString:
			if right, ok := rightValue.(ValueString); ok {
				return ValueBoolean(bytes.Compare(left, right) < 0), nil
			}
		}

		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support comparison", leftValue, rightValue), n.Position()}
	case OpEqual:
		rightValue, err := n.right.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		return ValueBoolean(leftValue.Equals(rightValue)), nil
	default:
		log.Fatal().Stringer("kind", ErrAssert).Msgf("unknown binary operator %s", n.String())
		return nil, err
	}
}

func (n NodeFunctionCall) Eval(scope *Scope, allowThunk bool) (Value, *Err) {
	fn, err := n.function.Eval(scope, false)
	if err != nil {
		return nil, err
	}

	argResults := make([]Value, len(n.arguments))
	for i, arg := range n.arguments {
		argResults[i], err = arg.Eval(scope, false)
		if err != nil {
			return nil, err
		}
	}

	return evalInkFunction(fn, allowThunk, n.Position(), argResults...)
}

// call into an Ink callback function synchronously
func evalInkFunction(fn Value, allowThunk bool, position Pos, args ...Value) (Value, *Err) {
	switch fn := fn.(type) {
	case ValueFunction:
		argValueTable := ValueTable{}
		for i, argNode := range fn.defn.arguments {
			if i < len(args) {
				if identNode, isIdent := argNode.(NodeIdentifier); isIdent {
					argValueTable[identNode.val] = args[i]
				}
			}
		}

		// TCO: used for evaluating expressions that may be in tail positions
		// at the end of Nodes whose evaluation allocates another Scope
		// like ExpressionList and FunctionLiteral's body
		returnThunk := ValueFunctionCallThunk{
			vt:       argValueTable,
			function: fn,
		}

		if allowThunk {
			return returnThunk, nil
		}
		return unwrapThunk(returnThunk)
	case NativeFunctionValue:
		return fn.exec(fn.ctx, position, args)
	default:
		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("attempted to call a non-function value %s", fn), position}
	}
}

func (n NodeMatchClause) Eval(scope *Scope, allowThunk bool) (Value, *Err) {
	log.Fatal().Stringer("kind", ErrAssert).Msg("cannot Eval a MatchClauseNode")
	return nil, nil
}

func (n NodeMatchExpr) Eval(scope *Scope, allowThunk bool) (Value, *Err) {
	conditionVal, err := n.condition.Eval(scope, false)
	if err != nil {
		return nil, err
	}

	for _, cl := range n.clauses {
		targetVal, err := cl.target.Eval(scope, false)
		if err != nil {
			return nil, err
		}

		if conditionVal.Equals(targetVal) {
			return cl.expression.Eval(scope, allowThunk)
		}
	}

	return Null, nil
}

func (n NodeExprList) Eval(scope *Scope, allowThunk bool) (Value, *Err) {
	length := len(n.expressions)
	if length == 0 {
		return Null, nil
	}

	newScope := &Scope{
		parent: scope,
		vt:     ValueTable{},
	}
	for _, expr := range n.expressions[:length-1] {
		if _, err := expr.Eval(newScope, false); err != nil {
			return nil, err
		}
	}

	// return values of expression lists are tail call optimized,
	// so return a maybe ThunkValue
	return n.expressions[length-1].Eval(newScope, allowThunk)
}

func (n NodeIdentifierEmpty) Eval(*Scope, bool) (Value, *Err) {
	return ValueEmpty{}, nil
}

func (n NodeIdentifier) Eval(scope *Scope, _ bool) (Value, *Err) {
	val, ok := scope.Get(n.val)
	if !ok {
		return nil, &Err{nil, ErrRuntime, fmt.Sprintf("%s is not defined", n.val), n.Position()}
	}
	return val, nil
}

func (n NodeLiteralNumber) Eval(*Scope, bool) (Value, *Err) {
	return ValueNumber(n.val), nil
}

func (n NodeLiteralString) Eval(*Scope, bool) (Value, *Err) {
	return ValueString(n.val), nil
}

func (n NodeLiteralBoolean) Eval(*Scope, bool) (Value, *Err) {
	return ValueBoolean(n.val), nil
}

func (n NodeLiteralObject) Eval(scope *Scope, _ bool) (Value, *Err) {
	obj := make(ValueComposite, len(n.entries))
	for _, entry := range n.entries {
		keyStr, err := operandToStringKey(scope, entry.key)
		if err != nil {
			return nil, err
		}

		obj[keyStr], err = entry.val.Eval(scope, false)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func (n NodeObjectEntry) Eval(*Scope, bool) (Value, *Err) {
	log.Fatal().Stringer("kind", ErrAssert).Msg("cannot Eval an ObjectEntryNode")
	return nil, nil
}

func (n NodeLiteralList) Eval(scope *Scope, _ bool) (Value, *Err) {
	listVal := make(ValueComposite, len(n.vals))
	for i, n := range n.vals {
		var err *Err
		listVal[strconv.Itoa(i)], err = n.Eval(scope, false)
		if err != nil {
			return nil, err
		}
	}
	return listVal, nil
}

func (n NodeLiteralFunction) Eval(scope *Scope, _ bool) (Value, *Err) {
	return ValueFunction{
		defn:  &n,
		scope: scope,
	}, nil
}

// Engine is a single global context of Ink program execution.
//
// A single thread of execution may run within an Engine at any given moment,
// and this is ensured by an internal execution lock. An execution's Engine
// also holds all permission and debugging flags.
//
// Within an Engine, there may exist multiple Contexts that each contain different
// execution environments, running concurrently under a single lock.
type Engine struct {
	// Listeners keeps track of the concurrent threads of execution running in the Engine.
	// Call `Engine.Listeners.Wait()` to block until all concurrent execution threads finish on an Engine.
	Listeners sync.WaitGroup

	// Ink de-duplicates imported source files here, where
	// Contexts from imports are deduplicated keyed by the
	// canonicalized import path. This prevents recursive
	// imports from crashing the interpreter and allows other
	// nice functionality.
	Contexts map[string]*Context
	values   map[string]Value

	// Only a single function may write to the stack frames at any moment.
	mu sync.Mutex
}

func NewEngine() *Engine {
	return &Engine{
		Contexts:  map[string]*Context{},
		values:    map[string]Value{},
		mu:        sync.Mutex{},
		Listeners: sync.WaitGroup{},
	}
}

// CreateContext creates and initializes a new Context tied to a given Engine.
func (eng *Engine) CreateContext() *Context {
	ctx := &Context{
		Engine: eng,
		Scope: &Scope{
			parent: nil,
			vt:     ValueTable{},
		},
	}

	ctx.resetWd()
	ctx.LoadEnvironment()

	return ctx
}

// Context represents a single, isolated execution context with its global heap,
// imports, call stack, and working directory.
type Context struct {
	// WorkingDirectory is absolute path to current working dir (of module system)
	WorkingDirectory string
	// currently executing file's path, if any
	File   string
	Engine *Engine
	// Scope represents the Context's global heap
	Scope *Scope
	// TODO: store position stacke somewhere to use in error reports
}

func (ctx *Context) resetWd() {
	var err error
	ctx.WorkingDirectory, err = os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Stringer("kind", ErrSystem).Msg("could not identify current working directory")
	}
}

// Eval takes a channel of Nodes to evaluate, and executes the Ink programs defined
// in the syntax tree. Eval returns the last value of the last expression in the AST,
// or an error if there was a runtime error.
func (ctx *Context) Eval(nodes iter.Seq[Node]) (val Value, err *Err) {
	ctx.Engine.mu.Lock()
	defer ctx.Engine.mu.Unlock()

	for node := range nodes {
		if val, err = node.Eval(ctx.Scope, false); err != nil {
			LogError(err)
			break
		}
	}

	logScope(ctx.Scope)

	return
}

// ExecListener queues an asynchronous callback task to the Engine behind the Context.
// Callbacks registered this way will also run with the Engine's execution lock.
func (ctx *Context) ExecListener(callback func()) {
	ctx.Engine.Listeners.Add(1)
	go func() {
		defer ctx.Engine.Listeners.Done()

		ctx.Engine.mu.Lock()
		defer ctx.Engine.mu.Unlock()

		callback()
	}()
}

// ParseReader runs an Ink program defined by an io.Reader.
// This is the main way to invoke Ink programs from Go.
// ParseReader blocks until the Ink program exits.
func ParseReader(filename string, r io.Reader) iter.Seq[Node] {
	tokens := tokenize(filename, r)
	nodes := parse(tokens)
	return nodes
}

// ExecPath is a convenience function to Exec() a program file in a given Context.
func (ctx *Context) ExecPath(path string) (iter.Seq[Node], *Err) {
	// update Cwd for any potential import() calls this file will make
	ctx.File = path

	var r io.Reader
	if u, err := url.Parse(path); err == nil && u.Scheme != "" {
		ctx.WorkingDirectory = path
		resp, err := http.Get(path)
		if err != nil {
			return nil, &Err{nil, ErrSystem, fmt.Sprintf("could not GET %s for execution: %s", path, err.Error()), Pos{}}
		}
		defer resp.Body.Close()

		r = resp.Body
	} else {
		ctx.WorkingDirectory = filepath.Dir(path)
		file, err := os.Open(path)
		if err != nil {
			return nil, &Err{nil, ErrSystem, fmt.Sprintf("could not open %s for execution: %s", path, err.Error()), Pos{}}
		}
		defer file.Close()

		r = file
	}

	return ParseReader(path, r), nil
}
