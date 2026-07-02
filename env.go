package kv

import (
	"os"
	"strings"
)

// AllowEmpty controls whether [SetWithRestore] sets a key to "" rather than
// unsetting it when the supplied value is empty.
type AllowEmpty bool

// Value is an environment-variable value supplied by the caller — a fallback
// for the Or helpers or the value to set in [SetWithRestore].
type Value string

// Get hands back the value of the environment variable key, or "" if it isn't set.
func Get(key Name) string { return os.Getenv(string(key)) }

// Lookup hands back the value of the environment variable key along with whether it's set.
func Lookup(key Name) (string, bool) { return os.LookupEnv(string(key)) }

// GetTrimmed grabs the value of the environment variable key and trims off any
// leading and trailing whitespace.
func GetTrimmed(key Name) string { return strings.TrimSpace(os.Getenv(string(key))) }

// GetOr returns the value of the environment variable key, falling back to fallback when it isn't set.
func GetOr(key Name, fallback Value) string {
	if value, ok := os.LookupEnv(string(key)); ok {
		return value
	}
	return string(fallback)
}

// First returns the value of the first environment variable in keys that's actually set, or "".
func First(keys ...Name) string {
	value, _ := LookupFirst(keys...)
	return value
}

// FirstOr returns the value of the first environment variable in keys that's set,
// or fallback if none of them are.
func FirstOr(fallback Value, keys ...Name) string {
	if value, ok := LookupFirst(keys...); ok {
		return value
	}
	return string(fallback)
}

// LookupFirst walks keys and returns the value of the first one that's set,
// plus whether it found anything at all.
func LookupFirst(keys ...Name) (string, bool) {
	for _, key := range keys {
		if value, ok := os.LookupEnv(string(key)); ok {
			return value, true
		}
	}
	return "", false
}

// SetWithRestore sets the environment variable key to value and hands back a
// function that puts things back the way they were. If value is empty and
// isAllowEmpty is false, key gets unset rather than set to "".
func SetWithRestore(key Name, value Value, isAllowEmpty AllowEmpty) func() {
	previous, existed := os.LookupEnv(string(key))
	if value == "" && !isAllowEmpty {
		_ = os.Unsetenv(string(key))
	} else {
		_ = os.Setenv(string(key), string(value))
	}
	if existed {
		return func() { _ = os.Setenv(string(key), previous) }
	}
	return func() { _ = os.Unsetenv(string(key)) }
}
