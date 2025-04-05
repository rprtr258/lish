package internal

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"strconv"
	"unicode"
)

// Kind is the sum type of all possible types of tokens in an Ink program
// TODO: separate unary and binary kinds of operators
type Kind int

const (
	Separator Kind = iota

	ExprUnary
	ExprBinary
	ExprMatch
	MatchClause

	Identifier
	IdentifierEmpty

	FunctionCall

	LiteralNumber
	LiteralString
	LiteralObject
	LiteralList
	LiteralFunction

	LiteralTrue
	LiteralFalse

	// ambiguous operators and symbols
	OpAccessor

	// =
	OpEqual
	FunctionArrow

	// :
	KeyValueSeparator
	OpDefine
	MatchColon

	// -
	CaseArrow
	OpSubtract

	// single char, unambiguous
	OpNegation
	OpAdd
	OpMultiply
	OpDivide
	OpModulus
	OpGreaterThan
	OpLessThan

	OpLogicalAnd
	OpLogicalOr
	OpLogicalXor

	ParenLeft
	ParenRight
	BracketLeft
	BracketRight
	BraceLeft
	BraceRight
)

func (kind Kind) String() string {
	switch kind {
	case ExprUnary:
		return "unary expression"
	case ExprBinary:
		return "binary expression"
	case ExprMatch:
		return "match expression"
	case MatchClause:
		return "match clause"

	case Identifier:
		return "identifier"
	case IdentifierEmpty:
		return "'_'"

	case FunctionCall:
		return "function call"

	case LiteralNumber:
		return "number literal"
	case LiteralString:
		return "string literal"
	case LiteralObject:
		return "composite literal"
	case LiteralList:
		return "list composite literal"
	case LiteralFunction:
		return "function literal"

	case LiteralTrue:
		return "'true'"
	case LiteralFalse:
		return "'false'"

	case OpAccessor:
		return "'.'"

	case OpEqual:
		return "'='"
	case FunctionArrow:
		return "'=>'"

	case KeyValueSeparator:
		return "':'"
	case OpDefine:
		return "':='"
	case MatchColon:
		return "'::'"

	case CaseArrow:
		return "'->'"
	case OpSubtract:
		return "'-'"

	case OpNegation:
		return "'~'"
	case OpAdd:
		return "'+'"
	case OpMultiply:
		return "'*'"
	case OpDivide:
		return "'/'"
	case OpModulus:
		return "'%'"
	case OpGreaterThan:
		return "'>'"
	case OpLessThan:
		return "'<'"

	case OpLogicalAnd:
		return "'&'"
	case OpLogicalOr:
		return "'|'"
	case OpLogicalXor:
		return "'^'"

	case Separator:
		return "','"
	case ParenLeft:
		return "'('"
	case ParenRight:
		return "')'"
	case BracketLeft:
		return "'['"
	case BracketRight:
		return "']'"
	case BraceLeft:
		return "'{'"
	case BraceRight:
		return "'}'"

	default:
		return "unknown token"
	}
}

type Pos struct {
	File      string
	Line, Col int
}

func (p Pos) String() string {
	if p == (Pos{}) {
		return "??:??:??"
	}

	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Col)
}

// Token is the monomorphic struct representing all Ink program tokens in the lexer
type Token struct {
	Kind Kind
	Pos

	// for string/number literals
	Str string
	Num float64
}

func (tok Token) String() string {
	switch tok.Kind {
	case Identifier, LiteralString:
		return fmt.Sprintf("%s '%s' [%s]", tok.Kind, tok.Str, tok.Pos)
	case LiteralNumber:
		return fmt.Sprintf("%s '%s' [%s]", tok.Kind, nToS(tok.Num), tok.Pos)
	default:
		return fmt.Sprintf("%s [%s]", tok.Kind, tok.Pos)
	}
}

// tokenize takes an io.Reader and transforms it into a stream of Tok (tokens).
func tokenize(file string, r io.Reader) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		var buf, strbuf []rune
		var strbufStartLine, strbufStartCol int

		lastKind := Separator
		lineNo, colNo := 1, 1

		simpleCommit := func(tok Token) {
			lastKind = tok.Kind
			LogToken(tok)
			yield(tok) // TODO: break on false
		}
		simpleCommitChar := func(kind Kind) {
			simpleCommit(Token{
				Kind: kind,
				Pos:  Pos{file, lineNo, colNo},
			})
		}
		commitClear := func() {
			if len(buf) == 0 {
				// no need to commit empty token
				return
			}

			cbuf := buf
			buf = nil
			switch {
			case string(cbuf) == "true":
				simpleCommitChar(LiteralTrue)
			case string(cbuf) == "false":
				simpleCommitChar(LiteralFalse)
			case unicode.IsDigit(rune(cbuf[0])):
				f, err := strconv.ParseFloat(string(cbuf), 64)
				if err != nil {
					LogError(&Err{nil, ErrSyntax, fmt.Sprintf("can't parse number: %s", err.Error()), Pos{file, lineNo, colNo - len(cbuf)}})
				}
				simpleCommit(Token{
					Num:  f,
					Kind: LiteralNumber,
					Pos:  Pos{file, lineNo, colNo - len(cbuf)},
				})
			default:
				simpleCommit(Token{
					Str:  string(cbuf),
					Kind: Identifier,
					Pos:  Pos{file, lineNo, colNo - len(cbuf)},
				})
			}
		}
		commit := func(tok Token) {
			commitClear()
			simpleCommit(tok)
		}
		commitChar := func(kind Kind) {
			commit(Token{
				Kind: kind,
				Pos:  Pos{file, lineNo, colNo},
			})
		}
		ensureSeparator := func() {
			commitClear()
			switch lastKind {
			case Separator, ParenLeft, BracketLeft, BraceLeft,
				OpAdd, OpSubtract, OpMultiply, OpDivide, OpModulus, OpNegation,
				OpGreaterThan, OpLessThan, OpEqual, OpDefine, OpAccessor,
				KeyValueSeparator, FunctionArrow, MatchColon, CaseArrow:
				// do nothing
			default:
				commitChar(Separator)
			}
		}

		inStringLiteral := false
		br := bufio.NewReader(r)
		for {
			char, _, err := br.ReadRune()
			if err != nil {
				break
			}

		OUTER:
			switch {
			case char == '\'' && !inStringLiteral:
				strbufStartLine, strbufStartCol = lineNo, colNo
				inStringLiteral = true
			case char == '\'' && inStringLiteral:
				commit(Token{
					Str:  string(strbuf),
					Kind: LiteralString,
					Pos:  Pos{file, strbufStartLine, strbufStartCol},
				})
				strbuf = strbuf[:0]
				inStringLiteral = false
			case inStringLiteral:
				switch char {
				case '\n':
					lineNo++
					colNo = 0
					strbuf = append(strbuf, char)
				case '\\':
					c, _, err := br.ReadRune()
					if err != nil {
						break OUTER
					}
					switch c {
					case '\\', '\'':
						strbuf = append(strbuf, c)
					case 'n':
						strbuf = append(strbuf, '\n')
					case 'r':
						strbuf = append(strbuf, '\r')
					case 't':
						strbuf = append(strbuf, '\t')
					default:
						strbuf = append(strbuf, '\\', c)
					}
					colNo++
				default:
					strbuf = append(strbuf, char)
				}
			case char == '#': // single-line comment, keep taking until EOL
				for {
					if c, _, err := br.ReadRune(); err != nil || c == '\n' {
						if c == '\n' {
							br.UnreadRune()
						}
						break
					}
				}
				continue
			case char == '`': // multi-line block comment, keep taking until end of block
				nextChar, _, err := br.ReadRune()
				if err != nil {
					break
				}

				for nextChar != '`' {
					nextChar, _, err = br.ReadRune()
					if err != nil {
						break
					}

					if nextChar == '\n' {
						lineNo++
						colNo = 0
					}
					colNo++
				}
			case char == '\n':
				ensureSeparator()
				lineNo++
				colNo = 0
			case unicode.IsSpace(char):
				commitClear()
			case char == '_' && len(buf) == 0:
				commitChar(IdentifierEmpty)
			case char == '~':
				commitChar(OpNegation)
			case char == '+':
				commitChar(OpAdd)
			case char == '*':
				commitChar(OpMultiply)
			case char == '/':
				commitChar(OpDivide)
			case char == '%':
				commitChar(OpModulus)
			case char == '&':
				commitChar(OpLogicalAnd)
			case char == '|':
				commitChar(OpLogicalOr)
			case char == '^':
				commitChar(OpLogicalXor)
			case char == '<':
				commitChar(OpLessThan)
			case char == '>':
				commitChar(OpGreaterThan)
			case char == ',':
				commitChar(Separator)
			case char == '.':
				// only non-AccessorOp case is [Number token] . [Number],
				// so we commit and bail early if the buf is empty or contains
				// a clearly non-numeric token. Note that this means all numbers
				// must start with a digit. i.e. .5 is not 0.5 but a syntax error.
				// This is the case since we don't know what the last token was,
				// and I think streaming parse is worth the tradeoffs of losing
				// that context.
				committed := false
				for _, d := range buf {
					if !unicode.IsDigit(d) {
						commitChar(OpAccessor)
						committed = true
						break
					}
				}
				if !committed {
					if len(buf) == 0 {
						commitChar(OpAccessor)
					} else {
						buf = append(buf, '.')
					}
				}
			case char == ':':
				nextChar, _, err := br.ReadRune()
				if err != nil {
					break
				}

				colNo++
				switch nextChar {
				case '=':
					commitChar(OpDefine)
				case ':':
					commitChar(MatchColon)
				default:
					// key is parsed as expression, so make sure
					// we mark expression end (Separator)
					ensureSeparator()
					commitChar(KeyValueSeparator)
					br.UnreadRune()
				}
			case char == '=':
				nextChar, _, err := br.ReadRune()
				if err != nil {
					break
				}

				colNo++
				if nextChar == '>' {
					commitChar(FunctionArrow)
				} else {
					commitChar(OpEqual)
					br.UnreadRune()
				}
			case char == '-':
				nextChar, _, err := br.ReadRune()
				if err != nil {
					break
				}

				colNo++
				if nextChar == '>' {
					commitChar(CaseArrow)
				} else {
					commitChar(OpSubtract)
					br.UnreadRune()
				}
			case char == '(':
				commitChar(ParenLeft)
			case char == ')':
				ensureSeparator()
				commitChar(ParenRight)
			case char == '[':
				commitChar(BracketLeft)
			case char == ']':
				ensureSeparator()
				commitChar(BracketRight)
			case char == '{':
				commitChar(BraceLeft)
			case char == '}':
				ensureSeparator()
				commitChar(BraceRight)
			default:
				buf = append(buf, char)
			}
			colNo++
		}

		ensureSeparator()
	}
}
