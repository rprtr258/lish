package internal

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

// ValueTable is used anytime a map of names/labels to Ink Values is needed,
// and is notably used to represent stack frames / heaps and CompositeValue dictionaries.
type ValueTable = map[string]Value

// Scope represents the heap of variables local to a particular function call frame,
// and recursively references other parent Scopes internally.
type Scope struct {
	parent *Scope
	vt     ValueTable
}

// Get a value from scope chain
func (s *Scope) Get(name string) (Value, bool) {
	for s != nil {
		if val, ok := s.vt[name]; ok {
			return val, true
		}

		s = s.parent
	}

	return Null, false
}

// Set a value to the last scope
func (s *Scope) Set(name string, val Value) {
	s.vt[name] = val
}

// Update updates a value in the scope chain
func (s *Scope) Update(name string, val Value) {
	for s != nil {
		if _, ok := s.vt[name]; ok {
			s.vt[name] = val
			return
		}

		s = s.parent
	}

	log.Fatal().Stringer("kind", ErrAssert).Msgf("StackFrame.Up expected to find variable '%s' in frame but did not", name)
}

func (s *Scope) String() string {
	entries := make([]string, 0, len(s.vt))
	for k, v := range s.vt {
		vstr := v.String()
		if len(vstr) > maxPrintLen {
			vstr = vstr[:maxPrintLen] + ".."
		}
		entries = append(entries, fmt.Sprintf("%s -> %s", k, vstr))
	}

	return fmt.Sprintf("{\n\t%s\n} -prnt-> %s", strings.Join(entries, "\n\t"), s.parent)
}
