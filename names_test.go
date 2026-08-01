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

// TestOrSelfNeverLosesANonEmptyValue names orSelf's claim. Expansion runs over
// every key and value, and os.Expand yields the empty string for a reference to
// an unset variable — so without this fallback, a value of "${UNSET}" would
// expand to "" and the entry would silently become empty. Falling back to the
// original keeps the reference visible instead of destroying the data, which is
// the only safe answer when the environment cannot answer the question.
//
// A genuinely empty original must stay empty: the fallback is for expansion
// COLLAPSING something, not for inventing content.
func TestOrSelfNeverLosesANonEmptyValue(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		expanded text
		original text
		want     string
		why      string
	}{
		{expanded: "value", original: "${VAR}", want: "value", why: "a successful expansion wins"},
		{expanded: "", original: "${UNSET}", want: "${UNSET}", why: "a collapsed expansion falls back rather than losing the entry"},
		{expanded: "", original: "", want: "", why: "an empty original stays empty; the fallback invents nothing"},
		{expanded: "x", original: "", want: "x", why: "an expansion of an empty original is still the expansion"},
	} {
		assert.Equal(t, tc.want, orSelf(tc.expanded, tc.original),
			"orSelf(%q, %q): %s", tc.expanded, tc.original, tc.why)
	}
}
