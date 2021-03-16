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

	expanded := Names{
		"A":          "1",
		"B":          "${A}",       // this one resolves to A's value
		"C":          "${MISSING}", // nothing to resolve, so it's kept as-is
		"${MISSING}": "key-ref",    // same deal for the key — left untouched
	}.Expand()

	try.Equal("1", expanded["A"])
	try.Equal("1", expanded["B"])
	try.Equal("${MISSING}", expanded["C"])
	try.Equal("key-ref", expanded["${MISSING}"])
}

func TestNames_Replace(t *testing.T) {
	try := assert.New(t)
	got := Names{"A": "1"}.Replace("${A}", "plain")
	try.Equal([]string{"1", "plain"}, got)
}
