package kv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Setenv("KV_TEST_GET", "value")
	try := assert.New(t)
	try.Equal("value", Get("KV_TEST_GET"))
	try.Equal("", Get("KV_TEST_MISSING"))
}

func TestLookup(t *testing.T) {
	t.Setenv("KV_TEST_LOOKUP", "value")
	try := assert.New(t)

	value, ok := Lookup("KV_TEST_LOOKUP")
	try.True(ok)
	try.Equal("value", value)

	value, ok = Lookup("KV_TEST_MISSING")
	try.False(ok)
	try.Equal("", value)
}

func TestGetTrimmed(t *testing.T) {
	t.Setenv("KV_TEST_TRIM", "  spaced  ")
	try := assert.New(t)
	try.Equal("spaced", GetTrimmed("KV_TEST_TRIM"))
}

func TestGetOr(t *testing.T) {
	t.Setenv("KV_TEST_OR", "value")
	try := assert.New(t)
	try.Equal("value", GetOr("KV_TEST_OR", "fallback"))
	try.Equal("fallback", GetOr("KV_TEST_MISSING", "fallback"))
}

func TestFirst(t *testing.T) {
	t.Setenv("KV_TEST_SECOND", "second")
	try := assert.New(t)
	try.Equal("second", First("KV_TEST_MISSING", "KV_TEST_SECOND"))
	try.Equal("", First("KV_TEST_MISSING"))
}

func TestFirstOr(t *testing.T) {
	t.Setenv("KV_TEST_SECOND", "second")
	try := assert.New(t)
	try.Equal("second", FirstOr("fallback", "KV_TEST_MISSING", "KV_TEST_SECOND"))
	try.Equal("fallback", FirstOr("fallback", "KV_TEST_MISSING"))
}

func TestLookupFirst(t *testing.T) {
	t.Setenv("KV_TEST_SECOND", "second")
	try := assert.New(t)

	value, ok := LookupFirst("KV_TEST_MISSING", "KV_TEST_SECOND")
	try.True(ok)
	try.Equal("second", value)

	value, ok = LookupFirst("KV_TEST_MISSING")
	try.False(ok)
	try.Equal("", value)
}

func TestSetWithRestore(t *testing.T) {
	const key = "KV_TEST_RESTORE"

	t.Run("existing key is restored", func(t *testing.T) {
		t.Setenv(key, "original")
		try := assert.New(t)

		restore := SetWithRestore(key, "changed", false)
		try.Equal("changed", os.Getenv(key))

		restore()
		try.Equal("original", os.Getenv(key))
	})

	t.Run("absent key is unset on restore", func(t *testing.T) {
		try := assert.New(t)
		os.Unsetenv(key)
		t.Cleanup(func() { os.Unsetenv(key) })

		restore := SetWithRestore(key, "value", false)
		try.Equal("value", os.Getenv(key))

		restore()
		_, ok := os.LookupEnv(key)
		try.False(ok)
	})

	t.Run("empty value without allowEmpty unsets", func(t *testing.T) {
		t.Setenv(key, "original")
		try := assert.New(t)

		restore := SetWithRestore(key, "", false)
		_, ok := os.LookupEnv(key)
		try.False(ok)

		restore()
		try.Equal("original", os.Getenv(key))
	})

	t.Run("empty value with allowEmpty sets empty", func(t *testing.T) {
		t.Setenv(key, "original")
		try := assert.New(t)

		restore := SetWithRestore(key, "", true)
		value, ok := os.LookupEnv(key)
		try.True(ok)
		try.Equal("", value)

		restore()
		try.Equal("original", os.Getenv(key))
	})
}
