package kv

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLookupAs(t *testing.T) {
	const key = "KV_TEST_AS"
	try := assert.New(t)

	t.Run("unset key", func(_ *testing.T) {
		value, ok, err := LookupAs(key, strconv.Atoi)
		try.False(ok)
		try.NoError(err)
		try.Equal(0, value)
	})

	t.Run("converts a set value", func(t *testing.T) {
		t.Setenv(key, "8080")
		value, ok, err := LookupAs(key, strconv.Atoi)
		try.True(ok)
		try.NoError(err)
		try.Equal(8080, value)
	})

	t.Run("reports a conversion error", func(t *testing.T) {
		t.Setenv(key, "not-a-number")
		value, ok, err := LookupAs(key, strconv.Atoi)
		try.True(ok)
		try.Error(err)
		try.Equal(0, value)
	})
}

func TestGetOrAs(t *testing.T) {
	const key = "KV_TEST_AS"
	try := assert.New(t)

	t.Run("fallback when unset", func(_ *testing.T) {
		try.Equal(8080, GetOrAs(key, strconv.Atoi, 8080))
	})

	t.Run("converted value when set", func(t *testing.T) {
		t.Setenv(key, "true")
		try.True(GetOrAs(key, strconv.ParseBool, false))
	})

	t.Run("fallback on conversion error", func(t *testing.T) {
		t.Setenv(key, "nope")
		try.Equal(time.Second, GetOrAs(key, time.ParseDuration, time.Second))
	})
}
