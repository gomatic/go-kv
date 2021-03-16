package kv

import (
	"os"
	"strings"
)

// Get hands back the value of the environment variable key, or "" if it isn't set.
func Get(key string) string { return os.Getenv(key) }

// Lookup hands back the value of the environment variable key along with whether it's set.
func Lookup(key string) (string, bool) { return os.LookupEnv(key) }

// GetTrimmed grabs the value of the environment variable key and trims off any
// leading and trailing whitespace.
func GetTrimmed(key string) string { return strings.TrimSpace(os.Getenv(key)) }

// GetOr returns the value of the environment variable key, falling back to fallback when it isn't set.
func GetOr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// First returns the value of the first environment variable in keys that's actually set, or "".
func First(keys ...string) string {
	value, _ := LookupFirst(keys...)
	return value
}

// FirstOr returns the value of the first environment variable in keys that's set,
// or fallback if none of them are.
func FirstOr(fallback string, keys ...string) string {
	if value, ok := LookupFirst(keys...); ok {
		return value
	}
	return fallback
}

// LookupFirst walks keys and returns the value of the first one that's set,
// plus whether it found anything at all.
func LookupFirst(keys ...string) (string, bool) {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			return value, true
		}
	}
	return "", false
}

// SetWithRestore sets the environment variable key to value and hands back a
// function that puts things back the way they were. If value is empty and
// allowEmpty is false, key gets unset rather than set to "".
func SetWithRestore(key, value string, allowEmpty bool) func() {
	previous, existed := os.LookupEnv(key)
	if value == "" && !allowEmpty {
		_ = os.Unsetenv(key)
	} else {
		_ = os.Setenv(key, value)
	}
	if existed {
		return func() { _ = os.Setenv(key, previous) }
	}
	return func() { _ = os.Unsetenv(key) }
}
