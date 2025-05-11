package internal

import (
	"bytes"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type Opcode uint8

const (
	OpConstBoolean Opcode = iota
	OpConstNumber
	OpConstString
	OpConstEmpty
	OpConstNull
	OpConstFunction
	OpOperatorUnary
	OpOperatorBinary
	OpVar
	OpVarSet
	OpSubSet
	OpPop
	OpComposite
	OpList
	OpCall
	OpReturn
	OpMatch
	OpMatchClear
	OpJmpIfNotTrue
	OpJmpIfTrue
	OpJmp
	OpDup
	OpNop
	OpScopePush
	OpScopePop
	opCount
)

type Definition struct {
	Name   string
	Widths []byte
}

var definitions = [opCount]Definition{
	OpConstBoolean:   {"BOOL      ", []byte{8}},
	OpConstNumber:    {"NUM       ", []byte{8}},
	OpConstString:    {"STR       ", []byte{8 * 2}},
	OpConstEmpty:     {"EMPTY     ", []byte{}},
	OpConstNull:      {"NULL      ", []byte{}},
	OpConstFunction:  {"FUNC      ", []byte{1}},
	OpOperatorUnary:  {"UNARY     ", []byte{1}},
	OpOperatorBinary: {"BINARY    ", []byte{1}},
	OpVar:            {"VARGET    ", []byte{1}},
	OpVarSet:         {"VARSET    ", []byte{}},
	OpSubSet:         {"SUBSET    ", []byte{}}, // TODO: arg path len
	OpPop:            {"POP       ", []byte{}},
	OpComposite:      {"COMPOSITE ", []byte{1}},
	OpList:           {"LIST      ", []byte{1}},
	OpCall:           {"CALL      ", []byte{1}},
	OpReturn:         {"RET       ", []byte{}},
	OpMatch:          {"MATCH     ", []byte{}},
	OpMatchClear:     {"MATCHCLEAR", []byte{}},
	OpJmpIfNotTrue:   {"JNE       ", []byte{1}},
	OpJmpIfTrue:      {"JE        ", []byte{1}},
	OpJmp:            {"JMP       ", []byte{1}},
	OpDup:            {"DUP       ", []byte{}},
	OpNop:            {"NOP       ", []byte{}},
	OpScopePush:      {"SCOPE_PUSH", []byte{}},
	OpScopePop:       {"SCOPE_POP ", []byte{}},
}

type Instruction struct {
	Op   Opcode
	Args []any
	Pos  Pos
}

func (ins Instruction) String() string {
	var sb strings.Builder
	sb.WriteString(definitions[ins.Op].Name)
	sb.WriteByte(' ')
	for i, arg := range ins.Args {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%T", arg)
		sb.WriteByte('(')
		fmt.Fprintf(&sb, "%v", arg)
		sb.WriteByte(')')
	}
	return sb.String()
}

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

func Make(op Opcode, pos Pos, args ...any) Instruction {
	def := definitions[op]
	assert(len(args) == len(def.Widths), "args", args, "def.Widths", def.Widths)

	// instructionLen := uint8(1)
	// for _, w := range def.Widths {
	// 	instructionLen += w
	// }

	// instruction := make([]byte, instructionLen)
	// instruction[0] = byte(op)
	// for i, operand := range operands {
	// 	// assertb(len(operand) == int(def.Widths[i]))
	// 	instruction = append(instruction, operand...)
	// }
	return Instruction{op, args, pos}
}

type compiler struct {
	AST    *AST
	funcs  [][]Instruction
	fnid   fnID
	scopes []map[string]int
}

func (c *compiler) emit(op Opcode, pos Pos, args ...any) {
	c.funcs[c.fnid] = append(c.funcs[c.fnid], Make(op, pos, args...))
}

func (c *compiler) define(lhs, rhs NodeID) {
	_, isEmpty := c.AST.Nodes[rhs].(NodeIdentifierEmpty)
	assert(!isEmpty)

	var scopeAdd func(lhs NodeID)
	scopeAdd = func(lhs NodeID) {
		switch leftSide := c.AST.Nodes[lhs].(type) {
		case NodeIdentifier: // x = y
			c.scopeAdd(leftSide.Val)
		case NodeLiteralList: // list destructure: [a, b, c] = list // TODO: test complex cases like [[m11, m12, m13], ...] = m_3x3
			for _, ln := range leftSide.Vals {
				scopeAdd(ln)
			}
		case NodeLiteralComposite: // dict destructure: {log, format: f} = std// TODO: test complex cases like {x, y: [z, f]} := {x: 1, y: [2, 3]}
			for _, ln := range leftSide.Entries {
				scopeAdd(ln.Val)
			}
		}
	}

	var emitLhs func(lhs NodeID)
	emitLhs = func(lhs NodeID) {
		l := c.AST.Nodes[lhs]
		pos := l.Position(c.AST)
		switch leftSide := l.(type) {
		case NodeIdentifier: // x = y
			c.emit(OpConstString, pos, leftSide.Val)
		case NodeLiteralList: // list destructure: [a, b, c] = list // TODO: test complex cases like [[m11, m12, m13], ...] = m_3x3
			for _, ln := range leftSide.Vals {
				emitLhs(ln)
			}
			c.emit(OpList, pos, len(leftSide.Vals))
		case NodeLiteralComposite: // dict destructure: {log, format: f} = std// TODO: test complex cases like {x, y: [z, f]} := {x: 1, y: [2, 3]}
			for _, ln := range leftSide.Entries {
				c.emit(OpConstString, pos, c.AST.Nodes[ln.Key].(NodeIdentifier).Val)
				emitLhs(ln.Val)
			}
			c.emit(OpComposite, pos, len(leftSide.Entries))
		}
	}

	l := c.AST.Nodes[lhs]
	pos := l.Position(c.AST)
	switch leftSide := l.(type) {
	case NodeIdentifier, // x = y
		NodeLiteralList,      // list destructure: [a, b, c] = list
		NodeLiteralComposite: // dict destructure: {log, format: f} = std
		scopeAdd(lhs)
		c.compile(rhs)
		emitLhs(lhs)
		c.emit(OpVarSet, pos)
	case NodeExprBinary: // x.y = z
		assert(leftSide.Operator == OpAccessor)

		// emit value [...path] ident
		c.compile(rhs)
		pathlen := 1
		switch r := c.AST.Nodes[leftSide.Right].(type) {
		case NodeIdentifier:
			c.emit(OpConstString, pos, r.Val)
		default:
			c.compile(leftSide.Right)
		}
		for func() bool {
			n, ok := c.AST.Nodes[leftSide.Left].(NodeExprBinary)
			return ok && n.Operator == OpAccessor
		}() {
			leftSide = c.AST.Nodes[leftSide.Left].(NodeExprBinary)
			r := c.AST.Nodes[leftSide.Right]
			switch r := r.(type) {
			case NodeIdentifier:
				c.emit(OpConstString, r.Position(c.AST), r.Val)
			default:
				c.compile(leftSide.Right)
			}
			pathlen++
		}
		c.emit(OpList, pos, pathlen)
		c.compile(leftSide.Left)
		c.emit(OpSubSet, pos)
	default:
		assert(false, "type", fmt.Sprintf("%T", leftSide))
	}
}

func logScope(scope map[string]int) {
	if !_debugvm {
		return
	}

	m := make([]string, len(scope))
	for k, i := range scope {
		m[i] = k
	}
	fmt.Println(m)
}

func (c *compiler) scopePush() {
	c.scopes = append(c.scopes, map[string]int{})
}

func (c *compiler) scopePop() {
	logScope(c.scope())
	c.scopes = c.scopes[:len(c.scopes)-1]
}

func (c *compiler) scopeAdd(ident string) int {
	// TODO: do not allow shadowing
	for _, scope := range slices.Backward(c.scopes) {
		if v, ok := scope[ident]; ok {
			return v
		}
	}
	scope := c.scope()
	scope[ident] = len(scope)
	return scope[ident]
}

func (c *compiler) scope() map[string]int {
	return c.scopes[len(c.scopes)-1]
}

func (c *compiler) compile(n NodeID) {
	node := c.AST.Nodes[n]
	pos := node.Position(c.AST)
	switch n := node.(type) {
	case NodeLiteralBoolean:
		c.emit(OpConstBoolean, pos, n.Val)
	case NodeLiteralNumber:
		c.emit(OpConstNumber, pos, n.Val)
	case NodeLiteralString:
		c.emit(OpConstString, pos, n.Val)
	case NodeLiteralComposite:
		for _, e := range n.Entries {
			nodeKey := c.AST.Nodes[e.Key]
			pos := nodeKey.Position(c.AST)
			switch nk := nodeKey.(type) {
			case NodeIdentifier:
				c.emit(OpConstString, pos, nk.Val)
			case NodeLiteralNumber:
				c.emit(OpConstString, pos, fmt.Sprint(nk.Val))
			default:
				c.compile(e.Key)
			}
			c.compile(e.Val)
		}
		c.emit(OpComposite, pos, len(n.Entries))
	case NodeConstFunctionCall:
		for _, arg := range n.Arguments {
			c.compile(arg)
		}
		c.emit(OpConstFunction, pos, n.Function)
		c.emit(OpCall, pos, 1)
	case NodeExprBinary:
		switch n.Operator {
		case OpDefine:
			c.define(n.Left, n.Right)
		case OpAccessor:
			c.compile(n.Left)
			r := c.AST.Nodes[n.Right]
			pos := r.Position(c.AST)
			switch keyNode := r.(type) {
			case NodeIdentifier:
				c.emit(OpConstString, pos, keyNode.Val)
			case NodeLiteralString:
				// // TODO: check type of left, if list || string, require int
				// if n, err := strconv.Itoa(); err == nil {
				// } else {
				c.emit(OpConstString, pos, keyNode.Val)
				// }
			case NodeLiteralNumber:
				// TODO: check if this is needed
				// c.emit(OpConstString, nToS(keyNode.Val))
				c.emit(OpConstNumber, pos, keyNode.Val)
			case NodeConstFunctionCall, NodeFunctionCall, NodeExprBinary, NodeExprList:
				c.compile(n.Right)
			default:
				assert(false, "type", fmt.Sprintf("%T", c.AST.Nodes[n.Right]))
			}
			c.emit(OpOperatorBinary, pos, n.Operator)
		case OpLogicalAnd:
			c.compile(n.Left)
			c.emit(OpDup, pos)
			c.emit(OpJmpIfNotTrue, pos, 9999999)
			lastBackPatch := &c.funcs[c.fnid][len(c.funcs[c.fnid])-1].Args[0]
			c.compile(n.Right)
			c.emit(OpOperatorBinary, pos, OpLogicalAnd)
			c.emit(OpNop, pos)
			*lastBackPatch = len(c.funcs[c.fnid])
		case OpLogicalOr:
			c.compile(n.Left)
			c.emit(OpDup, pos)
			c.emit(OpJmpIfTrue, pos, 9999999)
			lastBackPatch := &c.funcs[c.fnid][len(c.funcs[c.fnid])-1].Args[0]
			c.compile(n.Right)
			c.emit(OpOperatorBinary, pos, OpLogicalOr)
			c.emit(OpNop, pos)
			*lastBackPatch = len(c.funcs[c.fnid])
		default:
			c.compile(n.Left)
			c.compile(n.Right)
			c.emit(OpOperatorBinary, pos, n.Operator)
		}
	case NodeExprList:
		if len(n.Expressions) == 0 {
			c.emit(OpConstNull, pos)
		} else {
			c.emit(OpScopePush, pos)
			c.scopePush()
			for i, expr := range n.Expressions {
				c.compile(expr)
				pos := c.AST.Nodes[expr].Position(c.AST)
				if i < len(n.Expressions)-1 {
					c.emit(OpPop, pos)
				}
			}
			c.scopePop()
			c.emit(OpScopePop, pos)
		}
	case NodeIdentifier:
		c.scopeAdd(n.Val)
		c.emit(OpVar, pos, n.Val)
	case NodeIdentifierEmpty:
		c.emit(OpConstEmpty, pos)
	case NodeLiteralFunction:
		oldFnid := c.fnid
		newFnid := fnID(len(c.funcs))
		c.fnid = newFnid
		c.funcs = append(c.funcs, []Instruction{})
		fnArgs := n.Arguments
		// TODO: check function type at application site: assertb(len(fnArgs) == len(args), "fnArgs", fnArgs, "args", args)
		c.emit(OpScopePush, pos)
		c.scopePush()
		for _, argIdentNode := range slices.Backward(fnArgs) {
			switch argIdent := c.AST.Nodes[argIdentNode].(type) { // TODO: args destructure
			case NodeIdentifier:
				c.emit(OpConstString, pos, argIdent.Val)
				c.emit(OpVarSet, pos)
				c.emit(OpPop, pos)
			case NodeIdentifierEmpty:
				c.emit(OpPop, pos)
			default:
				assert(false, "arg", fmt.Sprintf("%T", argIdent))
			}
		}
		c.compile(n.Body)
		c.emit(OpReturn, pos)
		c.scopePop()
		c.emit(OpScopePop, pos)
		c.fnid = oldFnid
		c.emit(OpConstFunction, pos, ValueFunction{newFnid, &NodeLiteralFunction{Arguments: n.Arguments}, nil})
	case NodeLiteralList:
		for _, e := range n.Vals {
			c.compile(e)
		}
		c.emit(OpList, pos, len(n.Vals))
	case NodeFunctionCall:
		for _, expr := range n.Arguments {
			c.compile(expr)
		}
		c.compile(n.Function)
		c.emit(OpCall, pos, len(n.Arguments))
	case NodeExprMatch:
		c.compile(n.Condition)
		backpatches := []*any{}
		var lastBackPatch *any
		for i, clause := range n.Clauses {
			if i > 0 {
				*lastBackPatch = len(c.funcs[c.fnid])
			}
			c.compile(clause.Target)
			c.emit(OpMatch, pos)
			c.emit(OpJmpIfNotTrue, pos, 9999999)
			lastBackPatch = &c.funcs[c.fnid][len(c.funcs[c.fnid])-1].Args[0]
			c.compile(clause.Expression)
			c.emit(OpJmp, pos, 9999999)
			backpatches = append(backpatches, &c.funcs[c.fnid][len(c.funcs[c.fnid])-1].Args[0])
		}
		*lastBackPatch = len(c.funcs[c.fnid])
		for _, paddr := range backpatches {
			*paddr = len(c.funcs[c.fnid])
		}
		c.emit(OpMatchClear, pos)
	default:
		assert(false, "type", fmt.Sprintf("%T", n))
	}
}

type fnID int

type frame struct {
	fnid  fnID
	ip    int
	scope *Scope
}

// TODO: merge w/ Engine/Context
type VM struct {
	ctx    *Context
	stack  []Value
	frames []frame
}

func (vm *VM) frame() *frame {
	return &vm.frames[len(vm.frames)-1]
}

func (vm *VM) framePush(f frame) {
	vm.frames = append(vm.frames, f)
}

func (vm *VM) framePop() frame {
	f := vm.frames[len(vm.frames)-1]
	vm.frames = vm.frames[:len(vm.frames)-1]
	return f
}

func (vm *VM) push(v Value) {
	vm.stack = append(vm.stack, v)
}

func (vm *VM) dumpStack() {
	fmt.Println("  ENV:")
	for f := vm.frame().scope; f.parent != nil; f = f.parent {
		fmt.Println("   ", f.vt)
	}

	fmt.Print("  STACK: ")
	for _, v := range vm.stack {
		if v == nil {
			fmt.Print("<NIL>VAHUI< ")
			continue
		}
		fmt.Print(v.String(), " ")
	}
	fmt.Println()

	fmt.Print("  FN STACK: ")
	for _, f := range vm.frames {
		fmt.Print(f.ip, " ")
	}
	fmt.Println()
}

func (vm *VM) pop() Value {
	assert(len(vm.stack) > 0)
	res := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	if err, ok := res.(ValueError); ok { // TODO: remove
		log.Fatal().Err(err).Msg("ERROR")
	}
	return res
}

func (vm *VM) peek() Value {
	assert(len(vm.stack) > 0)
	return vm.stack[len(vm.stack)-1]
}

func unary(op Kind, arg Value, pos Pos) Value {
	switch op {
	case OpNegation:
		if isErr(arg) {
			return arg
		}

		switch o := arg.(type) {
		case ValueNumber:
			return -o
		case ValueBoolean:
			return ValueBoolean(!o)
		default:
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot negate non-boolean and non-number value %s", o), pos}}
		}
	default:
		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("unrecognized unary operator %s", "" /*n*/), pos}}
	}
}

func binary(op Kind, lhs, rhs Value, pos Pos) Value {
	switch op {
	case OpDefine:
		assert(false)
	case OpAccessor:
		if isErr(lhs) {
			return lhs
		}

		switch left := lhs.(type) {
		case ValueComposite:
			if n, ok := rhs.(ValueNumber); ok { // TODO: fix maps of non-string keys
				b := fmt.Append(nil, float64(n))
				rhs = ValueString{&b}
			}

			key := string(*rhs.(ValueString).b)
			if _, ok := left[key]; !ok {
				return Null
			}

			v := left[key]
			if s, ok := v.(ValueString); ok { // TODO: remove kostyl copy value
				b := slices.Clone(*s.b)
				return ValueString{&b}
			}
			return v
		case ValueList:
			if s, ok := rhs.(ValueString); ok { // TODO: compile into number in the first place
				f, err := strconv.ParseFloat(string(*s.b), 64)
				if err != nil {
					return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("invalid list index: %q", string(*s.b)), pos}}
				}
				rhs = ValueNumber(f)
			}

			idx := int(rhs.(ValueNumber))
			if idx < 0 || idx >= len(*left.xs) {
				return Null
				// TODO: return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("out of bounds %d while accessing list %s at an index, found non-integer index %d", idx, left, idx), pos}}
			}

			v := (*left.xs)[idx]
			if s, ok := v.(ValueString); ok { // TODO: remove kostyl copy value
				b := slices.Clone(*s.b)
				return ValueString{&b}
			}
			return v
		case ValueString:
			idx := int(rhs.(ValueNumber))
			if idx < 0 || idx >= len(*left.b) {
				return Null
			}

			b := []byte{(*left.b)[idx]}
			return ValueString{&b}
		default:
			return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot access property %v of a non-list/string/composite value %v", rhs, left), pos}}
		}
	case OpAdd: // TODO: check string + list gives nothing // TODO: check ValueError values are shown explicitly, not ignored
		if isErr(rhs) {
			return rhs
		}

		switch left := lhs.(type) {
		case ValueNumber:
			if right, ok := rhs.(ValueNumber); ok {
				return left + right
			}
		case ValueString:
			if right, ok := rhs.(ValueString); ok {
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
			if right, ok := rhs.(ValueBoolean); ok {
				return ValueBoolean(left || right)
			}
		case ValueComposite: // dict + dict
			if right, ok := rhs.(ValueComposite); ok {
				res := make(ValueComposite, len(left)+len(right))
				maps.Copy(res, left)
				maps.Copy(res, right)
				return res
			}
		case ValueList: // list + list
			if right, ok := rhs.(ValueList); ok {
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

		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support addition", lhs, rhs), pos}}
	case OpSubtract:
		if isErr(rhs) {
			return rhs
		}

		switch left := lhs.(type) {
		case ValueNumber:
			if right, ok := rhs.(ValueNumber); ok {
				return left - right
			}
		}

		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support subtraction", lhs, rhs), pos}}
	case OpMultiply:
		if isErr(rhs) {
			return rhs
		}

		switch left := lhs.(type) {
		case ValueNumber:
			if right, ok := rhs.(ValueNumber); ok {
				return left * right
			}
		// TODO: remove, same as &
		case ValueBoolean:
			if right, ok := rhs.(ValueBoolean); ok {
				return ValueBoolean(left && right)
			}
		}

		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support multiplication", lhs, rhs), pos}}
	case OpDivide:
		if isErr(rhs) {
			return rhs
		}

		if leftNum, isNum := lhs.(ValueNumber); isNum {
			if right, ok := rhs.(ValueNumber); ok {
				if right == 0 {
					return ValueError{&Err{nil, ErrRuntime, "division by zero error", pos}}
				}

				return leftNum / right
			}
		}

		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support division", lhs, rhs), pos}}
	case OpModulus:
		if isErr(rhs) {
			return rhs
		}

		if leftNum, isNum := lhs.(ValueNumber); isNum {
			if right, ok := rhs.(ValueNumber); ok {
				if right == 0 {
					return ValueError{&Err{nil, ErrRuntime, "division by zero error in modulus", pos}}
				}

				if isInteger(right) {
					return ValueNumber(int(leftNum) % int(right))
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take modulus of non-integer value %s", right.String()), pos}}
			}
		}

		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support modulus", lhs, rhs), pos}}
	case OpLogicalAnd:
		// TODO: do not evaluate `right` here
		if isErr(rhs) {
			return rhs
		}

		switch left := lhs.(type) {
		case ValueNumber:
			if right, ok := rhs.(ValueNumber); ok {
				if isInteger(left) && isInteger(right) {
					return ValueNumber(int64(left) & int64(right))
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take logical & of non-integer values %s, %s", right.String(), left.String()), pos}}
			}
		case ValueString:
			if right, ok := rhs.(ValueString); ok {
				max := max(len(*left.b), len(*right.b))

				a, b := zeroExtend(*left.b, max), zeroExtend(*right.b, max)
				c := make([]byte, max)
				for i := range c {
					c[i] = a[i] & b[i]
				}
				return ValueString{&c}
			}
		case ValueBoolean:
			right, ok := rhs.(ValueBoolean)
			if !ok {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise & of %[1]T(%[1]v) and %[2]T(%[2]v)", lhs, rhs), pos}}
			}

			if !left { // false & x = false
				return ValueBoolean(false)
			}

			if isErr(rhs) {
				return rhs
			}

			return ValueBoolean(right)
		}

		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical &", lhs, rhs), pos}}
	case OpLogicalOr:
		// TODO: do not evaluate `right` here
		if isErr(rhs) {
			return rhs
		}

		switch left := lhs.(type) {
		case ValueNumber:
			if right, ok := rhs.(ValueNumber); ok {
				if !isInteger(left) {
					return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise | of non-integer values %s, %s", right.String(), left.String()), pos}}
				}

				return ValueNumber(int64(left) | int64(right))
			}
		case ValueString:
			if right, ok := rhs.(ValueString); ok {
				max := max(len(*left.b), len(*right.b))

				a, b := zeroExtend(*left.b, max), zeroExtend(*right.b, max)
				c := make([]byte, max)
				for i := range c {
					c[i] = a[i] | b[i]
				}
				return ValueString{&c}
			}
		case ValueBoolean:
			if isErr(rhs) {
				return rhs
			}

			if left { // true | x = true
				return ValueBoolean(true)
			}

			right, ok := rhs.(ValueBoolean)
			if !ok {
				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take bitwise | of %T and %T", left, right), pos}}
			}

			return ValueBoolean(right)
		}

		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical |", lhs, rhs), pos}}
	case OpLogicalXor:
		if isErr(rhs) {
			return rhs
		}

		switch left := lhs.(type) {
		case ValueNumber:
			if right, ok := rhs.(ValueNumber); ok {
				if isInteger(left) && isInteger(right) {
					return ValueNumber(int64(left) ^ int64(right))
				}

				return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("cannot take logical ^ of non-integer values %s, %s", right.String(), left.String()), pos}}
			}
		case ValueString:
			if right, ok := rhs.(ValueString); ok {
				max := max(len(*left.b), len(*right.b))

				a, b := zeroExtend(*left.b, max), zeroExtend(*right.b, max)
				c := make([]byte, max)
				for i := range c {
					c[i] = a[i] ^ b[i]
				}
				return ValueString{&c}
			}
		case ValueBoolean:
			if right, ok := rhs.(ValueBoolean); ok {
				return ValueBoolean(left != right)
			}
		}

		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("values %s and %s do not support bitwise or logical ^", lhs, rhs), pos}}
	case OpGreaterThan:
		if isErr(rhs) {
			return rhs
		}

		switch left := lhs.(type) {
		case ValueNumber:
			if right, ok := rhs.(ValueNumber); ok {
				return ValueBoolean(left > right)
			}
		case ValueString:
			if right, ok := rhs.(ValueString); ok {
				return ValueBoolean(bytes.Compare(*left.b, *right.b) > 0)
			}
		}

		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf(">: values %s and %s do not support comparison", lhs, rhs), pos}}
	case OpLessThan:
		if isErr(rhs) {
			return rhs
		}

		switch left := lhs.(type) {
		case ValueNumber:
			if right, ok := rhs.(ValueNumber); ok {
				return ValueBoolean(left < right)
			}
		case ValueString:
			if right, ok := rhs.(ValueString); ok {
				return ValueBoolean(bytes.Compare(*left.b, *right.b) < 0)
			}
		}

		return ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("<: values %s and %s do not support comparison", lhs, rhs), pos}}
	case OpEqual:
		if isErr(rhs) {
			return rhs
		}

		return ValueBoolean(lhs.Equals(rhs))
	}
	return ValueError{&Err{nil, ErrAssert, fmt.Sprintf("unknown binary operator %s", op.String()), pos}}
}

const _debugvm = false

// const _debugvm = true

func (vm *VM) done() bool {
	return vm.frame().ip == len(vm.ctx.Engine.Cmplr.funcs[vm.frame().fnid])
}

func (vm *VM) step() {
	f := vm.frame()
	ins := vm.ctx.Engine.Cmplr.funcs[f.fnid][f.ip]
	if _debugvm {
		vm.dumpStack()
		fmt.Println(ins.String())
	}
	switch ins.Op {
	case OpConstNumber:
		num := ins.Args[0].(float64)
		vm.push(ValueNumber(num))
	case OpConstString:
		str := ins.Args[0].(string)
		vm.push(vs(str))
	case OpConstBoolean:
		b := ins.Args[0].(bool)
		vm.push(ValueBoolean(b))
	case OpConstEmpty:
		vm.push(ValueEmpty{})
	case OpConstNull:
		vm.push(Null)
	case OpOperatorUnary:
		op := ins.Args[0].(Kind)
		arg := vm.pop()
		vm.push(unary(op, arg, ins.Pos))
	case OpOperatorBinary:
		op := ins.Args[0].(Kind)
		rhs := vm.pop()
		lhs := vm.pop()
		vm.push(binary(op, lhs, rhs, ins.Pos))
	case OpList:
		length := ins.Args[0].(int)
		res := make([]Value, length)
		for i := range length {
			res[length-i-1] = vm.pop()
		}
		vm.push(ValueList{&res})
	case OpVar:
		ident := ins.Args[0].(string)
		v, ok := f.scope.Get(ident)
		assert(ok, "ident", ident)
		vm.push(v)
	case OpVarSet:
		dest := vm.pop()
		val := vm.pop()

		var set func(dest, val Value)
		set = func(dest, val Value) {
			switch d := dest.(type) {
			case ValueString: // ident
				f.scope.Set(string(*d.b), val)
			case ValueList: // list destructure
				valList := val.(ValueList)
				assert(len(*d.xs) == len(*valList.xs), "len(*d.xs)", len(*d.xs), "len(*val.xs)", len(*valList.xs))
				for i, destItem := range *d.xs {
					set(destItem, (*valList.xs)[i])
				}
			case ValueComposite: // dict destructure
				val := val.(ValueComposite)
				for k, destItem := range d {
					valValue, ok := val[k]
					assert(ok, "val", val, "k", k, "dest", d)
					set(destItem, valValue)
				}
			default:
				assert(false, "dest", dest, "val", val)
			}
		}
		set(dest, val)

		vm.push(val)
	case OpSubSet:
		lhs := vm.pop()
		path := vm.pop().(ValueList)
		val := vm.pop()
		origlhs := lhs

		for _, idx := range slices.Backward((*path.xs)[1:]) {
			switch l := lhs.(type) {
			case ValueComposite:
				lhs = l[string(*idx.(ValueString).b)]
			case ValueList:
				assert(isInteger(idx.(ValueNumber)), "idx", idx)
				lhs = (*l.xs)[int(idx.(ValueNumber))]
			default:
				assert(false, "lhs", lhs)
			}
		}

		idx := (*path.xs)[0]
		switch l := lhs.(type) {
		case ValueComposite:
			l[string(*idx.(ValueString).b)] = val
		case ValueList:
			if s, ok := idx.(ValueString); ok { // TODO: compile into number in the first place
				f, err := strconv.ParseFloat(string(*s.b), 64)
				if err != nil {
					vm.push(ValueError{&Err{nil, ErrRuntime, fmt.Sprintf("invalid list index: %q", string(*s.b)), ins.Pos}})
					return
				}
				idx = ValueNumber(f)
			}

			assert(isInteger(idx.(ValueNumber)), "idx", idx)
			idx := int(idx.(ValueNumber))
			assert(idx >= 0 && idx <= len(*l.xs), "idx", idx)
			if idx < len(*l.xs) {
				(*l.xs)[idx] = val
			} else {
				*l.xs = append(*l.xs, val)
			}
		case ValueString:
			assert(isInteger(idx.(ValueNumber)), "idx", idx)
			idx := int(idx.(ValueNumber))
			assert(idx >= 0 && idx <= len(*l.b), "idx", idx)
			val := *val.(ValueString).b
			if idx < len(*l.b) {
				s0 := string(*l.b)
				s1 := string(val)
				*l.b = []byte(s0[:min(len(s0), idx)] + s1 + s0[min(len(s0), idx+len(s1)):])
			} else {
				*l.b = append(*l.b, val...)
			}
		default:
			assert(false, "type", fmt.Sprintf("%T", l))
		}

		vm.push(origlhs)
	case OpComposite:
		length := ins.Args[0].(int)
		res := make(map[string]Value, length)
		for range length {
			v := vm.pop()
			k := vm.pop()
			res[string(*k.(ValueString).b)] = v
		}
		vm.push(ValueComposite(res))
	case OpPop:
		_ = vm.pop()
	case OpMatch:
		matcher := vm.pop()
		x := vm.peek()
		vm.push(ValueBoolean(matcher.Equals(x)))
	case OpMatchClear:
		matchResult := vm.pop()
		_ = vm.pop()
		vm.push(matchResult)
	case OpConstFunction:
		fni := ins.Args[0].(ValueFunction)
		fni.scope = vm.frame().scope
		vm.push(fni)
	case OpScopePush:
		f.scope = &Scope{f.scope, ValueTable{}}
	case OpScopePop:
		f.scope = f.scope.parent
	case OpCall: // TODO: check that NodeFunctionCall args the same len as required using function type
		fn := vm.pop()
		switch fn := fn.(type) {
		case ValueFunction:
			vm.framePush(frame{fn.id, 0, fn.scope})
			f.ip--
		case NativeFunctionValue:
			nargs := ins.Args[0].(int)
			args := make([]Value, nargs)
			for i := range nargs {
				args[nargs-i-1] = vm.pop()
			}
			vm.push(fn.exec(vm.ctx, ins.Pos, args))
		default:
			assert(false, "arg", fmt.Sprintf("%T", fn))
		}
	case OpReturn:
		vm.framePop()
		vm.frame().ip++
	case OpJmp:
		addr := ins.Args[0].(int)
		f.ip = addr - 1 // NOTE: -1 required to cancel next increment
	case OpJmpIfNotTrue:
		addr := ins.Args[0].(int)
		condition, ok := vm.pop().(ValueBoolean)
		if ok && !bool(condition) {
			f.ip = addr - 1 // NOTE: -1 required to cancel next increment
		}
	case OpJmpIfTrue:
		addr := ins.Args[0].(int)
		condition, ok := vm.pop().(ValueBoolean)
		if ok && bool(condition) {
			f.ip = addr - 1 // NOTE: -1 required to cancel next increment
		}
	case OpDup:
		val := vm.peek()
		vm.push(val)
	case OpNop:
	default:
		assert(false, "op", ins.Op)
	}
	f.ip++
}

func (vm *VM) Execute() Value {
	if _debugvm {
		// in := bufio.NewScanner(os.Stdin)
		// in.Scan()
		// line := in.Text()
		// if n, ok := func() (int, bool) {
		// 	var n int
		// 	_, err := fmt.Sscanf(line, "n %d", &n)
		// 	return n, err == nil
		// }(); ok {
		// 	for range n {
		// 		vm.step()
		// 	}
		// } else if line == "c" {
		// 	for !vm.done() {
		// 		vm.step()
		// 	}
		// } else {
		// 	vm.step()
		// }
	}

	for !vm.done() {
		f := vm.frame()
		ins := vm.ctx.Engine.Cmplr.funcs[f.fnid][f.ip]
		if len(vm.frames) == 1 && ins.Op == OpReturn { // TODO: clean // for callbacks in builtin funcs
			break
		}
		vm.step()
	}
	if _debugvm {
		vm.dumpStack()
	}
	assert(len(vm.stack) == 1, "len(vm.stack)", len(vm.stack))
	return vm.pop().(Value)
}
