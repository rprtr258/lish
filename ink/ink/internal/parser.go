package internal

import (
	"bytes"
	"fmt"
	"strconv"
)

type Unit struct{}

var unit = Unit{}

// TODO: pass and use pos
type Parser[T any] func(*AST, []byte) ([]byte, T, errParse)

var parseByteAny Parser[byte] = func(_ *AST, in []byte) ([]byte, byte, errParse) {
	if len(in) < 1 {
		return nil, 0, errParse{&Err{nil, ErrSyntax, "unexpected EOF", Pos{}}}
	}
	return in[1:], in[0], errParse{}
}

func parseBytes(s string) Parser[string] {
	return func(_ *AST, b []byte) ([]byte, string, errParse) {
		if len(b) < len(s) || string(b[:len(s)]) != s {
			var msg string
			if len(b) == 0 {
				msg = fmt.Sprintf("unexpected EOF, expected %q", s)
			} else {
				msg = fmt.Sprintf("unexpected bytes %q, expected %q", b[:min(len(b), len(s))], s)
			}
			return nil, "", errParse{&Err{nil, ErrSyntax, msg, Pos{}}}
		}
		return b[len(s):], s, errParse{}
	}
}

func parseOr[T any](parsers ...Parser[T]) Parser[T] {
	return func(ast *AST, b []byte) ([]byte, T, errParse) {
		LogToken2(">OR", "%d", len(parsers))
		errs := ""
		for i, p := range parsers {
			out, v, err := p(ast, b)
			LogToken2("OR", "%d/%d %q %+v", i, len(parsers), string(out), err.Err)
			if err.Err == nil {
				return out, v, err
			}
			errs += err.Err.Error() + ", "
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

func skipSpaces(b []byte) []byte {
	for len(b) > 0 && bytes.Contains([]byte(" \t\r\n"), b[:1]) {
		b = b[1:]
	}
	return b
}

func parseAnd2[A, B, R any](
	p1 Parser[A],
	p2 Parser[B],
	f func(A, B) (R, errParse),
) Parser[R] {
	return func(ast *AST, in []byte) ([]byte, R, errParse) {
		ain, a, err := p1(ast, skipSpaces(in))
		if err.Err != nil {
			return nil, *new(R), err
		}
		bin, b, err := p2(ast, skipSpaces(ain))
		if err.Err != nil {
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
	return func(ast *AST, in []byte) ([]byte, R, errParse) {
		LogToken2("IN", "%q", string(in))
		ain, a, err := p1(ast, skipSpaces(in))
		if err.Err != nil {
			return nil, *new(R), err
		}
		LogToken2("AIN", "%q", string(ain))
		bin, b, err := p2(ast, skipSpaces(ain))
		if err.Err != nil {
			return nil, *new(R), err
		}
		LogToken2("BIN", "%q", string(bin))
		cin, c, err := p3(ast, skipSpaces(bin))
		if err.Err != nil {
			return nil, *new(R), err
		}
		LogToken2("CIN", "%q", string(cin))
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
	return func(ast *AST, in []byte) ([]byte, R, errParse) {
		ain, a, err := p1(ast, skipSpaces(in))
		if err.Err != nil {
			return nil, *new(R), err
		}
		bin, b, err := p2(ast, skipSpaces(ain))
		if err.Err != nil {
			return nil, *new(R), err
		}
		cin, c, err := p3(ast, skipSpaces(bin))
		if err.Err != nil {
			return nil, *new(R), err
		}
		din, d, err := p4(ast, skipSpaces(cin))
		if err.Err != nil {
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
	return func(ast *AST, in []byte) ([]byte, Option[T], errParse) {
		out, v, err := p(ast, in)
		if err.Err != nil {
			return in, Option[T]{}, errParse{}
		}
		return out, Option[T]{v, true}, errParse{}
	}
}

func parseMany0[T any](p Parser[T]) Parser[[]T] {
	return func(ast *AST, in []byte) ([]byte, []T, errParse) {
		var res []T
		LogToken2(">MANY", "%q", string(in))
		for {
			out, v, err := p(ast, in)
			LogToken2("MANY", "%q %q", string(in), string(out))
			if err.Err != nil {
				return in, res, errParse{}
			}
			res = append(res, v)
			if len(out) == 0 {
				return out, res, errParse{}
			}
			in = out
		}
	}
}

func parseMap[T, R any](p Parser[T], f func(T) (R, errParse)) Parser[R] {
	return func(ast *AST, in []byte) ([]byte, R, errParse) {
		out, v, err := p(ast, in)
		if err.Err != nil {
			return in, *new(R), err
		}

		vv, err := f(v)
		return out, vv, err
	}
}

func parseByte(c byte) Parser[byte] {
	return parseMap(parseByteAny, func(v byte) (byte, errParse) {
		if v != c {
			return 0, errParse{&Err{nil, ErrSyntax, fmt.Sprintf("expected %c, found %c", c, v), Pos{}}}
		}
		return v, errParse{}
	})
}

func parseIgnore[T any](p Parser[T]) Parser[struct{}] {
	return func(ast *AST, in []byte) ([]byte, struct{}, errParse) {
		out, _, err := p(ast, in)
		return out, struct{}{}, err
	}
}

// 123 | 123. | 123.456 | .456 | -123.456
func parseNumber(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseAnd4(
		parseOptional(parseIgnore(parseMinus)),
		parseMany0(parseDigit),
		parseOptional(parseDot),
		parseMany0(parseDigit),
		func(
			isNegative Option[struct{}],
			integerPart []byte,
			dot Option[byte],
			fractionalPart []byte,
		) (int, errParse) {
			if len(integerPart) == 0 && len(fractionalPart) == 0 || len(integerPart) != 0 && dot.Valid && len(fractionalPart) != 0 {
				return -1, errParse{&Err{nil, ErrSyntax, "invalid number", Pos{}}}
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
				return -1, errParse{&Err{nil, ErrSyntax, err.Error(), Pos{}}}
			}

			return ast.Append(NodeLiteralNumber{
				Val: f,
				Pos: Pos{},
			}), errParse{}
		})(ast, b)
}

var _charsEscaped = map[byte]byte{
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'\\': '\\',
	'\'': '\'',
}

// TODO: PIDARAS NA VSCODE ZAMENYAET MNE two ' to ”
// ” | 'abc' | 'abc\'def\n\\\r\tghi'
func parseString(ast *AST, b []byte) ([]byte, int, errParse) {
	if len(b) == 0 || b[0] != '\'' {
		return nil, -1, errParse{&Err{nil, ErrSyntax, "expected string", Pos{}}}
	}

	end := bytes.IndexByte(b[1:], '\'')
	if end == -1 {
		return nil, -1, errParse{&Err{nil, ErrSyntax, "expected string, but found EOF", Pos{}}}
	}

	return b[end+2:], ast.Append(NodeLiteralString{
		Val: string(b[1 : end+1]),
		Pos: Pos{},
	}), errParse{}
}

func parseBoolean(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseMap(parseOr(
		parseBytes("true"),
		parseBytes("false"),
	), func(s string) (int, errParse) {
		switch s {
		case "true":
			return _astTrueLiteralIdx, errParse{}
		case "false":
			return _astFalseLiteralIdx, errParse{}
		default:
			return -1, errParse{&Err{nil, ErrSyntax, "invalid boolean literal: " + s, Pos{}}}
		}
	})(ast, b)
}

var parseIdentifierEmpty = parseMap(
	parseByte('_'),
	func(byte) (int, errParse) {
		return _astEmptyIdentifierIdx, errParse{}
	})

func parseIdentifier(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseMap(
		parseMany0(parseMap(parseByteAny, func(c byte) (byte, errParse) {
			// TODO: other symbols
			if '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_' {
				return c, errParse{}
			}
			return 0, errParse{&Err{nil, ErrSyntax, fmt.Sprintf("invalid identifier character %c", c), Pos{}}}
		})),
		func(ident []byte) (int, errParse) {
			switch {
			case len(ident) == 0:
				return -1, errParse{&Err{nil, ErrSyntax, "empty identifier", Pos{}}}
			case '0' <= ident[0] && ident[0] <= '9':
				return -1, errParse{&Err{nil, ErrSyntax, "identifier cannot start with digit", Pos{}}}
			default:
				return ast.Append(NodeIdentifier{Pos{}, string(ident)}), errParse{}
			}
		})(ast, b)
}

func parseComment(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseMap(func(_ *AST, b []byte) ([]byte, Unit, errParse) {
		LogToken2("COMMENT", "%q", string(b))
		if len(b) > 0 && !bytes.Contains([]byte(" \t\r\n"), b[:1]) {
			return b, unit, errParse{&Err{nil, ErrSyntax, "not a comment(1)", Pos{}}}
		}

		for {
			for len(b) > 0 && bytes.Contains([]byte(" \t\r\n"), b[:1]) {
				b = b[1:]
			}

			if len(b) == 0 {
				return b, unit, errParse{}
			}

			switch b[0] {
			case '#':
				if i := bytes.IndexByte(b, '\n'); i == -1 {
					// no newline, comment till end of file
					return []byte(nil), unit, errParse{}
				} else {
					return b[i+1:], unit, errParse{}
				}
			case '\n':
				// empty line, skip
				continue
			default:
				return b, unit, errParse{}
			}
		}
	}, func(_ Unit) (int, errParse) {
		return _astEmptyIdentifierIdx, errParse{} // NOTE: since we have to return some node, return empty
	})(ast, b)
}

func parseAssignment(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseAnd3(
		parseOr(
			parseIdentifier,
			parseDict,
			parseList,
		),
		parseDefine,
		parseExpression,
		func(lvalue int, _ string, rvalue int) (int, errParse) {
			return ast.Append(NodeExprBinary{Pos{}, OpDefine, lvalue, rvalue}), errParse{}
		},
	)(ast, b)
}

func parseList(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseAnd3(
		parseBracketLeft,
		parseMany0(parseExpression), // TODO: separated by commas?
		parseBracketRight,
		func(_ byte, exprs []int, _ byte) (int, errParse) {
			return ast.Append(NodeExprList{Pos{}, exprs}), errParse{}
		},
	)(ast, b)
}

func parseLambda(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseAnd3(
		parseOr(
			parseMap(parseIdentifier, func(ident int) ([]int, errParse) { return []int{ident}, errParse{} }),
			parseMap(parseList, func(list int) ([]int, errParse) { return ast.Nodes[list].(NodeExprList).Expressions, errParse{} }),
		),
		parseFunctionArrow,
		parseExpression,
		func(args []int, _ string, body int) (int, errParse) {
			return ast.Append(NodeLiteralFunction{Pos{}, args, body}), errParse{}
		},
	)(ast, b)
}

func parseBlock(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseAnd3(
		parseParenLeft,
		parseMany0(parseExpression), // TODO: separated by commas?
		parseParenRight,
		func(_ byte, exprs []int, _ byte) (int, errParse) {
			if len(exprs) == 0 {
				return _astEmptyIdentifierIdx, errParse{} // TODO: () is "nil"
			}
			// TODO: remove empty identifiers
			return ast.Append(NodeExprList{Expressions: exprs}), errParse{}
		},
	)(ast, b)
}

func parseLiteral(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseOr(
		parseNumber,
		parseString,
		parseBoolean,
	)(ast, b)
}

func parseDict(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseAnd3(
		parseBraceLeft,
		parseMany0(parseAnd3(
			parseExpression,
			parseColon,
			parseExpression,
			func(k int, _ byte, v int) (NodeObjectEntry, errParse) {
				return NodeObjectEntry{Pos{}, k, v}, errParse{}
			},
		)), // TODO: separated by commas?
		parseBraceRight,
		func(_ byte, kvs []NodeObjectEntry, _ byte) (int, errParse) {
			return ast.Append(NodeLiteralObject{Pos{}, kvs}), errParse{}
		},
	)(ast, b)
}

func parseUnary(ast *AST, b []byte) ([]byte, int, errParse) {
	return parseAnd2(
		parseByte('-'),
		parseExpression,
		func(_ byte, n int) (int, errParse) {
			return ast.Append(NodeExprUnary{Pos{}, OpNegation, n}), errParse{}
		},
	)(ast, b)
}

func parseExpression(ast *AST, b []byte) ([]byte, int, errParse) {
	b, lhs, err := parseOr(
		parseComment,
		parseLiteral,
		parseIdentifier,
		parseAssignment,
		parseBlock,
		parseList,
		parseDict,
		parseLambda,
		parseUnary,
	)(ast, b)
	if err.Err != nil {
		return nil, -1, err
	}

	{
		b, call, err := parseAnd3(
			parseParenLeft,
			parseMany0(parseExpression), // TODO: separated by commas?
			parseParenRight,
			func(_ byte, args []int, _ byte) (int, errParse) {
				return ast.Append(NodeFunctionCall{lhs, args}), errParse{}
			},
		)(ast, b)
		if err.Err == nil {
			return b, call, errParse{}
		}
	}

	{
		b, opRhs, err := parseAnd2(
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
			parseExpression,
			func(op byte, rhs int) (int, errParse) {
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
					return -1, errParse{&Err{nil, ErrSyntax, fmt.Sprintf("invalid operator %c", op), Pos{}}}
				}
				return ast.Append(NodeExprBinary{Pos{}, opKind, lhs, rhs}), errParse{}
			},
		)(ast, b)
		if err.Err == nil {
			return b, opRhs, errParse{}
		}
	}

	return b, lhs, errParse{}
}
