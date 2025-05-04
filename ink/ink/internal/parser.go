package internal

import (
	"cmp"
	"fmt"
	"iter"
	"slices"

	"github.com/rprtr258/fun"
)

func guardUnexpectedInputEnd(tokens []Token, idx int) {
	switch {
	case idx < len(tokens):
		return
	case len(tokens) == 0:
		panic(errParse{&Err{nil, ErrSyntax, "unexpected end of input", Pos{}}}) // TODO: report filename and position
	default:
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected end of input at %s", tokens[len(tokens)-1]), tokens[len(tokens)-1].Pos}})
	}
}

// parse concurrently transforms a stream of Tok (tokens) to Node (AST nodes).
// This implementation uses recursive descent parsing.
func parse(s *AST, tokenStream iter.Seq[Token]) (nodes []NodeID, e errParse) {
	// TODO: parse stream if we can, hence making "one-pass" interpreter
	tokens := slices.Collect(tokenStream)

	for i := 0; i < len(tokens); {
		if tokens[i].Kind == Separator {
			// this sometimes happens when the repl receives comment inputs
			i++
			continue
		}

		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(errParse)
				if !ok {
					panic(r)
				}

				e = err
			}
		}()

		expr, incr := parseExpression(tokens[i:], s)

		i += incr

		LogNode(s.Nodes[expr])
		nodes = append(nodes, expr)
	}
	LogAST(s)
	return
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

	OpDefine: 1,
}

func getOpPriority(t Token) int {
	// higher == greater priority
	return cmp.Or(opPriority[t.Kind], -1)
}

func isBinaryOp(t Token) bool {
	return fun.Contains(t.Kind,
		OpAdd, OpSubtract, OpMultiply, OpDivide, OpModulus,
		OpLogicalAnd, OpLogicalOr, OpLogicalXor,
		OpGreaterThan, OpLessThan, OpEqual, OpDefine, OpAccessor,
	)
}

func parseBinaryExpression(
	leftOperand NodeID,
	operator Token,
	tokens []Token,
	previousPriority int,
	s *AST,
) (NodeID, int) {
	rightAtom, idx := parseAtom(tokens, s)

	incr := 0
	ops := []Token{operator}
	nodes := []NodeID{leftOperand, rightAtom}
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
		tree = s.Append(NodeExprBinary{
			Operator: ops[0].Kind,
			Left:     tree,
			Right:    nodes[0],
			Pos:      ops[0].Pos,
		})
	}
	return tree, idx
}

func parseExpression(tokens []Token, s *AST) (NodeID, int) {
	idx := 0
	consumeDanglingSeparator := func() {
		// bounds check in case parseExpress() called at some point
		// consumed end token
		if idx < len(tokens) && tokens[idx].Kind == Separator {
			idx++
		}
	}

	atom, incr := parseAtom(tokens[idx:], s)

	idx += incr

	guardUnexpectedInputEnd(tokens, idx)

	nextTok := tokens[idx]
	idx++

	switch nextTok.Kind {
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
		if idx < len(tokens) && tokens[idx].Kind == MatchColon {
			colonPos := tokens[idx].Pos
			idx++ // MatchColon

			clauses, incr := parseMatchBody(tokens[idx:], s)
			idx += incr

			consumeDanglingSeparator()
			return s.Append(NodeExprMatch{
				Condition: binExpr,
				Clauses:   clauses,
				Pos:       colonPos,
			}), idx
		}

		consumeDanglingSeparator()
		return binExpr, idx
	case MatchColon:
		clauses, incr := parseMatchBody(tokens[idx:], s)

		idx += incr

		consumeDanglingSeparator()
		return s.Append(NodeExprMatch{
			Condition: atom,
			Clauses:   clauses,
			Pos:       nextTok.Pos,
		}), idx
	default:
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected token %s following an expression", nextTok), nextTok.Pos}})
	}
}

func parseAtom(tokens []Token, s *AST) (NodeID, int) {
	guardUnexpectedInputEnd(tokens, 0)

	tok, idx := tokens[0], 1

	if tok.Kind == OpNegation {
		atom, idx := parseAtom(tokens[idx:], s)
		return s.Append(NodeExprUnary{
			Operator: tok.Kind,
			Operand:  atom,
			Pos:      tok.Pos,
		}), idx + 1
	}

	guardUnexpectedInputEnd(tokens, idx)

	var atom NodeID
	switch tok.Kind {
	case LiteralNumber:
		return s.Append(NodeLiteralNumber{tok.Pos, tok.Num}), idx
	case LiteralString:
		return s.Append(NodeLiteralString{tok.Pos, tok.Str}), idx
	case LiteralTrue:
		return s.Append(NodeLiteralBoolean{tok.Pos, true}), idx
	case LiteralFalse:
		return s.Append(NodeLiteralBoolean{tok.Pos, false}), idx
	case Identifier:
		if tokens[idx].Kind == FunctionArrow {
			atom, idx = parseFunctionLiteral(tokens, s)

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			idx--
		} else {
			atom = s.Append(NodeIdentifier{tok.Pos, tok.Str})
		}
		// may be called as a function, so flows beyond
		// switch block
	case IdentifierEmpty:
		if tokens[idx].Kind == FunctionArrow {
			atom, idx = parseFunctionLiteral(tokens, s)

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			return atom, idx - 1
		}

		return s.Append(NodeIdentifierEmpty{tok.Pos}), idx
	case ParenLeft:
		// grouped expression or function literal
		exprs := make([]NodeID, 0)
		for tokens[idx].Kind != ParenRight {
			expr, incr := parseExpression(tokens[idx:], s)

			idx += incr
			exprs = append(exprs, expr)

			guardUnexpectedInputEnd(tokens, idx)
		}
		idx++ // RightParen

		guardUnexpectedInputEnd(tokens, idx)

		if tokens[idx].Kind == FunctionArrow {
			atom, idx = parseFunctionLiteral(tokens, s)

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			idx--
		} else {
			atom = s.Append(NodeExprList{
				Expressions: exprs,
				Pos:         tok.Pos,
			})
		}
		// may be called as a function, so flows beyond switch block
	case BraceLeft:
		entries := make([]NodeCompositeKeyValue, 0)
		for tokens[idx].Kind != BraceRight {
			keyExpr, keyIncr := parseExpression(tokens[idx:], s)
			idx += keyIncr

			guardUnexpectedInputEnd(tokens, idx)

			var valExpr NodeID
			if tokens[idx].Kind == KeyValueSeparator { // "key: value" pair
				idx++

				guardUnexpectedInputEnd(tokens, idx)

				expr, valIncr := parseExpression(tokens[idx:], s)

				valExpr = expr
				idx += valIncr // Separator consumed by parseExpression
			} else if _, ok := s.Nodes[keyExpr].(NodeIdentifier); ok { // "key", shorthand for "key: key"
				valExpr = keyExpr
			} else {
				panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %s after composite key, found %s", KeyValueSeparator.String(), tokens[idx]), tok.Pos}})
			}

			entries = append(entries, NodeCompositeKeyValue{
				Key: keyExpr,
				Val: valExpr,
				Pos: s.Nodes[keyExpr].Position(s),
			})

			guardUnexpectedInputEnd(tokens, idx)
		}
		idx++ // RightBrace

		return s.Append(NodeLiteralComposite{
			Entries: entries,
			Pos:     tok.Pos,
		}), idx
	case BracketLeft:
		vals := make([]NodeID, 0)
		for tokens[idx].Kind != BracketRight {
			expr, incr := parseExpression(tokens[idx:], s)

			idx += incr
			vals = append(vals, expr)

			guardUnexpectedInputEnd(tokens, idx)
		}
		idx++ // RightBracket

		return s.Append(NodeLiteralList{
			Vals: vals,
			Pos:  tok.Pos,
		}), idx
	default:
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected start of atom, found %s", tok), tok.Pos}})
	}

	// bounds check here because parseExpression may have
	// consumed all tokens before this
	for idx < len(tokens) && tokens[idx].Kind == ParenLeft {
		var incr int
		atom, incr = parseFunctionCall(tokens[idx:], s, atom)

		idx += incr

		guardUnexpectedInputEnd(tokens, idx)
	}

	return atom, idx
}

// parses everything that follows MatchColon
//
//	does not consume dangling separator -- that's for parseExpression
func parseMatchBody(tokens []Token, ast *AST) ([]NodeMatchClause, int) {
	idx := 1 // LeftBrace

	guardUnexpectedInputEnd(tokens, idx)

	clauses := make([]NodeMatchClause, 0)
	for tokens[idx].Kind != BraceRight {
		clauseNode, incr := parseMatchClause(tokens[idx:], ast)

		idx += incr
		clauses = append(clauses, clauseNode)

		guardUnexpectedInputEnd(tokens, idx)
	}
	idx++ // RightBrace

	// NOTE: adding _ -> () as last case
	if len(clauses) > 0 && !func() bool {
		_, ok := ast.Nodes[clauses[len(clauses)-1].Target].(NodeIdentifierEmpty)
		return ok
	}() {
		clauses = append(clauses, NodeMatchClause{ast.Append(NodeIdentifierEmpty{}), ast.Append(NodeLiteralList{})})
	}

	return clauses, idx
}

func parseMatchClause(tokens []Token, ast *AST) (NodeMatchClause, int) {
	atom, idx := parseExpression(tokens, ast)

	guardUnexpectedInputEnd(tokens, idx)

	if tokens[idx].Kind != CaseArrow {
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %s, but got %s", CaseArrow, tokens[idx]), tokens[idx].Pos}})
	}
	idx++ // CaseArrow

	guardUnexpectedInputEnd(tokens, idx)

	expr, incr := parseExpression(tokens[idx:], ast)

	idx += incr

	return NodeMatchClause{
		Target:     atom,
		Expression: expr,
	}, idx
}

func parseFunctionLiteral(tokens []Token, ast *AST) (NodeID, int) {
	tok, idx := tokens[0], 1

	guardUnexpectedInputEnd(tokens, idx)

	arguments := make([]NodeID, 0)
	switch tok.Kind {
	case ParenLeft:
	LOOP:
		for {
			tk := tokens[idx]
			switch tk.Kind {
			case Identifier:
				arguments = append(arguments, ast.Append(NodeIdentifier{tk.Pos, tk.Str}))
			case IdentifierEmpty:
				arguments = append(arguments, ast.Append(NodeIdentifierEmpty{tk.Pos}))
			default:
				break LOOP
			}
			idx++

			guardUnexpectedInputEnd(tokens, idx)

			if tokens[idx].Kind != Separator {
				panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected arguments in a list separated by %s, found %s", Separator, tokens[idx]), tokens[idx].Pos}})
			}
			idx++ // Separator
		}

		guardUnexpectedInputEnd(tokens, idx)

		if tokens[idx].Kind != ParenRight {
			panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected arguments list to terminate with %s, found %s", ParenRight, tokens[idx]), tokens[idx].Pos}})
		}
		idx++ // RightParen
	case Identifier:
		arguments = append(arguments, ast.Append(NodeIdentifier{tok.Pos, tok.Str}))
	case IdentifierEmpty:
		arguments = append(arguments, ast.Append(NodeIdentifierEmpty{tok.Pos}))
	default:
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("malformed arguments list in function at %s", tok), tok.Pos}})
	}

	guardUnexpectedInputEnd(tokens, idx)

	if tokens[idx].Kind != FunctionArrow {
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %s but found %s", FunctionArrow, tokens[idx]), tokens[idx].Pos}})
	}

	idx++ // FunctionArrow

	body, incr := parseExpression(tokens[idx:], ast)

	idx += incr

	return ast.Append(NodeLiteralFunction{
		Arguments: arguments,
		Body:      body,
		Pos:       tokens[0].Pos,
	}), idx
}

func parseFunctionCall(tokens []Token, s *AST, function NodeID) (NodeID, int) {
	idx := 1

	guardUnexpectedInputEnd(tokens, idx)

	arguments := make([]NodeID, 0)
	for tokens[idx].Kind != ParenRight {
		expr, incr := parseExpression(tokens[idx:], s)

		idx += incr
		arguments = append(arguments, expr)

		guardUnexpectedInputEnd(tokens, idx)
	}

	idx++ // RightParen

	return s.Append(NodeFunctionCall{function, arguments}), idx
}
