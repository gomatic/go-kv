package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNames_Get(t *testing.T) {
	t.Setenv("KV_TEST_NAMES", "from-os")
	try := assert.New(t)

	try.Equal("", Names(nil).Get("anything"))
	try.Equal("", Names{}.Get(""))
	try.Equal("override", Names{"KV_TEST_NAMES": "override"}.Get("KV_TEST_NAMES"))
	try.Equal("from-os", Names{}.Get("KV_TEST_NAMES"))
}

func TestNames_Expand(t *testing.T) {
	try := assert.New(t)

	try.Equal(Names{}, Names{}.Expand())

	original := Names{
		"A":          "1",
		"B":          "${A}",       // this one resolves to A's value
		"C":          "${MISSING}", // nothing to resolve, so it's kept as-is
		"${MISSING}": "key-ref",    // same deal for the key — left untouched
	}
	expanded := original.Expand()

	try.Equal("1", expanded["A"])
	try.Equal("1", expanded["B"])
	try.Equal("${MISSING}", expanded["C"])
	try.Equal("key-ref", expanded["${MISSING}"])

	// Expand is a combinator: the caller's map is left untouched.
	try.Equal(Names{
		"A":          "1",
		"B":          "${A}",
		"C":          "${MISSING}",
		"${MISSING}": "key-ref",
	}, original)
}

func TestNames_Replace(t *testing.T) {
	try := assert.New(t)

	input := []string{"${A}", "plain"}
	got := Names{"A": "1"}.Replace(input...)
	try.Equal([]string{"1", "plain"}, got)

	// Replace returns a fresh slice: the caller's backing array is untouched.
	try.Equal([]string{"${A}", "plain"}, input)
}
