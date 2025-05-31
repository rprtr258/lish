package internal

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"unsafe"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/scuf"
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
	frames := []string{}
	for ; frame != nil; frame = frame.parent {
		frames = append(frames, scuf.NewString(func(sb scuf.Buffer) {
			sb.String(fmt.Sprintf("%p", frame)[6:], scuf.FgHiBlack)
			sb.String("{")
			if frame.parent == nil {
				sb.String("#root")
			} else if len(frame.vt) > 0 {
				keys := fun.Keys(frame.vt)
				sort.Strings(keys)

				for i, k := range keys {
					if i > 0 {
						sb.String(" ")
					}
					sb.String(k)
					sb.String("=", scuf.FgBlack)
					sb.String(frame.vt[k].String())
				}
			}
			sb.String("}")
		}))
	}
	frames = slices.Collect(func(yield func(string) bool) {
		for _, frame := range slices.Backward(frames) {
			if !yield(frame) {
				break
			}
		}
	})
	return strings.Join(frames, "\n")
}

func (frame *StackFrame) LeastCommonAncestor(f *StackFrame) *StackFrame {
	for frame != f {
		if uintptr(unsafe.Pointer(frame)) > uintptr(unsafe.Pointer(f)) {
			frame = frame.parent
		} else {
			f = f.parent
		}
	}
	return frame
}

func (frame *StackFrame) Rebase(head *StackFrame) *StackFrame {
	lca := frame.LeastCommonAncestor(head)
	return frame.rebase(head, lca)
}

func (frame *StackFrame) rebase(head, lca *StackFrame) *StackFrame {
	if head == lca {
		return frame
	}

	prehead := frame.rebase(head.parent, lca)
	return &StackFrame{prehead, head.vt}
}
