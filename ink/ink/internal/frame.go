package internal

import (
	"fmt"
	"iter"
	"slices"
	"sort"
	"strings"

	"github.com/rprtr258/fun"
	"github.com/rprtr258/scuf"
)

type StackFrameID int

type StackFrame struct {
	parent StackFrameID
	vt     map[string]Value
}

// TheStack represents the heap of variables local to a particular function call frame,
// and recursively references other parent StackFrames internally.
type TheStack struct {
	frames []StackFrame
}

func (frame StackFrameID) String() string {
	return fmt.Sprintf("@%d", frame)
}

func (stack TheStack) history(frameID StackFrameID) iter.Seq2[StackFrameID, StackFrame] {
	return func(yield func(StackFrameID, StackFrame) bool) {
		for id := frameID; id != -1; id = stack.frames[id].parent {
			if !yield(id, stack.frames[id]) {
				return
			}
		}
	}
}

func (stack TheStack) history2(frameID StackFrameID) []StackFrameID {
	res := []StackFrameID{}
	for id := range stack.history(frameID) {
		res = append(res, id)
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

// Get a value from the stack frame chain
func (stack TheStack) Get(frameID StackFrameID, name string) (Value, bool) {
	for _, frame := range stack.history(frameID) {
		val, ok := frame.vt[name]
		if ok {
			return val, true
		}
	}

	return Null, false
}

// Set a value to the most recent call stack frame
func (stack TheStack) Set(frameID StackFrameID, name string, val Value) {
	frame := stack.frames[frameID]
	frame.vt[name] = val
}

// Up updates a value in the stack frame chain
func (stack TheStack) Update(frameID StackFrameID, name string, val Value) {
	for _, frame := range stack.history(frameID) {
		if _, ok := frame.vt[name]; ok {
			frame.vt[name] = val
			return
		}
	}

	LogError(&Err{nil, ErrAssert, fmt.Sprintf("StackFrame.Up expected to find variable '%s' in frame but did not", name), Pos{}})
}

func (stack TheStack) String(frameID StackFrameID) string {
	frames := []string{}
	for id, frame := range stack.history(frameID) {
		frames = append(frames, scuf.NewString(func(sb scuf.Buffer) {
			sb.String(id.String(), scuf.FgHiBlack)
			sb.InBytePair('{', '}', func(b scuf.Buffer) {
				if frame.parent == id {
					sb.String("#root")
					// TODO: print not native funcs also, when they appear
					return
				}

				keys := fun.Keys(frame.vt)
				sort.Strings(keys)
				for i, k := range keys {
					v := frame.vt[k]
					if v, ok := v.(NativeFunctionValue); ok && v.name == k {
						continue
					}

					if i > 0 {
						sb.String(" ")
					}
					sb.String(k)
					sb.String("=", scuf.FgBlack)
					sb.String(v.String())
				}
			})
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

func (stack TheStack) string(frameID StackFrameID) string {
	var sb strings.Builder
	sb.WriteString(frameID.String())
	sb.WriteString("(")
	for _, id := range stack.history2(frameID) {
		sb.WriteString(id.String())
		sb.WriteString(" ")
	}
	sb.WriteString(")")
	return sb.String()
}

func (stack TheStack) LeastCommonAncestor(frame, f StackFrameID) StackFrameID {
	/*
		possible cases:
		1. 0 <- * <- (f) <- * <- (frame) // return f.parent
		2. 0 <- * <- (frame) <- * <- (f) // i hope this is actually impossible
		3.             v- * <- (f)
		   0 <- * <- (lca)
		               ^- * <- (frame)
		  // return lca
	*/

	frameHistory := stack.history2(frame)
	fHistory := stack.history2(f)
	for i := 0; i < len(frameHistory) && i < len(fHistory); i++ {
		if frameHistory[i] != fHistory[i] {
			if i == 0 {
				panic(fmt.Sprintf("unreachable LCA(%s, %s)", stack.string(frame), stack.string(f)))
			}
			return frameHistory[i-1]
		}
	}
	return stack.frames[f].parent
}

func (stack *TheStack) Rebase(frame, head StackFrameID) StackFrameID {
	if head == frame {
		return frame
	}

	lca := stack.LeastCommonAncestor(frame, head)
	fmt.Printf("LCA(\n  %s\n  %s\n) = %s\n", stack.string(frame), stack.string(head), stack.string(lca))
	return stack.rebase(frame, head, lca)
}

func (stack *TheStack) append(parentID StackFrameID, vt map[string]Value) StackFrameID {
	frame := StackFrame{parentID, vt}
	stack.frames = append(stack.frames, frame)
	return StackFrameID(len(stack.frames) - 1)
}

func (stack *TheStack) Append(parentID StackFrameID) StackFrameID {
	return stack.append(parentID, map[string]Value{})
}

func (stack *TheStack) rebase(frame, head, lca StackFrameID) StackFrameID {
	framesToRebase := []StackFrameID{}
	for head != lca {
		framesToRebase = append(framesToRebase, head)
		head = stack.frames[head].parent
	}

	res := frame
	for _, id := range slices.Backward(framesToRebase) {
		res = stack.append(res, stack.frames[id].vt)
	}
	return res
}
