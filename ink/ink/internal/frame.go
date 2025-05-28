package internal

import (
	"fmt"
	"sort"
	"strings"
)

// StackFrame represents the heap of variables local to a particular function call frame,
// and recursively references other parent StackFrames internally.
type StackFrame struct { // TODO: unembed
	parent *StackFrame
	vt     map[string]Value
}

// Get a value from the stack frame chain
func (frame *StackFrame) Get(name string) (Value, bool) {
	for ; frame != nil; frame = frame.parent {
		val, ok := frame.vt[name]
		if ok {
			return val, true
		}
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

func (frame *StackFrame) String() string {
	var sb strings.Builder
	for ; frame != nil; frame = frame.parent {
		if len(frame.vt) > 0 {
			entries := make([]string, 0, len(frame.vt))
			for k, v := range frame.vt {
				entries = append(entries, fmt.Sprintf("%s : %s", k, v.String()))
			}
			sort.Strings(entries)

			sb.WriteString(fmt.Sprintf("{\n\t%s\n}", strings.Join(entries, "\n\t")))
		} else {
			sb.WriteString(fmt.Sprintf("{}"))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
