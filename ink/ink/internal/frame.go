package internal

import (
	"fmt"
	"sort"
	"strings"
)

// StackFrame represents the heap of variables local to a particular function call frame,
// and recursively references other parent StackFrames internally.
type StackFrame struct {
	parent *StackFrame
	vt     map[string]Value
}

// Get a value from the stack frame chain
func (frame *StackFrame) Get(name string) (Value, bool) {
	for frame != nil {
		val, ok := frame.vt[name]
		if ok {
			return val, true
		}

		frame = frame.parent
	}

	return Null, false
}

// Set a value to the most recent call stack frame
func (frame *StackFrame) Set(name string, val Value) {
	frame.vt[name] = val
}

// Up updates a value in the stack frame chain
func (frame *StackFrame) Update(name string, val Value) {
	for ; frame != nil; frame = frame.parent {
		if _, ok := frame.vt[name]; ok {
			frame.vt[name] = val
			return
		}
	}

	LogError(&Err{nil, ErrAssert, fmt.Sprintf("StackFrame.Up expected to find variable '%s' in frame but did not", name), Pos{}})
}

func (s *StackFrame) String() string {
	entries := make([]string, 0, len(s.vt))
	for k, v := range s.vt {
		vstr := v.String()
		if len(vstr) > maxPrintLen {
			vstr = vstr[:maxPrintLen] + ".."
		}
		entries = append(entries, fmt.Sprintf("%s -> %s", k, vstr))
	}

	sort.Strings(entries)

	return fmt.Sprintf("{\n\t%s\n} -prnt-> %s", strings.Join(entries, "\n\t"), s.parent)
}
