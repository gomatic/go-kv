package kv

import (
	errs "github.com/gomatic/go-error"
)

// Every error the package can throw is declared as a const of the ecosystem's
// [errs.Const] sentinel type, so you can match each one with errors.Is rather
// than comparing strings.
const (
	// ErrNilReader is what you get when a nil reader is handed to a loader.
	ErrNilReader errs.Const = "nil reader"
	// ErrFileLoad is the leading sentinel we wrap when a source file just won't
	// load into the target environment variable.
	ErrFileLoad errs.Const = "could not load"
)
