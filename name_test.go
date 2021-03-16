package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName_Value(t *testing.T) {
	const name Name = "KV_TEST_NAME"
	t.Setenv(string(name), "value")
	try := assert.New(t)
	try.Equal("value", name.Value())
}

func TestName_ValueOr(t *testing.T) {
	const name Name = "KV_TEST_NAME"
	try := assert.New(t)

	t.Run("present", func(t *testing.T) {
		t.Setenv(string(name), "value")
		try.Equal("value", name.ValueOr("fallback"))
	})

	t.Run("absent", func(t *testing.T) {
		try.Equal("fallback", Name("KV_TEST_MISSING").ValueOr("fallback"))
	})
}

func TestName_Lookup(t *testing.T) {
	const name Name = "KV_TEST_NAME"
	t.Setenv(string(name), "value")
	try := assert.New(t)

	value, ok := name.Lookup()
	try.True(ok)
	try.Equal("value", value)

	value, ok = Name("KV_TEST_MISSING").Lookup()
	try.False(ok)
	try.Equal("", value)
}
