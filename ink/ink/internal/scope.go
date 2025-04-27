package internal

import (
	"fmt"
	"strings"
)

// Scope represents the heap of variables local to a particular function call frame,
// and recursively references other parent Scopes internally.
type Scope struct {
	vt [][]Value
}

// Get a value from scope chain
func (s *Scope) Get(idx scopeIndex) (Value, bool) {
	return s.vt[idx.scope][idx.var_], true
	// return Null, false
}

// Set a value to the last scope
func (s *Scope) Set(idx scopeIndex, val Value) {
	scope := s.vt[idx.scope]
	if idx.var_ == len(scope) {
		s.vt[idx.scope] = append(scope, val)
	} else {
		scope[idx.var_] = val
	}
}

// Update updates a value in the scope chain
func (s *Scope) Update(idx scopeIndex, val Value) {
	s.vt[idx.scope][idx.var_] = val
	// log.Fatal().
	// 	Stringer("kind", ErrAssert).
	// 	Msgf("StackFrame.Up expected to find variable '%s' in frame but did not", name)
}

func (s *Scope) String() string {
	entries := make([]string, 0, len(s.vt))
	for scopeIdx, v := range s.vt {
		// vstr := v.String()
		// if len(vstr) > maxPrintLen {
		// 	vstr = vstr[:maxPrintLen] + ".."
		// }
		// entries = append(entries, fmt.Sprintf("%d[%s]", scopeIdx, vstr))
		entries = append(entries, fmt.Sprintf("%d: %v", scopeIdx, v))
	}

	// return fmt.Sprintf("{\n\t%s\n} -prnt-> %s", strings.Join(entries, "\n\t"), s.parent)
	return strings.Join(entries, "\n\t")
}
