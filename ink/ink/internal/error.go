package internal

import "fmt"

type ErrorReason int8

func (r ErrorReason) String() string {
	switch r {
	case ErrSyntax:
		return "syntax"
	case ErrRuntime:
		return "runtime"
	case ErrSystem:
		return "system"
	case ErrAssert:
		return "assert"
	default:
		return "unknown"
	}
}

// Error reasons are enumerated here to be used in the Err struct,
// the error type shared across all Ink APIs.
const (
	ErrUnknown ErrorReason = iota
	ErrSyntax
	ErrRuntime
	ErrSystem
	ErrAssert
)

// Err constants represent possible errors that Ink interpreter binding functions may return.
type Err struct {
	Parent  *Err
	Reason  ErrorReason
	Message string
	Pos     Pos
}

func (e *Err) Error() string {
	// TODO: print beautiful stack trace
	// TODO: skip reasons down the stack
	message := e.Message
	if e.Parent != nil {
		message = message + ": " + e.Parent.Error()
	}

	return fmt.Sprintf("%s error: %s [%s]", e.Reason, message, e.Pos)
}
