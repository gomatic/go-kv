package kv

import (
	"testing"

	errs "github.com/gomatic/go-error"
	"github.com/stretchr/testify/assert"
)

// TestSentinelsAreErrsConst pins the fleet error contract: every sentinel the
// package can emit is a constant of the ecosystem's errs.Const type from
// github.com/gomatic/go-error — never a package-local error mechanism.
func TestSentinelsAreErrsConst(t *testing.T) {
	try := assert.New(t)
	for name, sentinel := range map[string]errs.Const{
		"ErrNilReader": ErrNilReader,
		"ErrFileLoad":  ErrFileLoad,
	} {
		try.NotEmpty(sentinel.Error(), name)
	}
}
