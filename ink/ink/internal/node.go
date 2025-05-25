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

func (s *AST) String() string {
	var sb strings.Builder
	for i, n := range s.Nodes {
		fmt.Fprintf(&sb, "%3d(%s): %s\n", i, n.Position(s).String(), n.String())
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
		switch n.Kind {
		case NodeKindLiteralBoolean:
			props["fillcolor"] = _colorLiteral
			val = fun.IF(n.Meta.(bool), "true", "false")
		case NodeKindLiteralNumber:
			props["fillcolor"] = _colorLiteral
			val = fmt.Sprint(n.Meta.(float64))
		case NodeKindLiteralString:
			props["fillcolor"] = _colorLiteral
			val = "'" + strings.Trim(strconv.Quote(strings.Trim(strconv.Quote(n.Meta.(string)), "\"")), "\"") + "'"
		case NodeKindLiteralList, NodeKindLiteralFunction:
			props["fillcolor"] = _colorLiteral
		case NodeKindIdentifierEmpty:
			props["fillcolor"] = _colorIdent
		case NodeKindIdentifier:
			props["fillcolor"] = _colorIdent
			val = n.Meta.(string)
		case NodeKindExprUnary:
			props["fillcolor"] = _colorExpr
			val = n.Meta.(Kind).String()
		case NodeKindExprBinary:
			props["fillcolor"] = _colorExpr
			val = n.Meta.(Kind).String()
		case NodeKindFunctionCall:
			props["fillcolor"] = _colorExpr
		case NodeKindExprMatch, NodeKindExprList:
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
		switch n.Kind {
		case NodeKindExprList:
			for k, j := range n.Children {
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"%d\"]\n", i, j, k)
			}
		case NodeKindExprUnary:
			fmt.Fprintf(&sb, "n%d -> n%d\n", i, n.Children[0])
		case NodeKindExprBinary:
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"left\"]\n", i, n.Children[0])
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"right\"]\n", i, n.Children[1])
		case NodeKindExprMatch:
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"cond\"]\n", i, n.Children[0])
			clauses := n.Children[1:]
			for k := 0; k < len(clauses); k += 2 {
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"target\"]\n", i, clauses[k])
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"expr\"]\n", i, clauses[k+1])
			}
		case NodeKindLiteralFunction:
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"body\"]\n", i, n.Children[0])
			for k, j := range n.Children[1:] {
				fmt.Fprintf(&sb, "n%d -> n%d [label=\"arg/%d\"]\n", i, j, k)
			}
		case NodeKindFunctionCall:
			fmt.Fprintf(&sb, "n%d -> n%d [label=\"func\"]\n", i, n.Children[0])
			for k, j := range n.Children[1:] {
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

type NodeKind uint8 // meta, [children]

const (
	NodeKindExprUnary        NodeKind = iota // Kind, [operand]
	NodeKindExprBinary                       // Kind, [lhs, rhs]
	NodeKindFunctionCall                     // nil, [function, args...]
	NodeKindExprMatch                        // nil, [condition, [target, expression]...]
	NodeKindExprList                         // nil, [expressions...]
	NodeKindIdentifierEmpty                  // nil, []
	NodeKindIdentifier                       // string, []
	NodeKindLiteralNumber                    // float64, []
	NodeKindLiteralString                    // string, []
	NodeKindLiteralBoolean                   // boolean, []
	NodeKindLiteralComposite                 // [n]Pos, [n][Key, Val]
	NodeKindLiteralList                      // nil, [vals...]
	NodeKindLiteralFunction                  // nil, [body, arguments...]
)

// Node represents an abstract syntax tree (AST) node in an Ink program.
type Node struct {
	Kind NodeKind
	Pos
	Meta     any
	Children []NodeID
}

func (n Node) String() string {
	switch n.Kind {
	case NodeKindExprUnary:
		return fmt.Sprintf("Unary %s #%d", n.Meta.(Kind).String(), n.Children[0])
	case NodeKindExprBinary:
		lhs, rhs := n.Children[0], n.Children[1]
		return fmt.Sprintf("Binary #%d %s #%d", lhs, n.Meta.(Kind).String(), rhs)
	case NodeKindFunctionCall:
		var sb strings.Builder
		sb.WriteString("Call ")
		sb.WriteString(n.Children[0].String())
		sb.WriteString(" on (")
		for i, a := range n.Children[1:] {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(a.String())
		}
		sb.WriteString(")")
		return sb.String()
	case NodeKindExprMatch:
		var sb strings.Builder
		sb.WriteString("Match on (")
		sb.WriteString(n.Children[0].String())
		sb.WriteString(") to {")
		clauses := n.Children[1:]
		for i := 0; i < len(clauses); i += 2 {
			if i > 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(fmt.Sprintf("Clause #%d -> #%d", clauses[i], clauses[i+1]))
		}
		sb.WriteString("}")
		return sb.String()
	case NodeKindExprList:
		var sb strings.Builder
		sb.WriteString("Expression list (")
		for i, expr := range n.Children {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(expr.String())
		}
		sb.WriteString(")")
		return sb.String()
	case NodeKindIdentifierEmpty:
		return "IDENTIFIER_EMPTY"
	case NodeKindIdentifier:
		return fmt.Sprintf("IDENTIFIER(%q)", n.Meta.(string))
	case NodeKindLiteralNumber:
		return fmt.Sprintf("Number %s", nToS(n.Meta.(float64)))
	case NodeKindLiteralString:
		return fmt.Sprintf("String %q", n.Meta.(string))
	case NodeKindLiteralBoolean:
		return fmt.Sprintf("Boolean %t", n.Meta.(bool))
	case NodeKindLiteralComposite:
		entries := make([]string, len(n.Children)/2)
		for i := range entries {
			entries[i] = fmt.Sprintf("#%d: #%d", n.Children[i*2], n.Children[i*2+1])
		}
		return fmt.Sprintf("Object {%s}", strings.Join(entries, ", "))
	case NodeKindLiteralList:
		vals := make([]string, len(n.Children))
		for i, v := range n.Children {
			vals[i] = v.String()
		}
		return fmt.Sprintf("List [%s]", strings.Join(vals, ", "))
	case NodeKindLiteralFunction:
		args := make([]string, len(n.Children[1:]))
		for i, a := range n.Children[1:] {
			args[i] = a.String()
		}
		return fmt.Sprintf("Function(%s)=>#%d", strings.Join(args, ", "), n.Children[0])
	default:
		panic("unreachable")
	}
}

func (n Node) Position(ast *AST) Pos { // TODO: replace with []Node
	switch n.Kind {
	case NodeKindExprUnary, NodeKindExprBinary,
		NodeKindExprMatch, NodeKindExprList,
		NodeKindIdentifierEmpty, NodeKindIdentifier,
		NodeKindLiteralNumber, NodeKindLiteralString, NodeKindLiteralBoolean,
		NodeKindLiteralComposite, NodeKindLiteralList, NodeKindLiteralFunction:
		return n.Pos
	case NodeKindFunctionCall:
		return ast.Nodes[n.Children[0]].Position(ast)
	default:
		panic("unreachable")
	}
}

func NodeExprUnary(Pos Pos, Operator Kind, Operand NodeID) Node {
	return Node{NodeKindExprUnary, Pos, Operator, []NodeID{Operand}}
}
func NodeExprBinary(Pos Pos, Operator Kind, Left, Right NodeID) Node {
	return Node{NodeKindExprBinary, Pos, Operator, []NodeID{Left, Right}}
}
func NodeFunctionCall(Function NodeID, Arguments []NodeID) Node {
	return Node{NodeKindFunctionCall, Pos{}, nil, append([]NodeID{Function}, Arguments...)}
}

type NodeMatchClause struct {
	Target, Expression NodeID
}

func NodeExprMatch(Pos Pos, Condition NodeID, Clauses []NodeMatchClause) Node {
	children := make([]NodeID, 1+len(Clauses)*2)
	children[0] = Condition
	for i, clause := range Clauses {
		children[1+i*2] = clause.Target
		children[1+i*2+1] = clause.Expression
	}
	return Node{NodeKindExprMatch, Pos, nil, children}
}
func NodeExprList(Pos Pos, Expressions []NodeID) Node {
	return Node{NodeKindExprList, Pos, nil, Expressions}
}
func NodeIdentifierEmpty(Pos Pos) Node {
	return Node{NodeKindIdentifierEmpty, Pos, nil, nil}
}
func NodeIdentifier(Pos Pos, Val string) Node {
	return Node{NodeKindIdentifier, Pos, Val, nil}
}
func NodeLiteralNumber(Pos Pos, Val float64) Node {
	return Node{NodeKindLiteralNumber, Pos, Val, nil}
}
func NodeLiteralString(Pos Pos, Val string) Node {
	return Node{NodeKindLiteralString, Pos, Val, nil}
}
func NodeLiteralBoolean(Pos Pos, Val bool) Node {
	return Node{NodeKindLiteralBoolean, Pos, Val, nil}
}

type NodeCompositeKeyValue struct {
	Pos
	Key, Val NodeID
}

func NodeLiteralComposite(pos Pos, Entries []NodeCompositeKeyValue) Node {
	children := make([]NodeID, len(Entries)*2)
	poss := make([]Pos, len(Entries))
	for i, entry := range Entries {
		children[i*2] = entry.Key
		children[i*2+1] = entry.Val
		poss[i] = entry.Pos
	}
	return Node{NodeKindLiteralComposite, pos, poss, children}
}
func NodeLiteralList(Pos Pos, Vals ...NodeID) Node {
	return Node{NodeKindLiteralList, Pos, nil, Vals}
}
func NodeLiteralFunction(Pos Pos, Arguments []NodeID, Body NodeID) Node {
	return Node{NodeKindLiteralFunction, Pos, nil, append([]NodeID{Body}, Arguments...)}
}
