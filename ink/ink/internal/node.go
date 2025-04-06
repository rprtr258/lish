package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rprtr258/fun"
)

type NodeID int

func (n NodeID) String() string {
	return "#" + strconv.Itoa(int(n))
}

type AST struct {
	cache map[string]NodeID
	Nodes []Node
}

func NewAst() *AST {
	return &AST{
		cache: map[string]NodeID{},
		Nodes: []Node{},
	}
}

func (s *AST) Append(node Node) NodeID {
	ss := node.String()
	if _, ok := s.cache[ss]; ok {
		return s.cache[ss]
	}

	n := NodeID(len(s.Nodes))
	s.Nodes = append(s.Nodes, node)
	s.cache[ss] = n
	return n
}

func (s AST) String() string {
	var sb strings.Builder
	for i, n := range s.Nodes {
		fmt.Fprintf(&sb, "%3d(%s): %s\n", i, n.Position(&s).String(), n.String())
	}
	return sb.String()
}

func (s AST) Graph() string {
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
		case NodeExprMatch, NodeExprList:
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
		case NodeExprMatch:
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"cond\"]\n", i, n.Condition)
			for k, clause := range n.Clauses {
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"%d\"]\n", i, clause, k) // TODO: clause is not int
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"target\"]\n", i, clause.Target)
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"expr\"]\n", i, clause.Expression)
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
	Eval(*Scope, *AST, Cont) ValueThunk
	String() string
	Position(*AST) Pos
}

type NodeExprUnary struct {
	Pos
	Operator Kind
	Operand  NodeID
}

func (n NodeExprUnary) String() string {
	return fmt.Sprintf("Unary %s #%d", n.Operator, n.Operand)
}

func (n NodeExprUnary) Position(*AST) Pos {
	return n.Pos
}

type NodeExprBinary struct {
	Pos
	Operator    Kind
	Left, Right NodeID
}

func (n NodeExprBinary) String() string {
	return fmt.Sprintf("Binary #%d %s #%d", n.Left, n.Operator, n.Right)
}

func (n NodeExprBinary) Position(*AST) Pos {
	return n.Pos
}

type NodeFunctionCall struct {
	Function  NodeID
	Arguments []NodeID
}

func (n NodeFunctionCall) String() string {
	var sb strings.Builder
	sb.WriteString("Call ")
	sb.WriteString(n.Function.String())
	sb.WriteString(" on (")
	for i, a := range n.Arguments {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(a.String())
	}
	sb.WriteString(")")
	return sb.String()
}

func (n NodeFunctionCall) Position(s *AST) Pos {
	return s.Nodes[n.Function].Position(s)
}

type NodeMatchClause struct {
	Target, Expression NodeID
}

type NodeExprMatch struct {
	Pos
	Condition NodeID
	Clauses   []NodeMatchClause
}

func (n NodeExprMatch) String() string {
	var sb strings.Builder
	sb.WriteString("Match on (")
	sb.WriteString(n.Condition.String())
	sb.WriteString(") to {")
	for i, clause := range n.Clauses {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(fmt.Sprintf("Clause #%d -> #%d", clause.Target, clause.Expression))
	}
	sb.WriteString("}")
	return sb.String()
}

func (n NodeExprMatch) Position(*AST) Pos {
	return n.Pos
}

type NodeExprList struct {
	Pos
	Expressions []NodeID
}

func (n NodeExprList) String() string {
	var sb strings.Builder
	sb.WriteString("Expression list (")
	for i, expr := range n.Expressions {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(expr.String())
	}
	sb.WriteString(")")
	return sb.String()
}

func (n NodeExprList) Position(*AST) Pos {
	return n.Pos
}

type NodeIdentifierEmpty struct {
	Pos
}

func (n NodeIdentifierEmpty) String() string {
	return "Empty Identifier"
}

func (n NodeIdentifierEmpty) Position(*AST) Pos {
	return n.Pos
}

type NodeIdentifier struct {
	Pos
	Val string
}

func (n NodeIdentifier) String() string {
	return fmt.Sprintf("Identifier '%s'", n.Val)
}

func (n NodeIdentifier) Position(*AST) Pos {
	return n.Pos
}

type NodeLiteralNumber struct {
	Pos
	Val float64
}

func (n NodeLiteralNumber) String() string {
	return fmt.Sprintf("Number %s", nToS(n.Val))
}

func (n NodeLiteralNumber) Position(*AST) Pos {
	return n.Pos
}

type NodeLiteralString struct {
	Pos
	Val string
}

func (n NodeLiteralString) String() string {
	return fmt.Sprintf("String '%s'", n.Val)
}

func (n NodeLiteralString) Position(*AST) Pos {
	return n.Pos
}

type NodeLiteralBoolean struct {
	Pos
	Val bool
}

func (n NodeLiteralBoolean) String() string {
	return fmt.Sprintf("Boolean %t", n.Val)
}

func (n NodeLiteralBoolean) Position(*AST) Pos {
	return n.Pos
}

type NodeCompositeKeyValue struct {
	Pos
	Key, Val NodeID
}

type NodeLiteralComposite struct {
	Pos
	Entries []NodeCompositeKeyValue
}

func (n NodeLiteralComposite) String() string {
	entries := make([]string, len(n.Entries))
	for i, e := range n.Entries {
		entries[i] = fmt.Sprintf("#%d: #%d", e.Key, e.Val)
	}
	return fmt.Sprintf("Object {%s}",
		strings.Join(entries, ", "))
}

func (n NodeLiteralComposite) Position(*AST) Pos {
	return n.Pos
}

type NodeLiteralList struct {
	Pos
	Vals []NodeID
}

func (n NodeLiteralList) String() string {
	vals := make([]string, len(n.Vals))
	for i, v := range n.Vals {
		vals[i] = v.String()
	}
	return fmt.Sprintf("List [%s]", strings.Join(vals, ", "))
}

func (n NodeLiteralList) Position(*AST) Pos {
	return n.Pos
}

type NodeLiteralFunction struct {
	Pos
	Arguments []NodeID
	Body      NodeID
}

func (n NodeLiteralFunction) String() string {
	args := make([]string, len(n.Arguments))
	for i, a := range n.Arguments {
		args[i] = a.String()
	}
	return fmt.Sprintf("Function (%s) => #%d", strings.Join(args, ", "), n.Body)
}

func (n NodeLiteralFunction) Position(*AST) Pos {
	return n.Pos
}
