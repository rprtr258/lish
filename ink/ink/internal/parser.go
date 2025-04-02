package internal

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/rprtr258/fun"
)

type Unit = struct{}

var unit = Unit{}

type ctx struct {
	ast      *AST
	b        []byte
	offset   int
	filename string
}

func newCtx(b []byte, filename string) *ctx {
	return &ctx{NewAstSlice(), b, 0, filename}
}

// TODO: remove?
func (c *ctx) to(offset int) *ctx {
	return &ctx{c.ast, c.b, offset, c.filename}
}

func (c *ctx) skipSpaces(skipNewlines bool) *ctx {
	return c.to(c.offset + skipSpaces(c.bb(), skipNewlines))
}

func (c ctx) bb() []byte {
	return c.b[c.offset:]
}

func (c ctx) len() int {
	return len(c.bb())
}

func (c ctx) at(i int) byte {
	return c.bb()[i]
}

func (c ctx) pos() Pos {
	line := bytes.Count(c.b[:c.offset], []byte{'\n'})
	precol := bytes.LastIndexByte(c.b[:c.offset], '\n')
	return Pos{c.filename, line + 1, c.offset - precol - 1} // TODO: move +1s to Pos.String ?
}

// TODO: pass and use pos
type Parser[T any] func(*ctx) (int, T, errParse)

var parseByteAny Parser[byte] = func(c *ctx) (int, byte, errParse) {
	if c.len() < 1 {
		return 0, 0, errParse{&Err{nil, ErrSyntax, "unexpected EOF", c.pos()}}
	}
	return c.offset + 1, c.b[c.offset], errParse{}
}

func parseBytes(s string) Parser[string] {
	return func(c *ctx) (int, string, errParse) {
		if c.len() < len(s) || string(c.bb()[:len(s)]) != s {
			var msg string
			if c.len() == 0 {
				msg = fmt.Sprintf("unexpected EOF, expected %q", s)
			} else {
				msg = fmt.Sprintf("unexpected bytes %q, expected %q", c.bb()[:min(len(c.bb()), len(s))], s)
			}
			return 0, "", errParse{&Err{nil, ErrSyntax, msg, c.pos()}}
		}
		return c.offset + len(s), s, errParse{}
	}
}

func parseOr[T any](parsers ...Parser[T]) Parser[T] {
	return func(c *ctx) (int, T, errParse) {
		errs := ""
		for _, p := range parsers {
			out, v, err := p(c)
			if err.Err == nil {
				return out, v, err
			}
			errs += err.Err.Error() + ", "
		}
		return 0, *new(T), errParse{&Err{nil, ErrSyntax, "none matched: " + errs, c.pos()}}
	}
}

var (
	parseDigit = parseOr(
		parseByte('0'),
		parseByte('1'),
		parseByte('2'),
		parseByte('3'),
		parseByte('4'),
		parseByte('5'),
		parseByte('6'),
		parseByte('7'),
		parseByte('8'),
		parseByte('9'),
	)
	parseQuote         = parseByte('\'')
	parseDot           = parseByte('.')
	parseEqual         = parseByte('=')
	parseFunctionArrow = parseBytes("=>")
	parseComma         = parseByte(',')
	parseColon         = parseByte(':')
	parseDefine        = parseBytes(":=")
	parseMatch         = parseBytes("::")
	parseArrow         = parseBytes("->")
	parseMinus         = parseByte('-')
	parseNegation      = parseByte('~')
	parseAdd           = parseByte('+')
	parseMultiply      = parseByte('*')
	parseDivide        = parseByte('/')
	parseModulus       = parseByte('%')
	parseGreaterThan   = parseByte('>')
	parseLessThan      = parseByte('<')
	parseLogicalAnd    = parseByte('&')
	parseLogicalOr     = parseByte('|')
	parseLogicalXor    = parseByte('^')

	parseParenLeft, parseParenRight     = parseByte('('), parseByte(')')
	parseBracketLeft, parseBracketRight = parseByte('['), parseByte(']')
	parseBraceLeft, parseBraceRight     = parseByte('{'), parseByte('}')
)

func parseAnd2[A, B, R any](
	p1 Parser[A],
	p2 Parser[B],
	f func(A, B) (R, errParse),
) Parser[R] {
	return func(c *ctx) (int, R, errParse) {
		cc := c.to(c.offset + skipSpaces(c.bb(), true))
		aoffset, a, err := p1(cc)
		if err.Err != nil {
			return 0, *new(R), err
		}
		cc = cc.to(aoffset).skipSpaces(true) // TODO: copy in or, not and
		boffset, b2, err := p2(cc)
		if err.Err != nil {
			return 0, *new(R), err
		}
		r, err := f(a, b2)
		return boffset, r, err
	}
}

func parseAnd3_[A, B, C, R any](
	p1 Parser[A],
	p2 Parser[B],
	p3 Parser[C],
	f func(A, B, C) (R, errParse),
) Parser[R] {
	return func(cc *ctx) (int, R, errParse) {
		ain, a, err := p1(cc)
		if err.Err != nil {
			return 0, *new(R), err
		}
		bin, b, err := p2(cc.to(ain))
		if err.Err != nil {
			return 0, *new(R), err
		}
		cin, c, err := p3(cc.to(bin))
		if err.Err != nil {
			return 0, *new(R), err
		}
		r, err := f(a, b, c)
		return cin, r, err
	}
}

func parseAnd3[A, B, C, R any](
	p1 Parser[A],
	p2 Parser[B],
	p3 Parser[C],
	f func(A, B, C) (R, errParse),
) Parser[R] {
	return func(cc *ctx) (int, R, errParse) {
		cc = cc.skipSpaces(true)
		ain, a, err := p1(cc)
		if err.Err != nil {
			return 0, *new(R), err
		}
		cc = cc.to(ain).skipSpaces(true)
		bin, b, err := p2(cc)
		if err.Err != nil {
			return 0, *new(R), err
		}
		cc = cc.to(bin).skipSpaces(true)
		cin, c, err := p3(cc)
		if err.Err != nil {
			return 0, *new(R), err
		}
		r, err := f(a, b, c)
		return cin, r, err
	}
}

func parseAnd4[A, B, C, D, R any](
	p1 Parser[A],
	p2 Parser[B],
	p3 Parser[C],
	p4 Parser[D],
	f func(A, B, C, D) (R, errParse),
) Parser[R] {
	return func(cc *ctx) (int, R, errParse) {
		cc = cc.skipSpaces(true) // TODO: copy in or, not and
		ain, a, err := p1(cc)
		if err.Err != nil {
			return 0, *new(R), err
		}
		cc = cc.to(ain).skipSpaces(true) // TODO: copy in or, not and
		bin, b, err := p2(cc)
		if err.Err != nil {
			return 0, *new(R), err
		}
		cc = cc.to(bin).skipSpaces(true) // TODO: copy in or, not and
		cin, c, err := p3(cc)
		if err.Err != nil {
			return 0, *new(R), err
		}
		cc = cc.to(cin).skipSpaces(true) // TODO: copy in or, not and
		din, d, err := p4(cc)
		if err.Err != nil {
			return 0, *new(R), err
		}
		r, err := f(a, b, c, d)
		return din, r, err
	}
}

func parseOptional[T any](p Parser[T]) Parser[fun.Option[T]] {
	return func(c *ctx) (int, fun.Option[T], errParse) {
		out, v, err := p(c.to(c.offset))
		if err.Err != nil {
			return c.offset, fun.Option[T]{}, errParse{}
		}
		return out, fun.Option[T]{v, true}, errParse{}
	}
}

// a,b,c, -> [a,b,c]
// a,b,c  -> [a,b,c]
func parseMany[T any](p Parser[T]) Parser[[]T] {
	return func(c *ctx) (int, []T, errParse) {
		var res []T
		b := c.offset
		for {
			b2, v, err := p(c.to(b))
			if err.Err != nil {
				return b, res, errParse{}
			}
			res = append(res, v)
			if b2 == len(c.b) {
				return b2, res, errParse{}
			}
			{ // TODO: consider ';' in block instead of ','
				b2 += skipSpaces(c.b[b2:], false)
				delim := parseOr(parseComma, parseByte('\n'))
				b3, _, err := delim(c.to(b2))
				if err.Err != nil {
					return b2, res, errParse{}
				}
				b2 = b3 + skipSpaces(c.b[b3:], true)
			}
			b = b2
		}
	}
}

func parseMap[T, R any](p Parser[T], f func(T) (R, errParse)) Parser[R] {
	return func(c *ctx) (int, R, errParse) {
		out, v, err := p(c)
		if err.Err != nil {
			return 0, *new(R), err
		}

		vv, err := f(v)
		return out, vv, err
	}
}

func parseByte(c byte) Parser[byte] {
	return func(cc *ctx) (int, byte, errParse) {
		return parseMap(parseByteAny, func(v byte) (byte, errParse) {
			if v != c {
				return 0, errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %c, found %c", c, v), cc.pos()}}
			}
			return v, errParse{}
		})(cc)
	}
}

func parseIgnore[T any](p Parser[T]) Parser[struct{}] {
	return func(c *ctx) (int, struct{}, errParse) {
		out, _, err := p(c)
		return out, struct{}{}, err
	}
}

// 123 | 456
// NOTE: RETURNS INT, NOT NODE INDEX
func parseInt(c *ctx) (int, int, errParse) {
	if c.len() == 0 || c.at(0) < '0' || c.at(0) > '9' {
		return 0, 0, errParse{&Err{nil, ErrSyntax, "EOF", c.pos()}}
	}

	res := 0
	offset := c.offset
	for offset < len(c.b) && '0' <= c.b[offset] && c.b[offset] <= '9' {
		d := int(c.b[offset] - '0')
		res = res*10 + d
		offset++
	}
	return offset, res, errParse{}
}

// 123 | 123. | 123.456 | .456 | ~123.456
func parseNumber(c *ctx) (int, int, errParse) {
	return parseAnd3_(
		parseOptional(parseIgnore(parseNegation)),
		parseInt,
		parseOptional(parseAnd2(
			parseIgnore(parseDot),
			parseInt,
			func(
				dot struct{},
				fractionalPart int,
			) (int, errParse) {
				return fractionalPart, errParse{}
			},
		)),
		func(
			isNegative fun.Option[struct{}],
			integerPart int,
			fractionalPart fun.Option[int],
		) (int, errParse) {
			s := strconv.Itoa(integerPart)
			if fractionalPart.Valid {
				s += "." + strconv.Itoa(fractionalPart.Value)
			}
			if isNegative.Valid {
				s = "-" + s
			}

			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return -1, errParse{&Err{nil, ErrSyntax, err.Error(), c.pos()}}
			}

			return c.ast.Append(NodeLiteralNumber{
				Val: f,
				Pos: c.pos(),
			}), errParse{}
		})(c)
}

var _charsEscaped = map[byte]byte{
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'\\': '\\',
	'\'': '\'',
}

func parseStringIn(c *ctx, quote byte) (int, int, errParse) {
	end := bytes.IndexByte(c.bb()[1:], quote)
	if end == -1 {
		return 0, -1, errParse{&Err{nil, ErrSyntax, "expected string, but found EOF", c.pos()}}
	}

	return c.offset + 1 + end + 1, c.ast.Append(NodeLiteralString{
		Val: string(c.bb()[1 : end+1]),
		Pos: c.pos(),
	}), errParse{}
}

// TODO: PIDARAS NA VSCODE ZAMENYAET MNE two ' to ”
// ” | 'abc' | 'abc\'def\n\\\r\tghi' | `aaaa`
func parseString(c *ctx) (int, int, errParse) {
	switch {
	case c.len() == 0:
		return 0, -1, errParse{&Err{nil, ErrSyntax, "string: unexpected EOF", c.pos()}}
	case c.at(0) == '\'':
		return parseStringIn(c, '\'')
	case c.at(0) == '`':
		return parseStringIn(c, '`')
	default:
		return 0, -1, errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected string, found %c", c.at(0)), c.pos()}}
	}
}

func parseBoolean(c *ctx) (int, int, errParse) {
	return parseMap(parseOr(
		parseBytes("true"),
		parseBytes("false"),
	), func(s string) (int, errParse) {
		switch s {
		case "true":
			return c.ast.Append(NodeLiteralBoolean{c.pos(), true}), errParse{}
		case "false":
			return c.ast.Append(NodeLiteralBoolean{c.pos(), false}), errParse{}
		default:
			return -1, errParse{&Err{nil, ErrSyntax, "invalid boolean literal: " + s, c.pos()}}
		}
	})(c)
}

func isValidIdentifierByte(c byte) bool {
	return '0' <= c && c <= '9' ||
		'a' <= c && c <= 'z' ||
		'A' <= c && c <= 'Z' ||
		c == '_' || c == '?'
}

func parseIdentifier(c *ctx) (int, int, errParse) {
	i := c.offset
	for i < len(c.b) && isValidIdentifierByte(c.b[i]) {
		i++
	}

	ident := c.b[c.offset:i]
	switch {
	case len(ident) == 0:
		return 0, -1, errParse{&Err{nil, ErrSyntax, "empty identifier", c.pos()}}
	case '0' <= ident[0] && ident[0] <= '9':
		return 0, -1, errParse{&Err{nil, ErrSyntax, "identifier cannot start with digit", c.pos()}}
	case string(ident) == "_":
		return i, c.ast.Append(NodeIdentifierEmpty{c.pos()}), errParse{}
	default:
		return i, c.ast.Append(NodeIdentifier{c.pos(), string(ident)}), errParse{}
	}
}

// TODO: replace with comments skip
func skipSpaces(b []byte, skipNewlines bool) int {
	i := 0
	for {
		for i < len(b) && (bytes.Contains([]byte(" \t\r"), b[i:][:1]) || skipNewlines && b[i] == '\n') {
			i++
		}

		if i == len(b) || b[i] != '#' {
			return i
		}

		for i < len(b) && b[i] != '\n' {
			i++
		}
		if !skipNewlines {
			return i
		}
	}
}

func parseLhs(c *ctx) (int, int, errParse) {
	i, res, err := parseOr(
		parseIdentifier,
		parseDict,
		parseList,
	)(c)
	if err.Err != nil {
		return 0, -1, errParse{&Err{err.Err, ErrSyntax, "assignment start", c.pos()}}
	}

	if _, ok := c.ast.Nodes[res].(NodeIdentifierEmpty); ok {
		return i, res, errParse{}
	}

	for i < len(c.b) && c.b[i] == '.' {
		b2, rhs, err := parseAnd2(
			parseIgnore(parseDot),
			parseOr(
				parseLiteral,
				parseMap(
					parseIdentifier,
					func(k int) (int, errParse) {
						if _, ok := c.ast.Nodes[k].(NodeIdentifierEmpty); ok {
							return k, errParse{}
						}
						return c.ast.Append(NodeLiteralString{c.pos(), c.ast.Nodes[k].(NodeIdentifier).Val}), errParse{}
					},
				),
				parseMap(
					parseBlock,
					func(b int) (int, errParse) {
						if exprs := c.ast.Nodes[b].(NodeExprList).Expressions; len(exprs) == 1 {
							return exprs[0], errParse{}
						}
						return -1, errParse{&Err{nil, ErrSyntax, "too many indices", c.pos()}}
					},
				),
			),
			func(_ Unit, k int) (int, errParse) {
				return k, errParse{}
			},
		)(c.to(i))
		if err.Err != nil {
			return i, res, errParse{}
		}
		i = b2
		res = c.ast.Append(NodeExprBinary{c.pos(), OpAccessor, res, rhs})
	}

	return i, res, errParse{}
}

func parseAssignment(c *ctx) (int, int, errParse) {
	return parseAnd3(
		parseLhs,
		parseDefine,
		parseExpression,
		func(lvalue int, _ string, rvalue int) (int, errParse) {
			return c.ast.Append(NodeExprBinary{c.pos(), OpDefine, lvalue, rvalue}), errParse{}
		},
	)(c)
}

func parseList(c *ctx) (int, int, errParse) {
	return parseAnd3(
		parseBracketLeft,
		parseMany(parseExpression),
		parseBracketRight,
		func(_ byte, exprs []int, _ byte) (int, errParse) {
			return c.ast.Append(NodeLiteralList{c.pos(), exprs}), errParse{}
		},
	)(c)
}

func parseBlock(c *ctx) (int, int, errParse) {
	return parseAnd3(
		parseParenLeft,
		parseMany(parseExpression),
		parseParenRight,
		func(_ byte, exprs []int, _ byte) (int, errParse) {
			exprs = fun.Filter(func(n int) bool {
				_, ok := c.ast.Nodes[n].(NodeIdentifierEmpty)
				return !ok
			}, exprs...)
			// TODO: optimization, now can't be commented back because used to parse lambda args
			// if len(exprs) == 1 {
			// 	return exprs[0], errParse{}
			// }
			// TODO: if len(exprs) == 0 return Nil/Unit
			return c.ast.Append(NodeExprList{c.pos(), exprs}), errParse{}
		},
	)(c)
}

func parseLiteral(c *ctx) (int, int, errParse) {
	return parseOr(
		parseNumber,
		parseString,
		parseBoolean,
	)(c)
}

func parseDict(c *ctx) (int, int, errParse) {
	return parseAnd3(
		parseBraceLeft,
		parseMany(parseOr(
			parseAnd3( // "key: value" pair
				parseExpression,
				parseColon,
				parseExpression,
				func(k int, _ byte, v int) (NodeObjectEntry, errParse) {
					return NodeObjectEntry{c.pos(), k, v}, errParse{}
				},
			),
			parseMap( // "key" meaning "key: key"
				parseIdentifier,
				func(k int) (NodeObjectEntry, errParse) {
					return NodeObjectEntry{c.pos(), k, k}, errParse{}
				},
			),
		)),
		parseBraceRight,
		func(_ byte, kvs []NodeObjectEntry, _ byte) (int, errParse) {
			return c.ast.Append(NodeLiteralObject{c.pos(), kvs}), errParse{}
		},
	)(c)
}

func parseUnary(c *ctx) (int, int, errParse) {
	return parseAnd2(
		parseByte('-'),
		parseExpression,
		func(_ byte, n int) (int, errParse) {
			return c.ast.Append(NodeExprUnary{c.pos(), OpNegation, n}), errParse{}
		},
	)(c)
}

func parseExpression(c *ctx) (int, int, errParse) {
	cc := c.skipSpaces(true)
	b, lhs, err := parseOr(
		parseAssignment,
		parseBlock,
		parseList,
		parseDict,
		parseUnary,
		parseLiteral,
		parseIdentifier,
	)(cc)
	if err.Err != nil {
		return 0, -1, err
	}

	if _, ok := c.ast.Nodes[lhs].(NodeIdentifierEmpty); ok {
		return b, lhs, err
	}

	for {
		if b == len(cc.b) { // TODO: must be just ==
			return b, lhs, errParse{}
		}

		switch c.ast.Nodes[lhs].(type) {
		case NodeIdentifier, NodeExprList:
			b2, lambda, err := parseAnd2(
				parseFunctionArrow,
				parseExpression,
				func(_ string, body int) (int, errParse) {
					var args []int
					switch n := cc.ast.Nodes[lhs].(type) {
					case NodeIdentifier:
						args = []int{lhs}
					case NodeExprList:
						args = n.Expressions
					}
					return c.ast.Append(NodeLiteralFunction{c.pos(), args, body}), errParse{}
				},
			)(cc.to(b))
			if err.Err == nil {
				lhs = lambda
				b = b2
				continue
			}
		}

		if b < len(c.b) && cc.b[b] != '\n' {
			b2, call, err := parseAnd3(
				parseParenLeft,
				parseMany(parseExpression),
				parseParenRight,
				func(_ byte, args []int, _ byte) (int, errParse) {
					LogAST(cc.ast)
					return c.ast.Append(NodeFunctionCall{lhs, args}), errParse{}
				},
			)(c.to(b))
			if err.Err == nil {
				lhs = call
				b = b2
				continue
			}
		}

		{
			b2, bin, err := parseAnd2(
				parseOr(
					// NOTE: ordered by priority, higher is higher priority
					parseByte('.'),
					parseByte('%'),
					parseByte('*'),
					parseByte('/'),
					parseByte('>'),
					parseByte('<'),
					parseByte('='),
					parseByte('&'),
					parseByte('^'),
					parseByte('|'),
					parseByte('+'),
					parseByte('-'),
				),
				parseExpression,
				func(op byte, rhs int) (int, errParse) {
					mp := map[byte]Kind{
						'.': OpAccessor,
						'%': OpModulus,
						'*': OpMultiply,
						'/': OpDivide,
						'>': OpGreaterThan,
						'<': OpLessThan,
						'=': OpEqual,
						'&': OpLogicalAnd,
						'^': OpLogicalXor,
						'|': OpLogicalOr,
						'+': OpAdd,
						'-': OpSubtract,
					}
					opKind, ok := mp[op]
					if !ok {
						return -1, errParse{&Err{nil, ErrSyntax, fmt.Sprintf("invalid operator %c", op), c.pos()}}
					}
					return c.ast.Append(NodeExprBinary{c.pos(), opKind, lhs, rhs}), errParse{}
				},
			)(c.to(b))
			if err.Err == nil {
				lhs = bin
				b = b2
				continue
			}
		}

		{
			b2, match, err := parseAnd4(
				parseMatch,
				parseBraceLeft,
				parseMany(parseAnd3(
					parseExpression,
					parseArrow,
					parseExpression,
					func(target int, _ string, expression int) (int, errParse) {
						return c.ast.Append(NodeMatchClause{target, expression}), errParse{}
					},
				)),
				parseBraceRight,
				func(_ string, _ byte, clauses []int, _ byte) (int, errParse) {
					return c.ast.Append(NodeExprMatch{Condition: lhs, Clauses: clauses}), errParse{}
				},
			)(c.to(b))
			if err.Err == nil {
				lhs = match
				b = b2
				continue
			}
		}

		return b, lhs, errParse{}
	}
}
