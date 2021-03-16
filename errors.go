package kv

// Error is kv's sentinel-error type. Every error the package can throw is
// declared as a const of this type, so you can match each one with errors.Is
// rather than comparing strings. It follows the same shape as the rest of the
// ecosystem's gloo.Error / repl.Error / workgroup.Error.
type Error string

func (e Error) Error() string { return string(e) }

var _ error = Error("")

const (
	// ErrNilReader is what you get when a nil reader is handed to a loader.
	ErrNilReader Error = "nil reader"
	// ErrFileLoad is the leading sentinel we wrap when a source file just won't
	// load into the target environment variable.
	ErrFileLoad Error = "could not load"
)
