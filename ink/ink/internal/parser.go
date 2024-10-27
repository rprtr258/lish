package internal

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/rprtr258/fun"
)

// Node represents an abstract syntax tree (AST) node in an Ink program.
type Node interface {
	String() string
	Position() Pos
	Eval(*Scope, bool) (Value, *Err)
}

type NodeExprUnary struct {
	Pos
	operator Kind
	operand  Node
}

func (n NodeExprUnary) String() string {
	return fmt.Sprintf("Unary %s (%s)", n.operator, n.operand)
}

func (n NodeExprUnary) Position() Pos {
	return n.Pos
}

type NodeExprBinary struct {
	Pos
	operator    Kind
	left, right Node
}

func (n NodeExprBinary) String() string {
	return fmt.Sprintf("Binary (%s) %s (%s)", n.left, n.operator, n.right)
}

func (n NodeExprBinary) Position() Pos {
	return n.Pos
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

func (n NodeFunctionCall) Position() Pos {
	return n.function.Position()
}

type NodeMatchClause struct {
	target, expression Node
}

func (n NodeMatchClause) String() string {
	return fmt.Sprintf("Clause (%s) -> (%s)", n.target, n.expression)
}

func (n NodeMatchClause) Position() Pos {
	return n.target.Position()
}

type NodeMatchExpr struct {
	Pos
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

func (n NodeMatchExpr) Position() Pos {
	return n.Pos
}

type NodeExprList struct {
	Pos
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

func (n NodeExprList) Position() Pos {
	return n.Pos
}

type NodeIdentifierEmpty struct {
	Pos
}

func (n NodeIdentifierEmpty) String() string {
	return "Empty Identifier"
}

func (n NodeIdentifierEmpty) Position() Pos {
	return n.Pos
}

type NodeIdentifier struct {
	Pos
	val string
}

func (n NodeIdentifier) String() string {
	return fmt.Sprintf("Identifier '%s'", n.val)
}

func (n NodeIdentifier) Position() Pos {
	return n.Pos
}

type NodeLiteralNumber struct {
	Pos
	val float64
}

func (n NodeLiteralNumber) String() string {
	return fmt.Sprintf("Number %s", nToS(n.val))
}

func (n NodeLiteralNumber) Position() Pos {
	return n.Pos
}

type NodeLiteralString struct {
	Pos
	val string
}

func (n NodeLiteralString) String() string {
	return fmt.Sprintf("String '%s'", n.val)
}

func (n NodeLiteralString) Position() Pos {
	return n.Pos
}

type NodeLiteralBoolean struct {
	Pos
	val bool
}

func (n NodeLiteralBoolean) String() string {
	return fmt.Sprintf("Boolean %t", n.val)
}

func (n NodeLiteralBoolean) Position() Pos {
	return n.Pos
}

type NodeLiteralObject struct {
	Pos
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

func (n NodeLiteralObject) Position() Pos {
	return n.Pos
}

type NodeObjectEntry struct {
	Pos
	key, val Node
}

func (n NodeObjectEntry) String() string {
	return fmt.Sprintf("(%s): (%s)", n.key, n.val)
}

type NodeLiteralList struct {
	Pos
	vals []Node
}

func (n NodeLiteralList) String() string {
	vals := make([]string, len(n.vals))
	for i, v := range n.vals {
		vals[i] = v.String()
	}
	return fmt.Sprintf("List [%s]", strings.Join(vals, ", "))
}

func (n NodeLiteralList) Position() Pos {
	return n.Pos
}

type NodeLiteralFunction struct {
	Pos
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

func (n NodeLiteralFunction) Position() Pos {
	return n.Pos
}

func guardUnexpectedInputEnd(tokens []Token, idx int) *Err {
	if idx < len(tokens) {
		return nil
	}

	if len(tokens) == 0 {
		return &Err{nil, ErrSyntax, fmt.Sprintf("unexpected end of input"), Pos{}} // TODO: report filename and position
	}

	return &Err{nil, ErrSyntax, fmt.Sprintf("unexpected end of input at %s", tokens[len(tokens)-1]), tokens[len(tokens)-1].Pos}
}

// parse concurrently transforms a stream of Tok (tokens) to Node (AST nodes).
// This implementation uses recursive descent parsing.
func parse(tokenStream iter.Seq[Token]) iter.Seq[Node] {
	// TODO: parse stream if we can, hence making "one-pass" interpreter
	tokens := slices.Collect(tokenStream)

	return func(yield func(Node) bool) {
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
			if !yield(expr) {
				return
			}
		}
	}
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
	return fun.Contains(t.kind,
		OpAdd, OpSubtract, OpMultiply, OpDivide, OpModulus,
		OpLogicalAnd, OpLogicalOr, OpLogicalXor,
		OpGreaterThan, OpLessThan, OpEqual, OpDefine, OpAccessor,
	)
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
	ops := []Token{operator}
	nodes := []Node{leftOperand, rightAtom}
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

			if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
				return nil, 0, err
			}

			rightAtom, incr, err = parseAtom(tokens[idx:])
			if err != nil {
				return nil, 0, err
			}

			nodes = append(nodes, rightAtom)
			idx += incr
		default:
			if err := guardUnexpectedInputEnd(tokens, idx+1); err != nil {
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
	for nodes := nodes[1:]; len(ops) > 0; nodes, ops = nodes[1:], ops[1:] {
		tree = NodeExprBinary{
			operator: ops[0].kind,
			left:     tree,
			right:    nodes[0],
			Pos:      ops[0].Pos,
		}
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
			colonPos := tokens[idx].Pos
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
				Pos:       colonPos,
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
			Pos:       nextTok.Pos,
		}, idx, nil
	default:
		return nil, 0, &Err{nil, ErrSyntax, fmt.Sprintf("unexpected token %s following an expression", nextTok), nextTok.Pos}
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
			Pos:      tok.Pos,
		}, idx + 1, nil
	}

	if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
		return nil, 0, err
	}

	var atom Node
	switch tok.kind {
	case LiteralNumber:
		return NodeLiteralNumber{tok.Pos, tok.num}, idx, nil
	case LiteralString:
		return NodeLiteralString{tok.Pos, tok.str}, idx, nil
	case LiteralTrue:
		return NodeLiteralBoolean{tok.Pos, true}, idx, nil
	case LiteralFalse:
		return NodeLiteralBoolean{tok.Pos, false}, idx, nil
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
			atom = NodeIdentifier{tok.Pos, tok.str}
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

		return NodeIdentifierEmpty{tok.Pos}, idx, nil
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
				Pos:         tok.Pos,
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

			var valExpr Node
			if tokens[idx].kind == KeyValueSeparator { // "key: value" pair
				idx++

				if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
					return nil, 0, err
				}

				expr, valIncr, err := parseExpression(tokens[idx:])
				if err != nil {
					return nil, 0, err
				}

				valExpr = expr
				idx += valIncr // Separator consumed by parseExpression
			} else if _, ok := keyExpr.(NodeIdentifier); ok { // "key", shorthand for "key: key"
				valExpr = keyExpr
			} else {
				return nil, 0, &Err{nil, ErrSyntax, fmt.Sprintf("expected %s after composite key, found %s", KeyValueSeparator.String(), tokens[idx]), tok.Pos}
			}

			entries = append(entries, NodeObjectEntry{
				key: keyExpr,
				val: valExpr,
				Pos: keyExpr.Position(),
			})

			if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
				return nil, 0, err
			}
		}
		idx++ // RightBrace

		return NodeLiteralObject{
			entries: entries,
			Pos:     tok.Pos,
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

			if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
				return nil, 0, err
			}
		}
		idx++ // RightBracket

		return NodeLiteralList{
			vals: vals,
			Pos:  tok.Pos,
		}, idx, nil
	default:
		return nil, 0, &Err{nil, ErrSyntax, fmt.Sprintf("unexpected start of atom, found %s", tok), tok.Pos}
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
		return NodeMatchClause{}, 0, &Err{nil, ErrSyntax, fmt.Sprintf("expected %s, but got %s", CaseArrow, tokens[idx]), tokens[idx].Pos}
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
				idNode := NodeIdentifier{tk.Pos, tk.str}
				arguments = append(arguments, idNode)
			case IdentifierEmpty:
				idNode := NodeIdentifierEmpty{tk.Pos}
				arguments = append(arguments, idNode)
			default:
				break LOOP
			}
			idx++

			if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
				return NodeLiteralFunction{}, 0, err
			}

			if tokens[idx].kind != Separator {
				return NodeLiteralFunction{}, 0, &Err{nil, ErrSyntax, fmt.Sprintf("expected arguments in a list separated by %s, found %s", Separator, tokens[idx]), tokens[idx].Pos}
			}
			idx++ // Separator
		}

		if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
			return NodeLiteralFunction{}, 0, err
		}
		if tokens[idx].kind != ParenRight {
			return NodeLiteralFunction{}, 0, &Err{nil, ErrSyntax, fmt.Sprintf("expected arguments list to terminate with %s, found %s", ParenRight, tokens[idx]), tokens[idx].Pos}
		}
		idx++ // RightParen
	case Identifier:
		idNode := NodeIdentifier{tok.Pos, tok.str}
		arguments = append(arguments, idNode)
	case IdentifierEmpty:
		idNode := NodeIdentifierEmpty{tok.Pos}
		arguments = append(arguments, idNode)
	default:
		return NodeLiteralFunction{}, 0, &Err{nil, ErrSyntax, fmt.Sprintf("malformed arguments list in function at %s", tok), tok.Pos}
	}

	if err := guardUnexpectedInputEnd(tokens, idx); err != nil {
		return NodeLiteralFunction{}, 0, err
	}

	if tokens[idx].kind != FunctionArrow {
		return NodeLiteralFunction{}, 0, &Err{nil, ErrSyntax, fmt.Sprintf("expected %s but found %s", FunctionArrow, tokens[idx]), tokens[idx].Pos}
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
		Pos:       tokens[0].Pos,
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
