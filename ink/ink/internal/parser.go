package internal

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/rprtr258/fun"
)

type astSlice struct {
	numbers     map[float64]int
	identifiers map[string]int
	strings     map[string]int
	nodes       []Node
}

func newAstSlice() *astSlice {
	return &astSlice{
		numbers:     map[float64]int{},
		identifiers: map[string]int{},
		strings:     map[string]int{},
		nodes: []Node{
			NodeIdentifierEmpty{},
		},
	}
}

func (s *astSlice) append(node Node) Node {
	// TODO: position is lost
	switch node := node.(type) {
	case NodeIdentifierEmpty:
		// one empty identifier for all
		return s.nodes[0]
	case NodeLiteralNumber:
		if _, ok := s.numbers[node.val]; !ok {
			n := len(s.nodes)
			s.nodes = append(s.nodes, node)
			s.numbers[node.val] = n
		}
		return s.nodes[s.numbers[node.val]]
	case NodeLiteralString:
		if _, ok := s.strings[node.val]; !ok {
			n := len(s.nodes)
			s.nodes = append(s.nodes, node)
			s.strings[node.val] = n
		}
		return s.nodes[s.strings[node.val]]
	case NodeIdentifier:
		if _, ok := s.identifiers[node.val]; !ok {
			n := len(s.nodes)
			s.nodes = append(s.nodes, node)
			s.identifiers[node.val] = n
		}
		return s.nodes[s.identifiers[node.val]]
	default:
		n := len(s.nodes)
		s.nodes = append(s.nodes, node)
		return s.nodes[n]
	}
}

func (s astSlice) String() string {
	var sb strings.Builder
	if len(s.numbers) > 0 {
		fmt.Fprintln(&sb, "Numbers:")
		for x, i := range s.numbers {
			fmt.Fprintln(&sb, "\t", x, i)
		}
	}
	if len(s.strings) > 0 {
		fmt.Fprintln(&sb, "Strings:")
		for s, i := range s.strings {
			fmt.Fprintln(&sb, "\t", s, i)
		}
	}
	if len(s.identifiers) > 0 {
		fmt.Fprintln(&sb, "Identifiers:")
		for id, i := range s.identifiers {
			fmt.Fprintln(&sb, "\t", id, i)
		}
	}
	if len(s.nodes) > 0 {
		fmt.Fprintln(&sb, "AST:")
		for i, n := range s.nodes {
			fmt.Fprintln(&sb, "\t", i, n.String())
		}
	}
	return sb.String()
}

type errParse struct {
	err *Err
}

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

func guardUnexpectedInputEnd(tokens []Token, idx int) {
	if idx < len(tokens) {
		return
	}

	if len(tokens) == 0 {
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected end of input"), Pos{}}}) // TODO: report filename and position
	}

	panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected end of input at %s", tokens[len(tokens)-1]), tokens[len(tokens)-1].Pos}})
}

// parse concurrently transforms a stream of Tok (tokens) to Node (AST nodes).
// This implementation uses recursive descent parsing.
func parse(tokenStream iter.Seq[Token]) iter.Seq[Node] {
	// TODO: parse stream if we can, hence making "one-pass" interpreter
	tokens := slices.Collect(tokenStream)

	return func(yield func(Node) bool) {
		s := newAstSlice()
		for i := 0; i < len(tokens); {
			if tokens[i].kind == Separator {
				// this sometimes happens when the repl receives comment inputs
				i++
				continue
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						err, ok := r.(errParse)
						if !ok {
							panic(r)
						}

						LogError(err.err)
					}
				}()

				expr, incr := parseExpression(tokens[i:], s)

				i += incr

				logNode(expr)
				if !yield(expr) {
					return
				}
			}()
		}
		logAST(s)
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
	s *astSlice,
) (Node, int) {
	rightAtom, idx := parseAtom(tokens, s)

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

			guardUnexpectedInputEnd(tokens, idx)

			rightAtom, incr = parseAtom(tokens[idx:], s)

			nodes = append(nodes, rightAtom)
			idx += incr
		default:
			guardUnexpectedInputEnd(tokens, idx+1)

			// Priority is higher than previous ops,
			// so make it a right-heavy tree
			subtree, incr := parseBinaryExpression(
				nodes[len(nodes)-1],
				tokens[idx],
				tokens[idx+1:],
				getOpPriority(ops[len(ops)-1]),
				s,
			)

			nodes[len(nodes)-1] = subtree
			idx += incr + 1
		}
	}

	// ops, nodes -> left-biased binary expression tree
	tree := nodes[0]
	for nodes := nodes[1:]; len(ops) > 0; nodes, ops = nodes[1:], ops[1:] {
		tree = s.append(NodeExprBinary{
			operator: ops[0].kind,
			left:     tree,
			right:    nodes[0],
			Pos:      ops[0].Pos,
		})
	}
	return tree, idx
}

func parseExpression(tokens []Token, s *astSlice) (Node, int) {
	idx := 0
	consumeDanglingSeparator := func() {
		// bounds check in case parseExpress() called at some point
		// consumed end token
		if idx < len(tokens) && tokens[idx].kind == Separator {
			idx++
		}
	}

	atom, incr := parseAtom(tokens[idx:], s)

	idx += incr

	guardUnexpectedInputEnd(tokens, idx)

	nextTok := tokens[idx]
	idx++

	switch nextTok.kind {
	case Separator:
		// consuming dangling separator
		return atom, idx
	case ParenRight, KeyValueSeparator, CaseArrow:
		// these belong to the parent atom that contains this expression,
		// so return without consuming token (idx - 1)
		return atom, idx - 1
	case OpAdd, OpSubtract, OpMultiply, OpDivide, OpModulus,
		OpLogicalAnd, OpLogicalOr, OpLogicalXor,
		OpGreaterThan, OpLessThan, OpEqual, OpDefine, OpAccessor:
		binExpr, incr := parseBinaryExpression(atom, nextTok, tokens[idx:], -1, s)
		idx += incr

		// Binary expressions are often followed by a match
		// TODO: support empty match expression ((true by default) :: {n < 1 -> ...})
		if idx < len(tokens) && tokens[idx].kind == MatchColon {
			colonPos := tokens[idx].Pos
			idx++ // MatchColon

			clauses, incr := parseMatchBody(tokens[idx:], s)
			idx += incr

			consumeDanglingSeparator()
			return s.append(NodeMatchExpr{
				condition: binExpr,
				clauses:   clauses,
				Pos:       colonPos,
			}), idx
		}

		consumeDanglingSeparator()
		return binExpr, idx
	case MatchColon:
		clauses, incr := parseMatchBody(tokens[idx:], s)

		idx += incr

		consumeDanglingSeparator()
		return s.append(NodeMatchExpr{
			condition: atom,
			clauses:   clauses,
			Pos:       nextTok.Pos,
		}), idx
	default:
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected token %s following an expression", nextTok), nextTok.Pos}})
	}
}

func parseAtom(tokens []Token, s *astSlice) (Node, int) {
	guardUnexpectedInputEnd(tokens, 0)

	tok, idx := tokens[0], 1

	if tok.kind == OpNegation {
		atom, idx := parseAtom(tokens[idx:], s)
		return s.append(NodeExprUnary{
			operator: tok.kind,
			operand:  atom,
			Pos:      tok.Pos,
		}), idx + 1
	}

	guardUnexpectedInputEnd(tokens, idx)

	var atom Node
	switch tok.kind {
	case LiteralNumber:
		return s.append(NodeLiteralNumber{tok.Pos, tok.num}), idx
	case LiteralString:
		return s.append(NodeLiteralString{tok.Pos, tok.str}), idx
	case LiteralTrue:
		return s.append(NodeLiteralBoolean{tok.Pos, true}), idx
	case LiteralFalse:
		return s.append(NodeLiteralBoolean{tok.Pos, false}), idx
	case Identifier:
		if tokens[idx].kind == FunctionArrow {
			atom, idx = parseFunctionLiteral(tokens, s)

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			idx--
		} else {
			atom = s.append(NodeIdentifier{tok.Pos, tok.str})
		}
		// may be called as a function, so flows beyond
		// switch block
	case IdentifierEmpty:
		if tokens[idx].kind == FunctionArrow {
			atom, idx = parseFunctionLiteral(tokens, s)

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			return atom, idx - 1
		}

		return s.append(NodeIdentifierEmpty{tok.Pos}), idx
	case ParenLeft:
		// grouped expression or function literal
		exprs := make([]Node, 0)
		for tokens[idx].kind != ParenRight {
			expr, incr := parseExpression(tokens[idx:], s)

			idx += incr
			exprs = append(exprs, expr)

			guardUnexpectedInputEnd(tokens, idx)
		}
		idx++ // RightParen

		guardUnexpectedInputEnd(tokens, idx)

		if tokens[idx].kind == FunctionArrow {
			atom, idx = parseFunctionLiteral(tokens, s)

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			idx--
		} else {
			atom = s.append(NodeExprList{
				expressions: exprs,
				Pos:         tok.Pos,
			})
		}
		// may be called as a function, so flows beyond switch block
	case BraceLeft:
		entries := make([]NodeObjectEntry, 0)
		for tokens[idx].kind != BraceRight {
			keyExpr, keyIncr := parseExpression(tokens[idx:], s)
			idx += keyIncr

			guardUnexpectedInputEnd(tokens, idx)

			var valExpr Node
			if tokens[idx].kind == KeyValueSeparator { // "key: value" pair
				idx++

				guardUnexpectedInputEnd(tokens, idx)

				expr, valIncr := parseExpression(tokens[idx:], s)

				valExpr = expr
				idx += valIncr // Separator consumed by parseExpression
			} else if _, ok := keyExpr.(NodeIdentifier); ok { // "key", shorthand for "key: key"
				valExpr = keyExpr
			} else {
				panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %s after composite key, found %s", KeyValueSeparator.String(), tokens[idx]), tok.Pos}})
			}

			entries = append(entries, NodeObjectEntry{
				key: keyExpr,
				val: valExpr,
				Pos: keyExpr.Position(),
			})

			guardUnexpectedInputEnd(tokens, idx)
		}
		idx++ // RightBrace

		return s.append(NodeLiteralObject{
			entries: entries,
			Pos:     tok.Pos,
		}), idx
	case BracketLeft:
		vals := make([]Node, 0)
		for tokens[idx].kind != BracketRight {
			expr, incr := parseExpression(tokens[idx:], s)

			idx += incr
			vals = append(vals, expr)

			guardUnexpectedInputEnd(tokens, idx)
		}
		idx++ // RightBracket

		return s.append(NodeLiteralList{
			vals: vals,
			Pos:  tok.Pos,
		}), idx
	default:
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected start of atom, found %s", tok), tok.Pos}})
	}

	// bounds check here because parseExpression may have
	// consumed all tokens before this
	for idx < len(tokens) && tokens[idx].kind == ParenLeft {
		var incr int
		atom, incr = parseFunctionCall(atom, tokens[idx:], s)

		idx += incr

		guardUnexpectedInputEnd(tokens, idx)
	}

	return atom, idx
}

// parses everything that follows MatchColon
//
//	does not consume dangling separator -- that's for parseExpression
func parseMatchBody(tokens []Token, s *astSlice) ([]NodeMatchClause, int) {
	idx := 1 // LeftBrace

	guardUnexpectedInputEnd(tokens, idx)

	clauses := make([]NodeMatchClause, 0)
	for tokens[idx].kind != BraceRight {
		clauseNode, incr := parseMatchClause(tokens[idx:], s)

		idx += incr
		clauses = append(clauses, clauseNode)

		guardUnexpectedInputEnd(tokens, idx)
	}
	idx++ // RightBrace
	return clauses, idx
}

func parseMatchClause(tokens []Token, s *astSlice) (NodeMatchClause, int) {
	atom, idx := parseExpression(tokens, s)

	guardUnexpectedInputEnd(tokens, idx)

	if tokens[idx].kind != CaseArrow {
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %s, but got %s", CaseArrow, tokens[idx]), tokens[idx].Pos}})
	}
	idx++ // CaseArrow

	guardUnexpectedInputEnd(tokens, idx)

	expr, incr := parseExpression(tokens[idx:], s)

	idx += incr

	return s.append(NodeMatchClause{
		target:     atom,
		expression: expr,
	}).(NodeMatchClause), idx
}

func parseFunctionLiteral(tokens []Token, s *astSlice) (NodeLiteralFunction, int) {
	tok, idx := tokens[0], 1

	guardUnexpectedInputEnd(tokens, idx)

	arguments := make([]Node, 0)
	switch tok.kind {
	case ParenLeft:
	LOOP:
		for {
			tk := tokens[idx]
			switch tk.kind {
			case Identifier:
				idNode := s.append(NodeIdentifier{tk.Pos, tk.str})
				arguments = append(arguments, idNode)
			case IdentifierEmpty:
				idNode := s.append(NodeIdentifierEmpty{tk.Pos})
				arguments = append(arguments, idNode)
			default:
				break LOOP
			}
			idx++

			guardUnexpectedInputEnd(tokens, idx)

			if tokens[idx].kind != Separator {
				panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected arguments in a list separated by %s, found %s", Separator, tokens[idx]), tokens[idx].Pos}})
			}
			idx++ // Separator
		}

		guardUnexpectedInputEnd(tokens, idx)

		if tokens[idx].kind != ParenRight {
			panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected arguments list to terminate with %s, found %s", ParenRight, tokens[idx]), tokens[idx].Pos}})
		}
		idx++ // RightParen
	case Identifier:
		idNode := s.append(NodeIdentifier{tok.Pos, tok.str})
		arguments = append(arguments, idNode)
	case IdentifierEmpty:
		idNode := s.append(NodeIdentifierEmpty{tok.Pos})
		arguments = append(arguments, idNode)
	default:
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("malformed arguments list in function at %s", tok), tok.Pos}})
	}

	guardUnexpectedInputEnd(tokens, idx)

	if tokens[idx].kind != FunctionArrow {
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %s but found %s", FunctionArrow, tokens[idx]), tokens[idx].Pos}})
	}

	idx++ // FunctionArrow

	body, incr := parseExpression(tokens[idx:], s)

	idx += incr

	return s.append(NodeLiteralFunction{
		arguments: arguments,
		body:      body,
		Pos:       tokens[0].Pos,
	}).(NodeLiteralFunction), idx
}

func parseFunctionCall(function Node, tokens []Token, s *astSlice) (NodeFunctionCall, int) {
	idx := 1

	guardUnexpectedInputEnd(tokens, idx)

	arguments := make([]Node, 0)
	for tokens[idx].kind != ParenRight {
		expr, incr := parseExpression(tokens[idx:], s)

		idx += incr
		arguments = append(arguments, expr)

		guardUnexpectedInputEnd(tokens, idx)
	}

	idx++ // RightParen

	return s.append(NodeFunctionCall{
		function:  function,
		arguments: arguments,
	}).(NodeFunctionCall), idx
}
