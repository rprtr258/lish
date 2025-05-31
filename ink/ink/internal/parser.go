package internal

import (
	"cmp"
	"fmt"

	"github.com/rprtr258/fun"
)

type parser struct {
	ast    *AST
	tokens []Token
	idx    int
}

func (p *parser) peek() Token {
	switch {
	case p.idx < len(p.tokens):
		tok := p.tokens[p.idx]
		return tok
	case len(p.tokens) == 0:
		panic(errParse{&Err{nil, ErrSyntax, "unexpected end of input", Pos{}}}) // TODO: report filename and position
	default:
		tok := p.tokens[len(p.tokens)-1]
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected end of input at %s", tok), tok.Pos}})
	}
}

func (p *parser) consume() Token {
	tok := p.peek()
	p.idx++
	return tok
}

func (p *parser) consumeKind(kind Kind) Token {
	tok := p.consume()
	if tok.Kind != kind {
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %s, but got %s", kind, tok), tok.Pos}})
	}
	return tok
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

func (p *parser) parseBinaryExpression(
	leftOperand NodeID,
	operator Token,
	previousPriority int,
) NodeID {
	ast, tokens := p.ast, p.tokens
	rightAtom := p.parseAtom()

	ops := []Token{operator}
	nodes := []NodeID{leftOperand, rightAtom}
	// build up a list of binary operations, with tree nodes
	// where there are higher-priority binary ops
LOOP:
	for len(tokens) > p.idx && isBinaryOp(tokens[p.idx]) {
		switch {
		case previousPriority >= getOpPriority(tokens[p.idx]):
			// Priority is lower than the calling function's last op,
			//  so return control to the parent binary op
			break LOOP
		case getOpPriority(ops[len(ops)-1]) >= getOpPriority(tokens[p.idx]):
			// Priority is lower than the previous op (but higher than parent),
			// so it's ok to be left-heavy in this tree
			op := p.consume()
			rightAtom = p.parseAtom()
			ops = append(ops, op)
			nodes = append(nodes, rightAtom)
		default:
			op := p.consume()
			// Priority is higher than previous ops,
			// so make it a right-heavy tree
			nodes[len(nodes)-1] = p.parseBinaryExpression(
				nodes[len(nodes)-1],
				op,
				getOpPriority(ops[len(ops)-1]),
			)
		}
	}

	// ops, nodes -> left-biased binary expression tree
	tree := nodes[0]
	for i := range len(ops) {
		tree = ast.Append(NodeExprBinary(ops[i].Pos, ops[i].Kind, tree, nodes[i+1]))
	}
	return tree
}

func (p *parser) consumeDanglingSeparator() {
	// bounds check in case parseExpression() called at some point consumed end token
	if p.idx < len(p.tokens) && p.tokens[p.idx].Kind == Separator {
		p.idx++
	}
}

func (p *parser) parseExpression() NodeID {
	atom := p.parseAtom()

	nextTok := p.consume()

	switch nextTok.Kind {
	case Separator:
		// consuming dangling separator
		return atom
	case ParenRight, KeyValueSeparator, CaseArrow:
		// these belong to the parent atom that contains this expression,
		// so return without consuming token
		p.idx--
		return atom
	case OpAdd, OpSubtract, OpMultiply, OpDivide, OpModulus,
		OpLogicalAnd, OpLogicalOr, OpLogicalXor,
		OpGreaterThan, OpLessThan, OpEqual, OpDefine, OpAccessor:
		binExpr := p.parseBinaryExpression(atom, nextTok, -1)

		// Binary expressions followed by a match
		// TODO: support empty match expression ((true by default) :: {n < 1 -> ...})
		if p.idx < len(p.tokens) && p.tokens[p.idx].Kind == MatchColon {
			colonPos := p.consumeKind(MatchColon).Pos
			clauses := p.parseMatchBody()
			p.consumeDanglingSeparator()
			return p.ast.Append(NodeExprMatch(colonPos, binExpr, clauses))
		}
		p.consumeDanglingSeparator()
		return binExpr
	case MatchColon:
		clauses := p.parseMatchBody()
		p.consumeDanglingSeparator()
		return p.ast.Append(NodeExprMatch(nextTok.Pos, atom, clauses))
	default:
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected token %s following an expression", nextTok), nextTok.Pos}})
	}
}

func (p *parser) parseAtom() NodeID {
	ast, tokens := p.ast, p.tokens
	tok := p.consume()

	if tok.Kind == OpNegation {
		atom := p.parseAtom()
		return ast.Append(NodeExprUnary(tok.Pos, tok.Kind, atom))
	}

	var atom NodeID
	switch tok.Kind {
	case LiteralNumber:
		return ast.Append(NodeLiteralNumber(tok.Pos, tok.Num))
	case LiteralString:
		return ast.Append(NodeLiteralString(tok.Pos, tok.Str))
	case LiteralTrue:
		return ast.Append(NodeLiteralBoolean(tok.Pos, true))
	case LiteralFalse:
		return ast.Append(NodeLiteralBoolean(tok.Pos, false))
	case Identifier:
		atom = ast.Append(NodeIdentifier(tok.Pos, tok.Str))
		// may be called as a function, so flows beyond switch block
	case IdentifierEmpty:
		return ast.Append(NodeIdentifierEmpty(tok.Pos))
	case ParenLeft:
		// grouped expression or function literal
		exprs := make([]NodeID, 0)
		for p.peek().Kind != ParenRight {
			expr := p.parseExpression()
			exprs = append(exprs, expr)
		}
		p.consumeKind(ParenRight)

		if tokens[p.idx].Kind == FunctionArrow {
			atom = p.parseFunctionLiteral(exprs)
			// parseAtom should not consume trailing Separators, but
			// 	parseFunctionLiteral does because it ends with expressions.
			// 	so we backtrack one token.
			p.idx--
		} else {
			atom = ast.Append(NodeExprList(tok.Pos, exprs))
		}
		// may be called as a function, so flows beyond switch block
	case CurlyParenLeft:
		return p.parseCompositeLiteral(tok.Pos)
	case SquareParenLeft:
		return p.parseListLiteral(tok.Pos)
	default:
		panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected start of atom, found %s", tok), tok.Pos}})
	}

	// bounds check here because parseExpression may have
	// consumed all tokens before this
	for p.idx < len(tokens) && tokens[p.idx].Kind == ParenLeft {
		atom = p.parseFunctionCall(atom)
	}

	return atom
}

func (p *parser) parseCompositeLiteral(pos Pos) NodeID {
	entries := make([]NodeCompositeKeyValue, 0)
	for p.peek().Kind != CurlyParenRight {
		keyExpr := p.parseExpression()

		var valExpr NodeID
		if p.tokens[p.idx].Kind == KeyValueSeparator { // "key: value" pair
			p.idx++
			valExpr = p.parseExpression() // Separator consumed by parseExpression
		} else if p.ast.Nodes[keyExpr].Kind == NodeKindIdentifier { // "key", shorthand for "key: key"
			valExpr = keyExpr
		} else {
			panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %s after composite key, found %s", KeyValueSeparator.String(), p.tokens[p.idx]), pos}})
		}

		entries = append(entries, NodeCompositeKeyValue{
			Key: keyExpr,
			Val: valExpr,
			Pos: p.ast.Nodes[keyExpr].Position(p.ast),
		})
	}
	_ = p.consumeKind(CurlyParenRight)

	return p.ast.Append(NodeLiteralComposite(pos, entries))
}

func (p *parser) parseListLiteral(pos Pos) NodeID {
	vals := make([]NodeID, 0)
	for p.peek().Kind != SquareParenRight {
		expr := p.parseExpression()
		vals = append(vals, expr)
	}
	p.consumeKind(SquareParenRight)
	return p.ast.Append(NodeLiteralList(pos, vals...))
}

// parses everything that follows MatchColon does not consume dangling separator -- that's for parseExpression
func (p *parser) parseMatchBody() []NodeMatchClause {
	_ = p.consumeKind(CurlyParenLeft)
	clauses := make([]NodeMatchClause, 0)
	for p.peek().Kind != CurlyParenRight {
		atom := p.parseExpression()
		_ = p.consumeKind(CaseArrow)
		expr := p.parseExpression()

		clauses = append(clauses, NodeMatchClause{
			Target:     atom,
			Expression: expr,
		})
	}
	p.consumeKind(CurlyParenRight)
	return clauses
}

func (p *parser) parseFunctionLiteral(arguments []NodeID) NodeID {
	// arguments := make([]NodeID, 0)
	// switch tok.Kind {
	// case ParenLeft:
	// LOOP:
	// 	for {
	// 		tk := tokens[p.idx]
	// 		switch tk.Kind {
	// 		case Identifier:
	// 			arguments = append(arguments, ast.Append(NodeIdentifier(tk.Pos, tk.Str)))
	// 		case IdentifierEmpty:
	// 			arguments = append(arguments, ast.Append(NodeIdentifierEmpty(tk.Pos)))
	// 		default:
	// 			break LOOP
	// 		}
	// 		p.idx++

	// 		_ = p.consumeKind(Separator)
	// 	}

	// 	_ = p.consumeKind(ParenRight)

	// case Identifier:
	// 	arguments = append(arguments, ast.Append(NodeIdentifier(tok.Pos, tok.Str)))
	// case IdentifierEmpty:
	// 	arguments = append(arguments, ast.Append(NodeIdentifierEmpty(tok.Pos)))
	// default:
	// 	panic(errParse{&Err{nil, ErrSyntax, fmt.Sprintf("malformed arguments list in function at %s", tok), tok.Pos}})
	// }

	pos := p.consumeKind(FunctionArrow).Pos

	body := p.parseExpression()

	return p.ast.Append(NodeLiteralFunction(pos, arguments, body))
}

func (p *parser) parseFunctionCall(function NodeID) NodeID {
	p.idx++

	arguments := make([]NodeID, 0)
	for p.peek().Kind != ParenRight {
		expr := p.parseExpression()
		arguments = append(arguments, expr)
	}
	p.consumeKind(ParenRight)

	return p.ast.Append(NodeFunctionCall(function, arguments))
}

// parse concurrently transforms a stream of Tok (tokens) to Node (AST nodes).
// This implementation uses recursive descent parsing.
func parse(ast *AST, tokens []Token) (nodes []NodeID, e errParse) {
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

		p := &parser{ast, tokens, i}
		expr := p.parseExpression()
		i = p.idx

		LogNode(ast.Nodes[expr])
		nodes = append(nodes, expr)
	}
	return
}
