package internal

import (
	"fmt"
	"strings"
)

// Node represents an abstract syntax tree (AST) node in an Ink program.
type Node interface {
	String() string
	Position() position
	Eval(*Scope, bool) (Value, *Err)
}

type NodeExprUnary struct {
	position
	operator Kind
	operand  Node
}

func (n NodeExprUnary) String() string {
	return fmt.Sprintf("Unary %s (%s)", n.operator, n.operand)
}

func (n NodeExprUnary) Position() position {
	return n.position
}

type NodeExprBinary struct {
	position
	operator    Kind
	left, right Node
}

func (n NodeExprBinary) String() string {
	return fmt.Sprintf("Binary (%s) %s (%s)", n.left, n.operator, n.right)
}

func (n NodeExprBinary) Position() position {
	return n.position
}

type NodeFunctionCall struct {
	function  Node
	arguments []Node
}

func (n NodeFunctionCall) String() string {
	var sb strings.Builder
	sb.WriteString("Call (")
	sb.WriteString(n.function.String())
	sb.WriteString(") on (")
	for i, a := range n.arguments {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(a.String())
	}
	sb.WriteString(")")
	return sb.String()
}

func (n NodeFunctionCall) Position() position {
	return n.function.Position()
}

type NodeMatchClause struct {
	target, expression Node
}

func (n NodeMatchClause) String() string {
	return fmt.Sprintf("Clause (%s) -> (%s)", n.target, n.expression)
}

func (n NodeMatchClause) Position() position {
	return n.target.Position()
}

type NodeMatchExpr struct {
	position
	condition Node
	clauses   []NodeMatchClause
}

func (n NodeMatchExpr) String() string {
	var sb strings.Builder
	sb.WriteString("Match on (")
	sb.WriteString(n.condition.String())
	sb.WriteString(") to {")
	for i, a := range n.clauses {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(a.String())
	}
	sb.WriteString("}")
	return sb.String()
}

func (n NodeMatchExpr) Position() position {
	return n.position
}

type NodeExprList struct {
	position
	expressions []Node
}

func (n NodeExprList) String() string {
	var sb strings.Builder
	sb.WriteString("Expression list (")
	for i, expr := range n.expressions {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(expr.String())
	}
	sb.WriteString(")")
	return sb.String()
}

func (n NodeExprList) Position() position {
	return n.position
}

type NodeIdentifierEmpty struct {
	position
}

func (n NodeIdentifierEmpty) String() string {
	return "Empty Identifier"
}

func (n NodeIdentifierEmpty) Position() position {
	return n.position
}

type NodeIdentifier struct {
	position
	val string
}

func (n NodeIdentifier) String() string {
	return fmt.Sprintf("Identifier '%s'", n.val)
}

func (n NodeIdentifier) Position() position {
	return n.position
}

type NodeLiteralNumber struct {
	position
	val float64
}

func (n NodeLiteralNumber) String() string {
	return fmt.Sprintf("Number %s", nToS(n.val))
}

func (n NodeLiteralNumber) Position() position {
	return n.position
}

type NodeLiteralString struct {
	position
	val string
}

func (n NodeLiteralString) String() string {
	return fmt.Sprintf("String '%s'", n.val)
}

func (n NodeLiteralString) Position() position {
	return n.position
}

type NodeLiteralBoolean struct {
	position
	val bool
}

func (n NodeLiteralBoolean) String() string {
	return fmt.Sprintf("Boolean %t", n.val)
}

func (n NodeLiteralBoolean) Position() position {
	return n.position
}

type NodeLiteralObject struct {
	position
	entries []NodeObjectEntry
}

func (n NodeLiteralObject) String() string {
	entries := make([]string, len(n.entries))
	for i, e := range n.entries {
		entries[i] = e.String()
	}
	return fmt.Sprintf("Object {%s}",
		strings.Join(entries, ", "))
}

func (n NodeLiteralObject) Position() position {
	return n.position
}

type NodeObjectEntry struct {
	position
	key, val Node
}

func (n NodeObjectEntry) String() string {
	return fmt.Sprintf("(%s): (%s)", n.key, n.val)
}

type NodeLiteralList struct {
	position
	vals []Node
}

func (n NodeLiteralList) String() string {
	vals := make([]string, len(n.vals))
	for i, v := range n.vals {
		vals[i] = v.String()
	}
	return fmt.Sprintf("List [%s]", strings.Join(vals, ", "))
}

func (n NodeLiteralList) Position() position {
	return n.position
}

type NodeLiteralFunction struct {
	position
	arguments []Node
	body      Node
}

func (n NodeLiteralFunction) String() string {
	args := make([]string, len(n.arguments))
	for i, a := range n.arguments {
		args[i] = a.String()
	}
	return fmt.Sprintf("Function (%s) => (%s)", strings.Join(args, ", "), n.body)
}

func (n NodeLiteralFunction) Position() position {
	return n.position
}

func guardUnexpectedInputEnd(tokens []Token, idx int) *Err {
	if idx < len(tokens) {
		return nil
	}

	if len(tokens) == 0 {
		return &Err{ErrSyntax, fmt.Sprintf("unexpected end of input"), position{}} // TODO: report filename and position
	}

	return &Err{ErrSyntax, fmt.Sprintf("unexpected end of input at %s", tokens[len(tokens)-1]), tokens[len(tokens)-1].position}
}

// parse concurrently transforms a stream of Tok (tokens) to Node (AST nodes).
// This implementation uses recursive descent parsing.
func parse(tokenStream <-chan Token) <-chan Node {
	tokens := make([]Token, 0)
	for tok := range tokenStream {
		tokens = append(tokens, tok)
	}

	nodes := make(chan Node)
	go func() {
		defer close(nodes)

		for i := 0; i < len(tokens); {
			if tokens[i].kind == Separator {
				// this sometimes happens when the repl receives comment inputs
				i++
				continue
			}

			expr, incr, err := parseExpression(tokens[i:])
			if err != nil {
				LogError(err)
				return
			}

			i += incr

			LogNode(expr)
			nodes <- expr
		}
	}()
	return nodes
}

var opPriority = map[Kind]int{
	OpAccessor: 100,
	OpModulus:  80,

	OpMultiply: 50, OpDivide: 50,
	OpAdd: 40, OpSubtract: 40,

	OpGreaterThan: 30, OpLessThan: 30, OpEqual: 30,

	OpLogicalAnd: 20,
	OpLogicalXor: 15,
	OpLogicalOr:  10,

	OpDefine: 0,
}

func getOpPriority(t Token) int {
	// higher == greater priority
	priority, ok := opPriority[t.kind]
	if !ok {
		return -1
	}
	return priority
}

func isBinaryOp(t Token) bool {
	switch t.kind {
	case OpAdd, OpSubtract, OpMultiply, OpDivide, OpModulus,
		OpLogicalAnd, OpLogicalOr, OpLogicalXor,
		OpGreaterThan, OpLessThan, OpEqual, OpDefine, OpAccessor:
		return true
	default:
		return false
	}
}

func parseBinaryExpression(
	leftOperand Node,
	operator Token,
	tokens []Token,
	previousPriority int,
) (Node, int, *Err) {
	rightAtom, idx, err := parseAtom(tokens)
	if err != nil {
		return nil, 0, err
	}

	incr := 0
	ops := make([]Token, 1)
	nodes := make([]Node, 2)
	ops[0] = operator
	nodes[0] = leftOperand
	nodes[1] = rightAtom
	// build up a list of binary operations, with tree nodes
	// where there are higher-priority binary ops
LOOP:
	for len(tokens) > idx && isBinaryOp(tokens[idx]) {
		switch {
		case previousPriority >= getOpPriority(tokens[idx]):
			// Priority is lower than the calling function's last op,
			//  so return control to the parent binary op
			break LOOP
		case getOpPriority(ops[len(ops)-1]) >= getOpPriority(tokens[idx]):
			// Priority is lower than the previous op (but higher than parent),
			// so it's ok to be left-heavy in this tree
			ops = append(ops, tokens[idx])
			idx++

			err := guardUnexpectedInputEnd(tokens, idx)
			if err != nil {
				return nil, 0, err
			}

			rightAtom, incr, err = parseAtom(tokens[idx:])
			if err != nil {
				return nil, 0, err
			}
			nodes = append(nodes, rightAtom)
			idx += incr
		default:
			err := guardUnexpectedInputEnd(tokens, idx+1)
			if err != nil {
				return nil, 0, err
			}

			// Priority is higher than previous ops,
			// so make it a right-heavy tree
			subtree, incr, err := parseBinaryExpression(
				nodes[len(nodes)-1],
				tokens[idx],
				tokens[idx+1:],
				getOpPriority(ops[len(ops)-1]),
			)
			if err != nil {
				return nil, 0, err
			}
			nodes[len(nodes)-1] = subtree
			idx += incr + 1
		}
	}

	// ops, nodes -> left-biased binary expression tree
	tree := nodes[0]
	nodes = nodes[1:]
	for len(ops) > 0 {
		tree = NodeExprBinary{
			operator: ops[0].kind,
			left:     tree,
			right:    nodes[0],
			position: ops[0].position,
		}
		ops = ops[1:]
		nodes = nodes[1:]
	}

	return tree, idx, nil
}

func parseExpression(tokens []Token) (Node, int, *Err) {
	idx := 0
	consumeDanglingSeparator := func() {
		// bounds check in case parseExpress() called at some point
		// consumed end token
		if idx < len(tokens) && tokens[idx].kind == Separator {
			idx++
		}
	}

	atom, incr, err := parseAtom(tokens[idx:])
	if err != nil {
		return nil, 0, err
	}
	idx += incr

	if err = guardUnexpectedInputEnd(tokens, idx); err != nil {
		return nil, 0, err
	}

	nextTok := tokens[idx]
	idx++

	switch nextTok.kind {
	case Separator:
		// consuming dangling separator
		return atom, idx, nil

	case ParenRight, KeyValueSeparator, CaseArrow:
		// these belong to the parent atom that contains this expression,
		// so return without consuming token (idx - 1)
		return atom, idx - 1, nil

	case OpAdd, OpSubtract, OpMultiply, OpDivide, OpModulus,
		OpLogicalAnd, OpLogicalOr, OpLogicalXor,
		OpGreaterThan, OpLessThan, OpEqual, OpDefine, OpAccessor:
		binExpr, incr, err := parseBinaryExpression(atom, nextTok, tokens[idx:], -1)
		if err != nil {
			return nil, 0, err
		}
		idx += incr

		// Binary expressions are often followed by a match
		if idx < len(tokens) && tokens[idx].kind == MatchColon {
			colonPos := tokens[idx].position
			idx++ // MatchColon

			clauses, incr, err := parseMatchBody(tokens[idx:])
			if err != nil {
				return nil, 0, err
			}
			idx += incr

			consumeDanglingSeparator()
			return NodeMatchExpr{
				condition: binExpr,
				clauses:   clauses,
				position:  colonPos,
			}, idx, nil
		}

		consumeDanglingSeparator()
		return binExpr, idx, nil
	case MatchColon:
		clauses, incr, err := parseMatchBody(tokens[idx:])
		if err != nil {
			return nil, 0, err
		}
		idx += incr

		consumeDanglingSeparator()
		return NodeMatchExpr{
			condition: atom,
			clauses:   clauses,
			position:  nextTok.position,
		}, idx, nil
	default:
		return nil, 0, &Err{ErrSyntax, fmt.Sprintf("unexpected token %s following an expression", nextTok), nextTok.position}
	}
}

func parseAtom(tokens []Token) (Node, int, *Err) {
	if err := guardUnexpectedInputEnd(tokens, 0); err != nil {
		return nil, 0, err
	}

	tok, idx := tokens[0], 1

	if tok.kind == OpNegation {
		atom, idx, err := parseAtom(tokens[idx:])
		if err != nil {
			return nil, 0, err
		}
		return NodeExprUnary{
			operator: tok.kind,
			operand:  atom,
			position: tok.position,
		}, idx + 1, nil
	}

	if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
		return nil, 0, err
	}

	var atom Node
	switch tok.kind {
	case LiteralNumber:
		return NodeLiteralNumber{tok.position, tok.num}, idx, nil
	case LiteralString:
		return NodeLiteralString{tok.position, tok.str}, idx, nil
	case LiteralTrue:
		return NodeLiteralBoolean{tok.position, true}, idx, nil
	case LiteralFalse:
		return NodeLiteralBoolean{tok.position, false}, idx, nil
	case Identifier:
		if tokens[idx].kind == FunctionArrow {
			var err *Err
			atom, idx, err = parseFunctionLiteral(tokens)
			if err != nil {
				return nil, 0, err
			}

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			idx--
		} else {
			atom = NodeIdentifier{tok.position, tok.str}
		}
		// may be called as a function, so flows beyond
		// switch block
	case IdentifierEmpty:
		if tokens[idx].kind == FunctionArrow {
			var err *Err
			atom, idx, err = parseFunctionLiteral(tokens)
			if err != nil {
				return nil, 0, err
			}

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			return atom, idx - 1, nil
		}

		return NodeIdentifierEmpty{tok.position}, idx, nil
	case ParenLeft:
		// grouped expression or function literal
		exprs := make([]Node, 0)
		for tokens[idx].kind != ParenRight {
			expr, incr, err := parseExpression(tokens[idx:])
			if err != nil {
				return nil, 0, err
			}

			idx += incr
			exprs = append(exprs, expr)

			if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
				return nil, 0, err
			}
		}
		idx++ // RightParen

		if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
			return nil, 0, err
		}

		if tokens[idx].kind == FunctionArrow {
			var err *Err
			atom, idx, err = parseFunctionLiteral(tokens)
			if err != nil {
				return nil, 0, err
			}

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			idx--
		} else {
			atom = NodeExprList{
				expressions: exprs,
				position:    tok.position,
			}
		}
		// may be called as a function, so flows beyond
		// switch block
	case BraceLeft:
		entries := make([]NodeObjectEntry, 0)
		for tokens[idx].kind != BraceRight {
			keyExpr, keyIncr, err := parseExpression(tokens[idx:])
			if err != nil {
				return nil, 0, err
			}

			idx += keyIncr

			if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
				return nil, 0, err
			}

			if tokens[idx].kind != KeyValueSeparator {
				return nil, 0, &Err{ErrSyntax, fmt.Sprintf("expected %s after composite key, found %s", KeyValueSeparator.String(), tokens[idx]), tok.position}
			}

			idx++

			if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
				return nil, 0, err
			}

			valExpr, valIncr, err := parseExpression(tokens[idx:])
			if err != nil {
				return nil, 0, err
			}

			// Separator consumed by parseExpression
			idx += valIncr
			entries = append(entries, NodeObjectEntry{
				key:      keyExpr,
				val:      valExpr,
				position: keyExpr.Position(),
			})

			if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
				return nil, 0, err
			}
		}
		idx++ // RightBrace

		return NodeLiteralObject{
			entries:  entries,
			position: tok.position,
		}, idx, nil
	case BracketLeft:
		vals := make([]Node, 0)
		for tokens[idx].kind != BracketRight {
			expr, incr, err := parseExpression(tokens[idx:])
			if err != nil {
				return nil, 0, err
			}

			idx += incr
			vals = append(vals, expr)

			err = guardUnexpectedInputEnd(tokens, idx)
			if err != nil {
				return nil, 0, err
			}
		}
		idx++ // RightBracket

		return NodeLiteralList{
			vals:     vals,
			position: tok.position,
		}, idx, nil
	default:
		return nil, 0, &Err{ErrSyntax, fmt.Sprintf("unexpected start of atom, found %s", tok), tok.position}
	}

	// bounds check here because parseExpression may have
	// consumed all tokens before this
	for idx < len(tokens) && tokens[idx].kind == ParenLeft {
		var incr int
		var err *Err
		atom, incr, err = parseFunctionCall(atom, tokens[idx:])
		if err != nil {
			return nil, 0, err
		}

		idx += incr

		if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
			return nil, 0, err
		}
	}

	return atom, idx, nil
}

// parses everything that follows MatchColon
//
//	does not consume dangling separator -- that's for parseExpression
func parseMatchBody(tokens []Token) ([]NodeMatchClause, int, *Err) {
	idx := 1 // LeftBrace

	if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
		return nil, 0, err
	}

	clauses := make([]NodeMatchClause, 0)
	for tokens[idx].kind != BraceRight {
		clauseNode, incr, err := parseMatchClause(tokens[idx:])
		if err != nil {
			return nil, 0, err
		}
		idx += incr

		clauses = append(clauses, clauseNode)

		if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
			return nil, 0, err
		}
	}
	idx++ // RightBrace

	return clauses, idx, nil
}

func parseMatchClause(tokens []Token) (NodeMatchClause, int, *Err) {
	atom, idx, err := parseExpression(tokens)
	if err != nil {
		return NodeMatchClause{}, 0, err
	}

	if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
		return NodeMatchClause{}, 0, err
	}

	if tokens[idx].kind != CaseArrow {
		return NodeMatchClause{}, 0, &Err{ErrSyntax, fmt.Sprintf("expected %s, but got %s", CaseArrow, tokens[idx]), tokens[idx].position}
	}
	idx++ // CaseArrow

	if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
		return NodeMatchClause{}, 0, err
	}

	expr, incr, err := parseExpression(tokens[idx:])
	if err != nil {
		return NodeMatchClause{}, 0, err
	}

	idx += incr

	return NodeMatchClause{
		target:     atom,
		expression: expr,
	}, idx, nil
}

func parseFunctionLiteral(tokens []Token) (NodeLiteralFunction, int, *Err) {
	tok, idx := tokens[0], 1

	if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
		return NodeLiteralFunction{}, 0, err
	}

	arguments := make([]Node, 0)
	switch tok.kind {
	case ParenLeft:
	LOOP:
		for {
			tk := tokens[idx]
			switch tk.kind {
			case Identifier:
				idNode := NodeIdentifier{tk.position, tk.str}
				arguments = append(arguments, idNode)
			case IdentifierEmpty:
				idNode := NodeIdentifierEmpty{tk.position}
				arguments = append(arguments, idNode)
			default:
				break LOOP
			}
			idx++

			if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
				return NodeLiteralFunction{}, 0, err
			}

			if tokens[idx].kind != Separator {
				return NodeLiteralFunction{}, 0, &Err{ErrSyntax, fmt.Sprintf("expected arguments in a list separated by %s, found %s", Separator, tokens[idx]), tokens[idx].position}
			}
			idx++ // Separator
		}

		if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
			return NodeLiteralFunction{}, 0, err
		}
		if tokens[idx].kind != ParenRight {
			return NodeLiteralFunction{}, 0, &Err{ErrSyntax, fmt.Sprintf("expected arguments list to terminate with %s, found %s", ParenRight, tokens[idx]), tokens[idx].position}
		}
		idx++ // RightParen
	case Identifier:
		idNode := NodeIdentifier{tok.position, tok.str}
		arguments = append(arguments, idNode)
	case IdentifierEmpty:
		idNode := NodeIdentifierEmpty{tok.position}
		arguments = append(arguments, idNode)
	default:
		return NodeLiteralFunction{}, 0, &Err{ErrSyntax, fmt.Sprintf("malformed arguments list in function at %s", tok), tok.position}
	}

	if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
		return NodeLiteralFunction{}, 0, err
	}

	if tokens[idx].kind != FunctionArrow {
		return NodeLiteralFunction{}, 0, &Err{ErrSyntax, fmt.Sprintf("expected %s but found %s", FunctionArrow, tokens[idx]), tokens[idx].position}
	}
	idx++ // FunctionArrow

	body, incr, err := parseExpression(tokens[idx:])
	if err != nil {
		return NodeLiteralFunction{}, 0, err
	}
	idx += incr

	return NodeLiteralFunction{
		arguments: arguments,
		body:      body,
		position:  tokens[0].position,
	}, idx, nil
}

func parseFunctionCall(function Node, tokens []Token) (NodeFunctionCall, int, *Err) {
	idx := 1

	if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
		return NodeFunctionCall{}, 0, err
	}

	arguments := make([]Node, 0)
	for tokens[idx].kind != ParenRight {
		expr, incr, err := parseExpression(tokens[idx:])
		if err != nil {
			return NodeFunctionCall{}, 0, err
		}

		idx += incr
		arguments = append(arguments, expr)

		if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
			return NodeFunctionCall{}, 0, err
		}
	}
	idx++ // RightParen

	return NodeFunctionCall{
		function:  function,
		arguments: arguments,
	}, idx, nil
}
