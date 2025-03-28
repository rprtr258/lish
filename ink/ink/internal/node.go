package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rprtr258/fun"
)

const (
	_astEmptyIdentifierIdx = 0
	_astFalseLiteralIdx    = 1
	_astTrueLiteralIdx     = 2
)

type AST struct {
	Numbers     map[float64]int
	Identifiers map[string]int
	Strings     map[string]int
	Nodes       []Node
}

func NewAstSlice() *AST {
	return &AST{
		Numbers:     map[float64]int{},
		Identifiers: map[string]int{},
		Strings:     map[string]int{},
		Nodes: []Node{
			NodeIdentifierEmpty{},
			NodeLiteralBoolean{Val: false},
			NodeLiteralBoolean{Val: true},
		},
	}
}

func (s *AST) Append(node Node) int {
	fmt.Fprint(os.Stderr, len(s.Nodes), " ")
	LogNode(node)
	// TODO: position is lost
	switch node := node.(type) {
	case NodeIdentifierEmpty:
		// one empty identifier for all
		return _astEmptyIdentifierIdx
	case NodeLiteralBoolean:
		return fun.IF(node.Val, _astTrueLiteralIdx, _astFalseLiteralIdx)
	case NodeLiteralNumber:
		if _, ok := s.Numbers[node.Val]; !ok {
			n := len(s.Nodes)
			s.Nodes = append(s.Nodes, node)
			s.Numbers[node.Val] = n
		}
		return s.Numbers[node.Val]
	case NodeLiteralString:
		if _, ok := s.Strings[node.Val]; !ok {
			n := len(s.Nodes)
			s.Nodes = append(s.Nodes, node)
			s.Strings[node.Val] = n
		}
		return s.Strings[node.Val]
	case NodeIdentifier:
		if _, ok := s.Identifiers[node.Val]; !ok {
			n := len(s.Nodes)
			s.Nodes = append(s.Nodes, node)
			s.Identifiers[node.Val] = n
		}
		return s.Identifiers[node.Val]
	default:
		n := len(s.Nodes)
		s.Nodes = append(s.Nodes, node)
		return n
	}
}

func (s AST) String() string {
	const (
		_colorLiteral = "#fffd87"
		_colorIdent   = "#88ff00"
		_colorExpr    = "#a1c5ff"
		_colorControl = "#ffa1a1"
	)

	var sb strings.Builder
	sb.WriteString("digraph G {\n")
	sb.WriteString(`node [
	style=filled
	shape=rect
	fillcolor=white
	fontname="Helvetica,Arial,sans-serif"
]
`)
	for i, n := range s.Nodes {
		typ := strings.TrimPrefix(fmt.Sprintf("%T", n), "internal.Node")
		var val string
		props := map[string]string{}
		switch n := n.(type) {
		case NodeLiteralBoolean:
			props["fillcolor"] = _colorLiteral
			val = fun.IF(n.Val, "true", "false")
		case NodeLiteralNumber:
			props["fillcolor"] = _colorLiteral
			val = fmt.Sprint(n.Val)
		case NodeLiteralString:
			props["fillcolor"] = _colorLiteral
			val = "'" + strings.Trim(strconv.Quote(strings.Trim(strconv.Quote(n.Val), "\"")), "\"") + "'"
		case NodeLiteralList, NodeLiteralFunction:
			props["fillcolor"] = _colorLiteral
		case NodeIdentifierEmpty:
			props["fillcolor"] = _colorIdent
		case NodeIdentifier:
			props["fillcolor"] = _colorIdent
			val = n.Val
		case NodeExprUnary:
			props["fillcolor"] = _colorExpr
			val = n.Operator.String()
		case NodeExprBinary:
			props["fillcolor"] = _colorExpr
			val = n.Operator.String()
		case NodeFunctionCall:
			props["fillcolor"] = _colorExpr
		case NodeMatchClause, NodeMatchExpr, NodeExprList:
			props["fillcolor"] = _colorControl
		}
		props["label"] = fmt.Sprintf(`#%[1]d %[2]s\n%s`, i, typ, val)

		fmt.Fprintf(&sb, "n%d [\n", i)
		for k, v := range props {
			fmt.Fprintf(&sb, "  %s=\"%s\"\n", k, v)
		}
		fmt.Fprint(&sb, "]\n")
	}
	for i, n := range s.Nodes {
		switch n := n.(type) {
		case NodeExprList:
			for k, j := range n.Expressions {
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"%d\"]\n", i, j, k)
			}
		case NodeExprUnary:
			fmt.Fprintf(&sb, "n%d -> n%d\n", i, n.Operand)
		case NodeExprBinary:
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"left\"]\n", i, n.Left)
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"right\"]\n", i, n.Right)
		case NodeMatchClause:
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"target\"]\n", i, n.Target)
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"expr\"]\n", i, n.Expression)
		case NodeMatchExpr:
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"cond\"]\n", i, n.Condition)
			for k, j := range n.Clauses {
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"%d\"]\n", i, j, k)
			}
		case NodeLiteralFunction:
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"body\"]\n", i, n.Body)
			for k, j := range n.Arguments {
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"arg/%d\"]\n", i, j, k)
			}
		case NodeFunctionCall:
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"func\"]\n", i, n.Function)
			for k, j := range n.Arguments {
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"app/%d\"]\n", i, j, k)
			}
		}
	}
	sb.WriteString("}")
	return sb.String()
}

type errParse struct {
	*Err
}

// Node represents an abstract syntax tree (AST) node in an Ink program.
type Node interface {
	String() string
	Position() Pos
	Eval(*Scope, *AST, bool) (Value, *Err)
}

type NodeExprUnary struct {
	Pos
	Operator Kind
	Operand  int
}

func (n NodeExprUnary) String() string {
	return fmt.Sprintf("Unary %s #%d", n.Operator, n.Operand)
}

func (n NodeExprUnary) Position() Pos {
	return n.Pos
}

type NodeExprBinary struct {
	Pos
	Operator    Kind
	Left, Right int
}

func (n NodeExprBinary) String() string {
	return fmt.Sprintf("Binary #%d %s #%d", n.Left, n.Operator, n.Right)
}

func (n NodeExprBinary) Position() Pos {
	return n.Pos
}

type NodeFunctionCall struct {
	Function  int
	Arguments []int
}

func (n NodeFunctionCall) String() string {
	var sb strings.Builder
	sb.WriteString("Call (#")
	sb.WriteString(strconv.Itoa(n.Function))
	sb.WriteString(") on (")
	for i, a := range n.Arguments {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("#")
		sb.WriteString(strconv.Itoa(a))
	}
	sb.WriteString(")")
	return sb.String()
}

func (n NodeFunctionCall) Position() Pos {
	return Pos{} // TODO: n.function.Position()
}

type NodeMatchClause struct {
	Target, Expression int
}

func (n NodeMatchClause) String() string {
	return fmt.Sprintf("Clause #%d -> #%d", n.Target, n.Expression)
}

func (n NodeMatchClause) Position() Pos {
	return Pos{} // TODO: return n.target.Position()
}

type NodeMatchExpr struct {
	Pos
	Condition int
	Clauses   []int // TODO: inline clauses
}

func (n NodeMatchExpr) String() string {
	var sb strings.Builder
	sb.WriteString("Match on (#")
	sb.WriteString(strconv.Itoa(n.Condition))
	sb.WriteString(") to {")
	for i, a := range n.Clauses {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("#")
		sb.WriteString(strconv.Itoa(a))
	}
	sb.WriteString("}")
	return sb.String()
}

func (n NodeMatchExpr) Position() Pos {
	return n.Pos
}

type NodeExprList struct {
	Pos
	Expressions []int
}

func (n NodeExprList) String() string {
	var sb strings.Builder
	sb.WriteString("Expression list (")
	for i, expr := range n.Expressions {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("#")
		sb.WriteString(strconv.Itoa(expr))
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
	Val string
}

func (n NodeIdentifier) String() string {
	return fmt.Sprintf("Identifier '%s'", n.Val)
}

func (n NodeIdentifier) Position() Pos {
	return n.Pos
}

type NodeLiteralNumber struct {
	Pos
	Val float64
}

func (n NodeLiteralNumber) String() string {
	return fmt.Sprintf("Number %s", nToS(n.Val))
}

func (n NodeLiteralNumber) Position() Pos {
	return n.Pos
}

type NodeLiteralString struct {
	Pos
	Val string
}

func (n NodeLiteralString) String() string {
	return fmt.Sprintf("String '%s'", n.Val)
}

func (n NodeLiteralString) Position() Pos {
	return n.Pos
}

type NodeLiteralBoolean struct {
	Pos
	Val bool
}

func (n NodeLiteralBoolean) String() string {
	return fmt.Sprintf("Boolean %t", n.Val)
}

func (n NodeLiteralBoolean) Position() Pos {
	return n.Pos
}

type NodeLiteralObject struct {
	Pos
	Entries []NodeObjectEntry
}

func (n NodeLiteralObject) String() string {
	entries := make([]string, len(n.Entries))
	for i, e := range n.Entries {
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
	Key, Val int
}

func (n NodeObjectEntry) String() string {
	return fmt.Sprintf("#%d: #%d", n.Key, n.Val)
}

type NodeLiteralList struct {
	Pos
	Vals []int
}

func (n NodeLiteralList) String() string {
	vals := make([]string, len(n.Vals))
	for i, v := range n.Vals {
		vals[i] = "#" + strconv.Itoa(v)
	}
	return fmt.Sprintf("List [%s]", strings.Join(vals, ", "))
}

func (n NodeLiteralList) Position() Pos {
	return n.Pos
}

type NodeLiteralFunction struct {
	Pos
	Arguments []int
	Body      int
}

func (n NodeLiteralFunction) String() string {
	args := make([]string, len(n.Arguments))
	for i, a := range n.Arguments {
		args[i] = "#" + strconv.Itoa(a)
	}
	return fmt.Sprintf("Function (%s) => #%d", strings.Join(args, ", "), n.Body)
}

func (n NodeLiteralFunction) Position() Pos {
	return n.Pos
}
