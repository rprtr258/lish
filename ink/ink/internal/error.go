package internal

import "fmt"

type ErrorReason int

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
	reason  ErrorReason
	message string
}

func (e Err) Error() string {
	return fmt.Sprintf("%s error: %s", e.reason, e.message)
}
