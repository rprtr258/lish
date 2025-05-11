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
func parse(ast *AST, tokenStream iter.Seq[Token]) (nodes []NodeID, e errParse) {
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

		expr, incr := parseExpression(tokens[i:], ast)

		i += incr

		LogNode(ast.Nodes[expr])
		nodes = append(nodes, expr)
	}
	LogAST(ast)
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
	tokens []Token,
	ast *AST,
	leftOperand NodeID,
	operator Token,
	previousPriority int,
) (NodeID, int) {
	rightAtom, idx := parseAtom(tokens, ast)

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

			rightAtom, incr = parseAtom(tokens[idx:], ast)

			nodes = append(nodes, rightAtom)
			idx += incr
		default:
			guardUnexpectedInputEnd(tokens, idx+1)

			// Priority is higher than previous ops,
			// so make it a right-heavy tree
			subtree, incr := parseBinaryExpression(
				tokens[idx+1:],
				ast,
				nodes[len(nodes)-1],
				tokens[idx],
				getOpPriority(ops[len(ops)-1]),
			)

			nodes[len(nodes)-1] = subtree
			idx += incr + 1
		}
	}

	// ops, nodes -> left-biased binary expression tree
	tree := nodes[0]
	for nodes := nodes[1:]; len(ops) > 0; nodes, ops = nodes[1:], ops[1:] {
		tree = ast.AppendOp(ops[0].Kind, ops[0].Pos, tree, nodes[0])
	}
	return tree, idx
}

func parseExpression(tokens []Token, ast *AST) (NodeID, int) {
	idx := 0
	consumeDanglingSeparator := func() {
		// bounds check in case parseExpress() called at some point
		// consumed end token
		if idx < len(tokens) && tokens[idx].Kind == Separator {
			idx++
		}
	}

	atom, incr := parseAtom(tokens[idx:], ast)

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
		binExpr, incr := parseBinaryExpression(tokens[idx:], ast, atom, nextTok, -1)
		idx += incr

		// Binary expressions are often followed by a match
		// TODO: support empty match expression ((true by default) :: {n < 1 -> ...})
		if idx < len(tokens) && tokens[idx].Kind == MatchColon {
			colonPos := tokens[idx].Pos
			idx++ // MatchColon

			clauses, incr := parseMatchBody(tokens[idx:], ast)
			idx += incr

			consumeDanglingSeparator()
			return ast.Append(NodeExprMatch{
				Condition: binExpr,
				Clauses:   clauses,
				Pos:       colonPos,
			}), idx
		}

		consumeDanglingSeparator()
		return binExpr, idx
	case MatchColon:
		clauses, incr := parseMatchBody(tokens[idx:], ast)

		idx += incr

		consumeDanglingSeparator()
		return ast.Append(NodeExprMatch{
			Condition: atom,
			Clauses:   clauses,
			Pos:       nextTok.Pos,
		}), idx
	default:
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected token %s following an expression", nextTok), nextTok.Pos}})
	}
}

func parseAtom(tokens []Token, ast *AST) (NodeID, int) {
	guardUnexpectedInputEnd(tokens, 0)

	tok, idx := tokens[0], 1

	if tok.Kind == OpNegation {
		atom, idx := parseAtom(tokens[idx:], ast)
		return ast.Append(NodeConstFunctionCall{
			Function:  operatorFunc(tok.Kind),
			Arguments: []NodeID{atom},
			Pos:       tok.Pos,
		}), idx + 1
	}

	guardUnexpectedInputEnd(tokens, idx)

	var atom NodeID
	switch tok.Kind {
	case LiteralNumber:
		return ast.Append(NodeLiteralNumber{tok.Pos, tok.Num}), idx
	case LiteralString:
		return ast.Append(NodeLiteralString{tok.Pos, tok.Str}), idx
	case LiteralTrue:
		return ast.Append(NodeLiteralBoolean{tok.Pos, true}), idx
	case LiteralFalse:
		return ast.Append(NodeLiteralBoolean{tok.Pos, false}), idx
	case Identifier:
		if tokens[idx].Kind == FunctionArrow {
			atom, idx = parseFunctionLiteral(tokens, ast)

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			idx--
		} else {
			atom = ast.Append(NodeIdentifier{tok.Pos, tok.Str})
		}
		// may be called as a function, so flows beyond
		// switch block
	case IdentifierEmpty:
		if tokens[idx].Kind == FunctionArrow {
			atom, idx = parseFunctionLiteral(tokens, ast)

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			return atom, idx - 1
		}

		return ast.Append(NodeIdentifierEmpty{tok.Pos}), idx
	case ParenLeft:
		// grouped expression or function literal
		exprs := make([]NodeID, 0)
		for tokens[idx].Kind != ParenRight {
			expr, incr := parseExpression(tokens[idx:], ast)

			idx += incr
			exprs = append(exprs, expr)

			guardUnexpectedInputEnd(tokens, idx)
		}
		idx++ // RightParen

		guardUnexpectedInputEnd(tokens, idx)

		if tokens[idx].Kind == FunctionArrow {
			atom, idx = parseFunctionLiteral(tokens, ast)

			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			idx--
		} else {
			atom = ast.Append(NodeExprList{
				Expressions: exprs,
				Pos:         tok.Pos,
			})
		}
		// may be called as a function, so flows beyond switch block
	case CurlyParenLeft:
		entries := make([]NodeCompositeKeyValue, 0)
		for tokens[idx].Kind != CurlyParenRight {
			keyExpr, keyIncr := parseExpression(tokens[idx:], ast)
			idx += keyIncr

			guardUnexpectedInputEnd(tokens, idx)

			var valExpr NodeID
			if tokens[idx].Kind == KeyValueSeparator { // "key: value" pair
				idx++

				guardUnexpectedInputEnd(tokens, idx)

				expr, valIncr := parseExpression(tokens[idx:], ast)

				valExpr = expr
				idx += valIncr // Separator consumed by parseExpression
			} else if _, ok := ast.Nodes[keyExpr].(NodeIdentifier); ok { // "key", shorthand for "key: key"
				valExpr = keyExpr
			} else {
				panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %s after composite key, found %s", KeyValueSeparator.String(), tokens[idx]), tok.Pos}})
			}

			entries = append(entries, NodeCompositeKeyValue{
				Key: keyExpr,
				Val: valExpr,
				Pos: ast.Nodes[keyExpr].Position(ast),
			})

			guardUnexpectedInputEnd(tokens, idx)
		}
		idx++ // RightBrace

		return ast.Append(NodeLiteralComposite{
			Entries: entries,
			Pos:     tok.Pos,
		}), idx
	case SquareParenLeft:
		vals := make([]NodeID, 0)
		for tokens[idx].Kind != SquareParenRight {
			expr, incr := parseExpression(tokens[idx:], ast)

			idx += incr
			vals = append(vals, expr)

			guardUnexpectedInputEnd(tokens, idx)
		}
		idx++ // RightBracket

		return ast.Append(NodeLiteralList{
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
		atom, incr = parseFunctionCall(tokens[idx:], ast, atom)

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
	for tokens[idx].Kind != CurlyParenRight {
		atom, incr1 := parseExpression(tokens[idx:], ast)

		guardUnexpectedInputEnd(tokens[idx:], incr1)

		if tokens[idx+incr1].Kind != CaseArrow {
			panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %s, but got %s", CaseArrow, tokens[idx+incr1]), tokens[idx+incr1].Pos}})
		}
		incr1++ // CaseArrow

		guardUnexpectedInputEnd(tokens[idx:], incr1)

		expr, incr2 := parseExpression(tokens[idx+incr1:], ast)

		idx += incr1 + incr2
		clauses = append(clauses, NodeMatchClause{
			Target:     atom,
			Expression: expr,
		})

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

func parseFunctionCall(tokens []Token, ast *AST, function NodeID) (NodeID, int) {
	idx := 1

	guardUnexpectedInputEnd(tokens, idx)

	arguments := make([]NodeID, 0)
	for tokens[idx].Kind != ParenRight {
		expr, incr := parseExpression(tokens[idx:], ast)

		idx += incr
		arguments = append(arguments, expr)

		guardUnexpectedInputEnd(tokens, idx)
	}

	idx++ // RightParen

	return ast.Append(NodeFunctionCall{function, arguments}), idx
}
