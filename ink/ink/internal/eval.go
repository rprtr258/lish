package internal

import (
	"bytes"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"unsafe"

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

type Cont func(Value) ValueThunk

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

// ValueThunk is an internal representation of a lazy
// function evaluation used to implement tail call optimization.
type ValueThunk func() Value

func (v ValueThunk) String() string {
	var sb strings.Builder
	// for k, v := range v.vt {
	// 	fmt.Fprintf(&sb, "%s=%s, ", k, v.String())
	// }
	return fmt.Sprintf("Thunk[%s](%v)", sb.String(), unsafe.Pointer(&v))
}

func (v ValueThunk) Equals(other Value) bool {
	// switch ov := other.(type) {
	// case ValueEmpty:
	// 	return true
	// case ValueFunctionCallThunk:
	// 	// to compare structs containing slices, we really want
	// 	// a pointer comparison, not a value comparison
	// 	return &v.vt == &ov.vt && &v.function == &ov.function
	// default:
	return false
	// }
}

func (n NodeExprUnary) Eval(scope *Scope, ast *AST, k Cont) ValueThunk {
	switch n.Operator {
	case OpNegation:
		return ast.Nodes[n.Operand].Eval(scope, ast, func(operand Value) ValueThunk {
			if isErr(operand) {
				return k(operand)
			}

			switch o := operand.(type) {
			case ValueNumber:
				return k(-o)
			case ValueBoolean:
				return k(ValueBoolean(!o))
			default:
				return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot negate non-boolean and non-number value %s", o), ast.Nodes[n.Operand].Position(ast)}})
			}
		})
	default:
		return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("unrecognized unary operator %s", n), n.Position(ast)}})
	}
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
		var rightEvaluatedValue Value
		_ = trampoline(keyOperand.Eval(scope, ast, func(v Value) ValueThunk {
			rightEvaluatedValue = v
			return func() Value {
				return Null
			}
		}))
		if err, ok := rightEvaluatedValue.(ValueError); ok {
			return "", err.Err
		}

		switch rv := rightEvaluatedValue.(type) {
		case ValueString:
			return string(rv), nil
		case ValueNumber:
			return rv.String(), nil
		default:
			return "", &Err{nil, ErrRuntime, fmt.Sprintf("cannot access invalid property name %s of a composite value", rightEvaluatedValue), keyOperand.Position(ast)}
		}
	}
}

func define(scope *Scope, ast *AST, leftNode Node, rightValue Value, k Cont) ValueThunk {
	if _, isEmpty := rightValue.(ValueEmpty); isEmpty {
		return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot assign an empty value to %s (actually anything)", leftNode), leftNode.Position(ast)}})
	}

	switch leftSide := leftNode.(type) {
	case NodeIdentifier:
		scope.Set(leftSide.Val, rightValue)
		return k(rightValue)
	case NodeExprBinary:
		if leftSide.Operator != OpAccessor {
			return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot assign value to %s", leftSide), leftNode.Position(ast)}})
		}

		return ast.Nodes[leftSide.Left].Eval(scope, ast, func(leftValue Value) ValueThunk {
			if isErr(leftValue) {
				return k(leftValue)
			}

			leftKey, err := operandToStringKey(scope, ast, ast.Nodes[leftSide.Right])
			if err != nil {
				return k(ValueError{err})
			}

			switch left := leftValue.(type) {
			case ValueComposite:
				left[leftKey] = rightValue
				return k(left)
			case ValueList:
				rightNum, errr := strconv.Atoi(leftKey)
				if errr != nil {
					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing list %s at an index, found non-integer index %s", left, leftKey), ast.Nodes[leftSide.Right].Position(ast)}})
				}

				if rightNum < 0 || rightNum > len(*left.xs) {
					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("out of bounds %d while accessing list %s at an index, found non-integer index %s", rightNum, left, leftKey), ast.Nodes[leftSide.Right].Position(ast)}})
				}

				if rightNum == len(*left.xs) { // append
					*left.xs = append(*left.xs, rightValue)
				} else { // set
					(*left.xs)[rightNum] = rightValue
				}
				return k(left)
			case ValueString:
				leftIdent, isLeftIdent := ast.Nodes[leftSide.Left].(NodeIdentifier)
				if !isLeftIdent {
					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot set string %s at index because string is not an identifier", left), ast.Nodes[leftSide.Right].Position(ast)}})
				}

				rightString, isString := rightValue.(ValueString)
				if !isString {
					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot set part of string to a non-character %s", rightValue), leftNode.Position(ast)}}) // TODO: put right position
				}

				rightNum, errr := strconv.Atoi(leftKey)
				if errr != nil {
					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing string %s at an index, found non-integer index %s", left, leftKey), ast.Nodes[leftSide.Right].Position(ast)}})
				}

				switch rn := rightNum; {
				case 0 <= rn && rn < len(left):
					for i, r := range rightString {
						if rn+i < len(left) {
							left[rn+i] = r
						} else {
							left = append(left, r)
						}
					}
					scope.Update(leftIdent.Val, left)
					return k(left)
				case rn == len(left):
					left = append(left, rightString...)
					scope.Update(leftIdent.Val, left)
					return k(left)
				default:
					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("tried to modify string %s at out of bounds index %s", left, leftKey), ast.Nodes[leftSide.Right].Position(ast)}})
				}
			default:
				return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot set property of a non-composite value %s", leftValue), ast.Nodes[leftSide.Left].Position(ast)}})
			}
		})
	case NodeLiteralList: // list destructure: [a, b, c] = list
		rightList, isList := rightValue.(ValueList)
		if !isList {
			return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure non-list value %s into list", rightValue), leftNode.Position(ast)}})
		} else if len(leftSide.Vals) != len(*rightList.xs) {
			return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure list into different length: %d value into %d", len(*rightList.xs), len(leftSide.Vals)), leftNode.Position(ast)}})
		}

		xs := make([]Value, len(leftSide.Vals))
		res := ValueList{&xs}
		var k_ func(int) ValueThunk
		k_ = func(i int) ValueThunk {
			if i < len(leftSide.Vals) {
				leftSide := leftSide.Vals[i]
				return define(scope, ast, ast.Nodes[leftSide], (*rightList.xs)[i], func(v Value) ValueThunk {
					if isErr(v) {
						return k(v)
					}
					(*res.xs)[i] = v
					return k_(i + 1)
				})
			} else {
				return k(res)
			}
		}
		return k_(0)
	case NodeLiteralComposite: // dict destructure: {log, format: f} = std
		rightComposite, isComposite := rightValue.(ValueComposite)
		if !isComposite {
			return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure non-dict value %s into dict", rightValue), leftNode.Position(ast)}})
		}

		res := make(ValueComposite, len(leftSide.Entries))
		var k_ func(int) ValueThunk
		k_ = func(i int) ValueThunk {
			if i < len(leftSide.Entries) {
				entry := leftSide.Entries[i]
				key, err := operandToStringKey(scope, ast, ast.Nodes[entry.Key])
				if err != nil {
					return k(ValueError{&Err{err, ErrRuntime, "invalid key in dict destructure assignment", entry.Pos}})
				}

				rightSide, ok := rightComposite[key]
				if !ok {
					knownKeys := fun.Keys(rightComposite)
					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot destructure unknown key %s in dict, known keys are: %v", key, knownKeys), ast.Nodes[entry.Key].Position(ast)}})
				}

				return define(scope, ast, ast.Nodes[entry.Val], rightSide, func(v Value) ValueThunk {
					if isErr(v) {
						return k(v)
					}
					res[key] = v
					return k_(i + 1)
				})
			} else {
				return k(res)
			}
		}
		return k_(0)
	default:
		// TODO: show node as-is, store position start and end instead of just start
		return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot assign value to non-identifier %s", leftNode), leftNode.Position(ast)}})
	}
}

func (n NodeExprBinary) Eval(scope *Scope, ast *AST, k Cont) ValueThunk {
	left := ast.Nodes[n.Left]
	right := ast.Nodes[n.Right]
	return func() Value {
		switch n.Operator {
		case OpDefine:
			return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
				if err, ok := rightValue.(ValueError); ok {
					return k(ValueError{&Err{err.Err, ErrRuntime, "cannot evaluate right-side of assignment", ast.Nodes[n.Left].Position(ast)}})
				}

				return define(scope, ast, left, rightValue, k)
			})
		case OpAccessor:
			return left.Eval(scope, ast, func(leftValue Value) ValueThunk {
				if isErr(leftValue) {
					return k(leftValue)
				}

				rightValueStr, err := operandToStringKey(scope, ast, right)
				if err != nil {
					return k(ValueError{err})
				}

				switch left := leftValue.(type) {
				case ValueComposite:
					if v, ok := left[rightValueStr]; ok {
						return k(v)
					}

					return k(Null)
				case ValueList:
					rightNum, err := strconv.Atoi(rightValueStr)
					if err != nil {
						return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing list %s at an index, found non-integer index %s", left, rightValueStr), ast.Nodes[n.Right].Position(ast)}})
					}
					if rightNum < 0 || rightNum >= len(*left.xs) {
						return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("out of bounds %d while accessing list %s at an index, found non-integer index %s", rightNum, left, rightValueStr), ast.Nodes[n.Right].Position(ast)}})
					}

					v := (*left.xs)[rightNum]
					return k(v)
				case ValueString:
					rightNum, err := strconv.Atoi(rightValueStr)
					if err != nil {
						return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("while accessing string %s at an index, found non-integer index %s", left, rightValueStr), ast.Nodes[n.Right].Position(ast)}})
					}

					if rn := int(rightNum); 0 <= rn && rn < len(left) {
						return k(ValueString([]byte{left[rn]}))
					}

					return k(Null)
				default:
					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot access property %q of a non-list/composite value %v", rightValueStr, left), ast.Nodes[n.Right].Position(ast)}})
				}
			})
		}

		return left.Eval(scope, ast, func(leftValue Value) ValueThunk {
			if isErr(leftValue) {
				return k(leftValue)
			}

			switch n.Operator {
			case OpAdd:
				return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
					if isErr(rightValue) {
						return k(rightValue)
					}

					switch left := leftValue.(type) {
					case ValueNumber:
						if right, ok := rightValue.(ValueNumber); ok {
							return k(left + right)
						}
					case ValueString:
						if right, ok := rightValue.(ValueString); ok {
							// In this context, strings are immutable. i.e. concatenating
							// strings should produce a completely new string whose modifications
							// won't be observable by the original strings.
							base := make([]byte, 0, len(left)+len(right))
							base = append(base, left...)
							base = append(base, right...)
							return k(ValueString(base))
						}
					// TODO: remove, same as |
					case ValueBoolean:
						if right, ok := rightValue.(ValueBoolean); ok {
							return k(ValueBoolean(left || right))
						}
					case ValueComposite: // dict + dict
						if right, ok := rightValue.(ValueComposite); ok {
							res := make(ValueComposite, len(left)+len(right))
							maps.Copy(res, left)
							maps.Copy(res, right)
							return k(res)
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
							return k(ValueList{&xs})
						}
					}

					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support addition", leftValue, rightValue), n.Position(ast)}})
				})
			case OpSubtract:
				return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
					if isErr(rightValue) {
						return k(rightValue)
					}

					switch left := leftValue.(type) {
					case ValueNumber:
						if right, ok := rightValue.(ValueNumber); ok {
							return k(left - right)
						}
					}

					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support subtraction", leftValue, rightValue), n.Position(ast)}})
				})
			case OpMultiply:
				return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
					if isErr(rightValue) {
						return k(rightValue)
					}

					switch left := leftValue.(type) {
					case ValueNumber:
						if right, ok := rightValue.(ValueNumber); ok {
							return k(left * right)
						}
					// TODO: remove, same as &
					case ValueBoolean:
						if right, ok := rightValue.(ValueBoolean); ok {
							return k(ValueBoolean(left && right))
						}
					}

					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support multiplication", leftValue, rightValue), n.Position(ast)}})
				})
			case OpDivide:
				return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
					if isErr(rightValue) {
						return k(rightValue)
					}

					if leftNum, isNum := leftValue.(ValueNumber); isNum {
						if right, ok := rightValue.(ValueNumber); ok {
							if right == 0 {
								return k(ValueError{&Err{nil, ErrRuntime, "division by zero error", ast.Nodes[n.Right].Position(ast)}})
							}

							return k(leftNum / right)
						}
					}

					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support division", leftValue, rightValue), n.Position(ast)}})
				})
			case OpModulus:
				return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
					if isErr(rightValue) {
						return k(rightValue)
					}

					if leftNum, isNum := leftValue.(ValueNumber); isNum {
						if right, ok := rightValue.(ValueNumber); ok {
							if right == 0 {
								return k(ValueError{&Err{nil, ErrRuntime, "division by zero error in modulus", ast.Nodes[n.Right].Position(ast)}})
							}

							if isInteger(right) {
								return k(ValueNumber(int(leftNum) % int(right)))
							}

							return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take modulus of non-integer value %s", right.String()), ast.Nodes[n.Left].Position(ast)}})
						}
					}

					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support modulus", leftValue, rightValue), n.Position(ast)}})
				})
			case OpLogicalAnd:
				// TODO: do not evaluate `right` here
				fail := func() ValueThunk {
					return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
						if isErr(rightValue) {
							return k(rightValue)
						}

						return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical &", leftValue, rightValue), n.Position(ast)}})
					})
				}

				switch left := leftValue.(type) {
				case ValueNumber:
					return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
						if isErr(rightValue) {
							return k(rightValue)
						}

						if right, ok := rightValue.(ValueNumber); ok {
							if isInteger(left) && isInteger(right) {
								return k(ValueNumber(int64(left) & int64(right)))
							}

							return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take logical & of non-integer values %s, %s", right.String(), left.String()), n.Position(ast)}})
						}

						return fail()
					})
				case ValueString:
					return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
						if isErr(rightValue) {
							return k(rightValue)
						}

						if right, ok := rightValue.(ValueString); ok {
							max := max(len(left), len(right))

							a, b := zeroExtend(left, max), zeroExtend(right, max)
							c := make([]byte, max)
							for i := range c {
								c[i] = a[i] & b[i]
							}
							return k(ValueString(c))
						}

						return fail()
					})
				case ValueBoolean:
					if !left { // false & x = false
						return k(ValueBoolean(false))
					}

					return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
						if isErr(rightValue) {
							return k(rightValue)
						}

						right, ok := rightValue.(ValueBoolean)
						if !ok {
							return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise & of %T and %T", left, right), n.Position(ast)}})
						}

						return k(ValueBoolean(right))
					})
				}

				return fail()
			case OpLogicalOr:
				// TODO: do not evaluate `right` here
				fail := func() ValueThunk {
					return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
						if isErr(rightValue) {
							return k(rightValue)
						}

						return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical |", leftValue, rightValue), n.Position(ast)}})
					})
				}

				switch left := leftValue.(type) {
				case ValueNumber:
					return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
						if isErr(rightValue) {
							return k(rightValue)
						}

						if right, ok := rightValue.(ValueNumber); ok {
							if !isInteger(left) {
								return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise | of non-integer values %s, %s", right.String(), left.String()), n.Position(ast)}})
							}

							return k(ValueNumber(int64(left) | int64(right)))
						}
						return fail()
					})
				case ValueString:
					return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
						if isErr(rightValue) {
							return k(rightValue)
						}

						if right, ok := rightValue.(ValueString); ok {
							max := max(len(left), len(right))

							a, b := zeroExtend(left, max), zeroExtend(right, max)
							c := make([]byte, max)
							for i := range c {
								c[i] = a[i] | b[i]
							}
							return k(ValueString(c))
						}

						return fail()
					})
				case ValueBoolean:
					if left { // true | x = true
						return k(ValueBoolean(true))
					}

					return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
						if isErr(rightValue) {
							return k(rightValue)
						}

						right, ok := rightValue.(ValueBoolean)
						if !ok {
							return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise | of %T and %T", left, right), n.Position(ast)}})
						}

						return k(ValueBoolean(right))
					})
				}

				return fail()
			case OpLogicalXor:
				return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
					if isErr(rightValue) {
						return k(rightValue)
					}

					switch left := leftValue.(type) {
					case ValueNumber:
						if right, ok := rightValue.(ValueNumber); ok {
							if isInteger(left) && isInteger(right) {
								return k(ValueNumber(int64(left) ^ int64(right)))
							}

							return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take logical ^ of non-integer values %s, %s", right.String(), left.String()), n.Position(ast)}})
						}
					case ValueString:
						if right, ok := rightValue.(ValueString); ok {
							max := max(len(left), len(right))

							a, b := zeroExtend(left, max), zeroExtend(right, max)
							c := make([]byte, max)
							for i := range c {
								c[i] = a[i] ^ b[i]
							}
							return k(ValueString(c))
						}
					case ValueBoolean:
						if right, ok := rightValue.(ValueBoolean); ok {
							return k(ValueBoolean(left != right))
						}
					}

					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical ^", leftValue, rightValue), n.Position(ast)}})
				})
			case OpGreaterThan:
				return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
					if isErr(rightValue) {
						return k(rightValue)
					}

					switch left := leftValue.(type) {
					case ValueNumber:
						if right, ok := rightValue.(ValueNumber); ok {
							return k(ValueBoolean(left > right))
						}
					case ValueString:
						if right, ok := rightValue.(ValueString); ok {
							return k(ValueBoolean(bytes.Compare(left, right) > 0))
						}
					}

					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf(">: values %s and %s do not support comparison", leftValue, rightValue), n.Position(ast)}})
				})
			case OpLessThan:
				return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
					if isErr(rightValue) {
						return k(rightValue)
					}

					switch left := leftValue.(type) {
					case ValueNumber:
						if right, ok := rightValue.(ValueNumber); ok {
							return k(ValueBoolean(left < right))
						}
					case ValueString:
						if right, ok := rightValue.(ValueString); ok {
							return k(ValueBoolean(bytes.Compare(left, right) < 0))
						}
					}

					return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("<: values %s and %s do not support comparison", leftValue, rightValue), n.Position(ast)}})
				})
			case OpEqual:
				return right.Eval(scope, ast, func(rightValue Value) ValueThunk {
					if isErr(rightValue) {
						return k(rightValue)
					}

					return k(ValueBoolean(leftValue.Equals(rightValue)))
				})
			default:
				return k(ValueError{&Err{nil, ErrAssert, fmt.Sprintf("unknown binary operator %s", n.String()), Pos{}}})
			}
		})
	}
}

func (n NodeFunctionCall) Eval(scope *Scope, ast *AST, k Cont) ValueThunk {
	return func() Value {
		return ast.Nodes[n.Function].Eval(scope, ast, func(fn Value) ValueThunk {
			if isErr(fn) {
				return k(fn)
			}

			args := make([]Value, len(n.Arguments))
			var k_ func(int) ValueThunk
			k_ = func(i int) ValueThunk {
				if i < len(n.Arguments) {
					return ast.Nodes[n.Arguments[i]].Eval(scope, ast, func(arg Value) ValueThunk {
						if isErr(args[i]) {
							return k(args[i])
						}
						args[i] = arg
						return k_(i + 1)
					})
				} else {
					return evalInkFunction(ast, fn, n.Position(ast), k, args...)
				}
			}
			return k_(0)
		})
	}
}

func evalInkFunction(ast *AST, fn Value, pos Pos, k Cont, args ...Value) ValueThunk {
	// call into an Ink callback function synchronously
	switch fn := fn.(type) {
	case ValueFunction:
		vt := make(ValueTable, len(args))
		for j, argNode := range fn.defn.Arguments {
			if j < len(args) {
				if identNode, isIdent := ast.Nodes[argNode].(NodeIdentifier); isIdent {
					vt[identNode.Val] = args[j]
				}
			}
		}

		// // TCO: used for evaluating expressions that may be in tail positions
		// // at the end of Nodes whose evaluation allocates another Scope
		// // like ExpressionList and FunctionLiteral's body
		//
		// // expand out recursive structure of thunks into flat for loop control structure
		return ast.Nodes[fn.defn.Body].Eval(&Scope{
			parent: fn.scope,
			vt:     vt,
		}, ast, k)
	case NativeFunctionValue:
		return fn.exec(fn.ctx, pos, args, k)
	default:
		return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("attempted to call a non-function value %s", fn), pos}})
	}
}

func (n NodeExprMatch) Eval(scope *Scope, ast *AST, k Cont) ValueThunk {
	return ast.Nodes[n.Condition].Eval(scope, ast, func(conditionVal Value) ValueThunk {
		if isErr(conditionVal) {
			return k(conditionVal)
		}

		var k_ func(int) ValueThunk
		k_ = func(i int) ValueThunk {
			if i < len(n.Clauses) {
				clause := n.Clauses[i]
				return ast.Nodes[clause.Target].Eval(scope, ast, func(targetVal Value) ValueThunk {
					if isErr(targetVal) {
						return k(targetVal)
					}

					if conditionVal.Equals(targetVal) {
						return ast.Nodes[clause.Expression].Eval(scope, ast, k)
					}

					return k_(i + 1)
				})
			} else {
				return k(Null)
			}
		}
		return k_(0)
	})
}

func (n NodeExprList) Eval(scope *Scope, ast *AST, k Cont) ValueThunk {
	length := len(n.Expressions)
	if length == 0 {
		return k(Null)
	}

	newScope := &Scope{
		parent: scope,
		vt:     ValueTable{},
	}
	var k_ func(int) ValueThunk
	k_ = func(i int) ValueThunk {
		if i < length-1 {
			return ast.Nodes[n.Expressions[i]].Eval(newScope, ast, func(expr Value) ValueThunk {
				if isErr(expr) {
					return k(expr)
				}
				return k_(i + 1)
			})
		} else {
			// return values of expression lists are tail call optimized,
			// so return a maybe ThunkValue
			return ast.Nodes[n.Expressions[length-1]].Eval(newScope, ast, k)
		}
	}
	return k_(0)
}

func (n NodeIdentifierEmpty) Eval(_ *Scope, _ *AST, k Cont) ValueThunk {
	return k(ValueEmpty{})
}

func (n NodeIdentifier) Eval(scope *Scope, ast *AST, k Cont) ValueThunk {
	LogScope(scope)
	return func() Value {
		val, ok := scope.Get(n.Val)
		if !ok {
			return k(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("%s is not defined", n.Val), n.Position(ast)}})
		}
		return k(val)
	}
}

func (n NodeLiteralNumber) Eval(_ *Scope, _ *AST, k Cont) ValueThunk {
	return k(ValueNumber(n.Val))
}

func (n NodeLiteralString) Eval(_ *Scope, _ *AST, k Cont) ValueThunk {
	return k(ValueString(n.Val))
}

func (n NodeLiteralBoolean) Eval(_ *Scope, _ *AST, k Cont) ValueThunk {
	return k(ValueBoolean(n.Val))
}

func (n NodeLiteralComposite) Eval(scope *Scope, ast *AST, k Cont) ValueThunk {
	obj := make(ValueComposite, len(n.Entries))
	var k_ func(int) ValueThunk
	k_ = func(i int) ValueThunk {
		if i < len(n.Entries) {
			entry := n.Entries[i]
			keyStr, err := operandToStringKey(scope, ast, ast.Nodes[entry.Key])
			if err != nil {
				return k(ValueError{err})
			}

			return ast.Nodes[entry.Val].Eval(scope, ast, func(v Value) ValueThunk {
				obj[keyStr] = v
				if isErr(v) {
					return k(v)
				}
				return k_(i + 1)
			})
		} else {
			return k(obj)
		}
	}
	return k_(0)
}

func (n NodeLiteralList) Eval(scope *Scope, ast *AST, k Cont) ValueThunk {
	xs := make([]Value, len(n.Vals))
	listVal := ValueList{&xs}
	var k_ func(int) ValueThunk
	k_ = func(i int) ValueThunk {
		if i < len(n.Vals) {
			return ast.Nodes[n.Vals[i]].Eval(scope, ast, func(v Value) ValueThunk {
				if isErr(v) {
					return k(v)
				}
				(*listVal.xs)[i] = v
				return k_(i + 1)
			})
		} else {
			return k(listVal)
		}
	}
	return k_(0)
}

func (n NodeLiteralFunction) Eval(scope *Scope, _ *AST, k Cont) ValueThunk {
	return k(ValueFunction{
		defn:  &n,
		scope: scope,
	})
}
