package internal

import (
	"fmt"
	"strconv"
)

type Unit struct{}

// TODO: pass and use pos
type Parser[T any] func([]byte) ([]byte, T, errParse)

var parseByteAny Parser[byte] = func(in []byte) ([]byte, byte, errParse) {
	if len(in) < 1 {
		return nil, 0, errParse{&Err{nil, ErrSyntax, "unexpected EOF", Pos{}}}
	}
	return in[1:], in[0], errParse{}
}

func parseBytes(s string) Parser[string] {
	return func(in []byte) ([]byte, string, errParse) {
		if len(in) < len(s) || string(in[:len(s)]) != s {
			var msg string
			if len(in) == 0 {
				msg = fmt.Sprintf("unexpected EOF, expected %q", s)
			} else {
				msg = fmt.Sprintf("unexpected bytes %q, expected %q", in[:min(len(in), len(s))], s)
			}
			return nil, "", errParse{&Err{nil, ErrSyntax, msg, Pos{}}}
		}
		return in[len(s):], s, errParse{}
	}
}

func parseOr[T any](parsers ...Parser[T]) Parser[T] {
	return func(in []byte) ([]byte, T, errParse) {
		errs := ""
		for _, p := range parsers {
			out, v, err := p(in)
			if err.err == nil {
				return out, v, err
			}
			errs += err.err.Error() + ", "
		}
		return nil, *new(T), errParse{&Err{nil, ErrSyntax, "none matched: " + errs, Pos{}}}
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
	return func(in []byte) ([]byte, R, errParse) {
		ain, a, err := p1(in)
		if err.err != nil {
			return nil, *new(R), err
		}
		bin, b, err := p2(ain)
		if err.err != nil {
			return nil, *new(R), err
		}
		r, err := f(a, b)
		return bin, r, err
	}
}

func parseAnd3[A, B, C, R any](
	p1 Parser[A],
	p2 Parser[B],
	p3 Parser[C],
	f func(A, B, C) (R, errParse),
) Parser[R] {
	return func(in []byte) ([]byte, R, errParse) {
		ain, a, err := p1(in)
		if err.err != nil {
			return nil, *new(R), err
		}
		bin, b, err := p2(ain)
		if err.err != nil {
			return nil, *new(R), err
		}
		cin, c, err := p3(bin)
		if err.err != nil {
			return nil, *new(R), err
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
	return func(in []byte) ([]byte, R, errParse) {
		ain, a, err := p1(in)
		if err.err != nil {
			return nil, *new(R), err
		}
		bin, b, err := p2(ain)
		if err.err != nil {
			return nil, *new(R), err
		}
		cin, c, err := p3(bin)
		if err.err != nil {
			return nil, *new(R), err
		}
		din, d, err := p4(cin)
		if err.err != nil {
			return nil, *new(R), err
		}
		r, err := f(a, b, c, d)
		return din, r, err
	}
}

type Option[T any] struct {
	Value T
	Valid bool
}

func parseOptional[T any](p Parser[T]) Parser[Option[T]] {
	return func(in []byte) ([]byte, Option[T], errParse) {
		out, v, err := p(in)
		if err.err != nil {
			return in, Option[T]{}, errParse{}
		}
		return out, Option[T]{v, true}, errParse{}
	}
}

func parseMany[T any](p Parser[T]) Parser[[]T] {
	return func(in []byte) ([]byte, []T, errParse) {
		var res []T
		for {
			out, v, err := p(in)
			if err.err != nil {
				return in, res, errParse{}
			}
			res = append(res, v)
			in = out
		}
	}
}

func parseMap[T, R any](p Parser[T], f func(T) (R, errParse)) Parser[R] {
	return func(in []byte) ([]byte, R, errParse) {
		out, v, err := p(in)
		if err.err != nil {
			return in, *new(R), errParse{}
		}

		vv, err := f(v)
		return out, vv, err
	}
}

func parseByte(b byte) Parser[byte] {
	return parseMap(parseByteAny, func(v byte) (byte, errParse) {
		if v != b {
			return 0, errParse{&Err{nil, ErrSyntax, fmt.Sprintf("unexpected byte %c, expected %c", v, b), Pos{}}}
		}
		return v, errParse{}
	})
}

func parseIgnore[T any](p Parser[T]) Parser[struct{}] {
	return func(in []byte) ([]byte, struct{}, errParse) {
		out, _, err := p(in)
		return out, struct{}{}, err
	}
}

/*
123 | 123. | 123.456 | .456 | -123.456
*/
var parseNumber = parseAnd4(
	parseOptional(parseIgnore(parseMinus)),
	parseMany(parseDigit),
	parseOptional(parseDot),
	parseMany(parseDigit),
	func(
		isNegative Option[struct{}],
		integerPart []byte,
		dot Option[byte],
		fractionalPart []byte,
	) (NodeLiteralNumber, errParse) {
		if len(integerPart) == 0 && len(fractionalPart) == 0 || len(integerPart) != 0 && dot.Valid && len(fractionalPart) != 0 {
			return NodeLiteralNumber{}, errParse{&Err{nil, ErrSyntax, "invalid number", Pos{}}}
		}

		s := string(integerPart)
		if dot.Valid {
			s += "." + string(dot.Value)
		}
		if isNegative.Valid {
			s = "-" + s
		}

		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return NodeLiteralNumber{}, errParse{&Err{nil, ErrSyntax, err.Error(), Pos{}}}
		}

		return NodeLiteralNumber{
			val: f,
			Pos: Pos{},
		}, errParse{}
	})

var _charsEscaped = map[byte]byte{
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'\\': '\\',
	'\'': '\'',
}

// ” | 'abc' | 'abc\'def\n\\\r\tghi'
var parseString = parseAnd3(
	parseQuote,
	parseMany(parseMap(parseOr(
		parseAnd2(
			parseByte('\\'),
			parseByteAny,
			func(_ byte, b byte) (byte, errParse) {
				if _, ok := _charsEscaped[b]; !ok {
					return 0, errParse{
						&Err{nil, ErrSyntax, fmt.Sprintf("invalid escape sequence \\%c", b), Pos{}},
					}
				}
				return _charsEscaped[b], errParse{}
			},
		),
		parseByteAny,
	), func(b byte) (byte, errParse) {
		return b, errParse{}
	},
	)),
	parseQuote,
	func(
		_ byte,
		inside []byte,
		_ byte,
	) (NodeLiteralString, errParse) {
		return NodeLiteralString{
			val: string(inside),
			Pos: Pos{},
		}, errParse{}
	})

var parseBoolean = parseMap(parseOr(
	parseBytes("true"),
	parseBytes("false"),
), func(s string) (NodeLiteralBoolean, errParse) {
	switch s {
	case "true":
		return NodeLiteralBoolean{Pos{}, true}, errParse{}
	case "false":
		return NodeLiteralBoolean{Pos{}, false}, errParse{}
	default:
		return NodeLiteralBoolean{}, errParse{&Err{nil, ErrSyntax, "invalid boolean literal: " + s, Pos{}}}
	}
})

var parseIdentifierEmpty = parseMap(
	parseByte('_'),
	func(byte) (NodeIdentifierEmpty, errParse) {
		return NodeIdentifierEmpty{Pos{}}, errParse{}
	})

var parseIdentifier = parseMap(
	parseMany(parseMap(parseByteAny, func(b byte) (byte, errParse) {
		// TODO: other symbols
		if '0' <= b && b <= '9' || 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_' {
			return b, errParse{}
		}
		return 0, errParse{
			&Err{nil, ErrSyntax, fmt.Sprintf("invalid identifier character %c", b), Pos{}},
		}
	})),
	func(ident []byte) (NodeIdentifier, errParse) {
		switch {
		case len(ident) == 0:
			return NodeIdentifier{}, errParse{&Err{nil, ErrSyntax, "empty identifier", Pos{}}}
		case '0' <= ident[0] && ident[0] <= '9':
			return NodeIdentifier{}, errParse{&Err{nil, ErrSyntax, "identifier cannot start with digit", Pos{}}}
		default:
			return NodeIdentifier{Pos{}, string(ident)}, errParse{}
		}
	})

func parseNode[N Node](p Parser[N]) Parser[Node] {
	return parseMap(p, func(n N) (Node, errParse) { return n, errParse{} })
}

func parseExpression2(b []byte) ([]byte, Node, errParse) {
	parseList := parseAnd3(
		parseParenLeft,
		parseMany(parseExpression2), // TODO: separated by commas?
		parseParenRight,
		func(_ byte, exprs []Node, _ byte) (NodeExprList, errParse) {
			return NodeExprList{Pos{}, exprs}, errParse{}
		},
	)

	return parseOr(
		// literals
		parseNode(parseNumber),
		parseNode(parseString),
		parseNode(parseBoolean),
		// identifier
		parseNode(parseIdentifier),
		// assignment
		parseAnd3(
			parseExpression2,
			parseDefine,
			parseExpression2,
			func(lvalue Node, _ string, rvalue Node) (Node, errParse) {
				return NodeExprBinary{Pos{}, OpDefine, lvalue, rvalue}, errParse{}
			},
		),
		// block
		parseAnd3(
			parseParenLeft,
			parseMany(parseExpression2), // TODO: separated by commas?
			parseParenRight,
			func(_ byte, exprs []Node, _ byte) (Node, errParse) {
				if len(exprs) == 0 {
					return NodeIdentifierEmpty{Pos{}}, errParse{} // TODO: () is "nil"
				}
				res := exprs[0]
				for i := 1; i < len(exprs); i++ {
					// TODO: sequencing operator
					res = NodeExprBinary{Pos{}, OpSubtract, res, exprs[i]}
				}
				return res, errParse{}
			},
		),
		// list
		parseNode(parseList),
		// dict
		parseAnd3(
			parseBraceLeft,
			parseMany(parseAnd3(
				parseExpression2,
				parseColon,
				parseExpression2,
				func(k Node, _ byte, v Node) (NodeObjectEntry, errParse) {
					return NodeObjectEntry{Pos{}, k, v}, errParse{}
				},
			)), // TODO: separated by commas?
			parseBraceRight,
			func(_ byte, kvs []NodeObjectEntry, _ byte) (Node, errParse) {
				return NodeLiteralObject{Pos{}, kvs}, errParse{}
			},
		),
		// lambda
		parseNode(parseAnd3(
			parseOr(
				parseMap(parseIdentifier, func(ident NodeIdentifier) ([]Node, errParse) { return []Node{ident}, errParse{} }),
				parseMap(parseList, func(list NodeExprList) ([]Node, errParse) { return list.expressions, errParse{} }),
			),
			parseFunctionArrow,
			parseExpression2,
			func(args []Node, _ string, body Node) (NodeLiteralFunction, errParse) {
				return NodeLiteralFunction{Pos{}, args, body}, errParse{}
			},
		)),
		// function call
		parseAnd3(
			parseExpression2,
			parseParenLeft,
			parseMany(parseExpression2), // TODO: separated by commas?
			func(fn Node, _ byte, args []Node) (Node, errParse) {
				return NodeFunctionCall{fn, args}, errParse{}
			},
		),
		// unary expression
		parseAnd2(
			parseByte('-'),
			parseExpression2,
			func(_ byte, n Node) (Node, errParse) { return NodeExprUnary{Pos{}, OpNegation, n}, errParse{} },
		),
		// binary expression
		parseAnd3(
			parseExpression2,
			parseOr(
				// NOTE: ordered by priority, higher is higher priority
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
			parseExpression2,
			func(m Node, op byte, n Node) (Node, errParse) {
				mp := map[byte]Kind{
					'+': OpAdd,
					'-': OpSubtract,
					'*': OpMultiply,
					'/': OpDivide,
					'%': OpModulus,
					'&': OpLogicalAnd,
					'|': OpLogicalOr,
					'^': OpLogicalXor,
					'>': OpGreaterThan,
					'<': OpLessThan,
					'=': OpEqual,
				}
				opKind, ok := mp[op]
				if !ok {
					return nil, errParse{
						&Err{nil, ErrSyntax, fmt.Sprintf("invalid operator %c", op), Pos{}},
					}
				}
				return NodeExprBinary{Pos{}, opKind, n, m}, errParse{}
			},
		),
	)(b)
}
