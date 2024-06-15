package main

import (
	"regexp"
	"strconv"

	"github.com/rprtr258/fun"
)

var (
	_reInt   = regexp.MustCompile(`^-?\d+$`)
	_reFloat = regexp.MustCompile(`^-?\d+\.\d*$`)
)

func readAtom(token string) Atom {
	if token == "true" {
		return atomBool(true)
	} else if token == "false" {
		return atomBool(false)
	}

	if _reInt.MatchString(token) {
		n, _ := strconv.ParseInt(token, 10, 64)
		return atomInt(n)
	}

	if _reFloat.MatchString(token) {
		n, _ := strconv.ParseFloat(token, 64)
		return atomFloat(n)
	}

	if s, err := strconv.Unquote(token); err == nil {
		return atomString(s)
	}

	return atomSymbol(token)
}

type Frame struct {
	a             Atom
	isReaderMacro bool
}

type stack struct {
	data []Frame
}

func (s *stack) push(a Atom, isReaderMacro bool) {
	s.data = append(s.data, Frame{a, isReaderMacro})
}

func (s *stack) pop() (Frame, bool) {
	if len(s.data) == 0 {
		return Frame{}, false
	}

	res := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return res, true
}

// TODO: reader macro list, (add run-time)?
func read_form(tokens []string) Atom {
	lists_stack := stack{}
	append_item_to_last_stack_list := func(lists_stack *stack, item Atom) {
		n := len(lists_stack.data)

		if n == 0 {
			lists_stack.push(atomList(item), false)
			return
		}

		last := lists_stack.data[n-1]
		switch last.a.Kind {
		case AtomKindList:
			lists_stack.data[n-1].a = Atom{AtomKindList, append(last.a.Value.([]Atom), item)}
		default:
			panic("unimplemented")
		}
	}
	for i := 0; i < len(tokens); i++ {
		switch token := tokens[i]; token {
		case "(":
			lists_stack.push(atomNil, false)
		case ")":
			if i == len(tokens)-1 {
				continue
			}
			last_list, _ := lists_stack.pop()
			append_item_to_last_stack_list(&lists_stack, last_list.a)
		case "{":
			hashmap := map[string]Atom{}
			for i+1 < len(tokens) && tokens[i+1] != "}" {
				i++
				a := readAtom(tokens[i])
				if a.Kind != AtomKindString {
					panic("TODO: Not a valid key")
				}
				key := a.Value.(string)
				i++
				value := readAtom(tokens[i]) // TODO: eval
				hashmap[key] = value
			}
			if i+1 == len(tokens) || tokens[i+1] == "}" {
				i++ // "}"
			}
			append_item_to_last_stack_list(&lists_stack, atomHash(hashmap))
		case "'":
			lists_stack.push(atomList(atomSymbol("quote")), true)
		case "`":
			lists_stack.push(atomList(atomSymbol("quasiquote")), true)
		case ",":
			lists_stack.push(atomList(atomSymbol("unquote")), true)
		case ",@":
			lists_stack.push(atomList(atomSymbol("splice-unquote")), true)
		default:
			item := readAtom(token)
			append_item_to_last_stack_list(&lists_stack, item)
		}
		for len(lists_stack.data) > 1 && lists_stack.data[len(lists_stack.data)-1].isReaderMacro && len(lists_stack.data[len(lists_stack.data)-1].a.Value.([]Atom)[1:]) == 1 {
			last_list, _ := lists_stack.pop()
			append_item_to_last_stack_list(&lists_stack, last_list.a)
		}
	}
	for len(lists_stack.data) > 1 {
		last_list, _ := lists_stack.pop()
		append_item_to_last_stack_list(&lists_stack, last_list.a)
	}

	a, ok := lists_stack.pop()
	return fun.IF(ok, a.a, atomNil)
}

var RE = regexp.MustCompile(`\s*(,@|[{}()'` + "`" + `,^@]|"(?:\\.|[^\\"])*"|;.*|[^\s{}()'"` + "`" + `,;]*)\s*`)

func read(cmd string) Atom {
	tokens := []string{}
	for _, submatch := range RE.FindAllStringSubmatch(cmd, -1) {
		token := submatch[1]
		if token == "" || token[0] == ';' {
			continue
		}
		tokens = append(tokens, token)
	}
	return read_form(tokens)
}
