// Package kv is a grab-bag of helpers for reading and poking at key/value
// environment data: an [Environment] map, lookups with fallbacks, loading from
// JSON/YAML, and scoped set/restore of the process environment.
//
// Heads up on the parameter order of the lookup helpers — it isn't the same for
// both:
//   - GetOr(key, fallback)
//   - FirstOr(fallback, keys...)
package kv

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Environment is just a map of environment keys to values. It doesn't have to
// mirror the process environment (os.Environ) — you can use it as a standalone
// map that happens to get the same lookup and load operations.
type Environment map[string]string

// New returns an Environment filled in from the current process environment.
func New() Environment { return Parse(os.Environ()) }

// Parse builds an Environment from a slice of "KEY=VALUE" entries, like the ones
// os.Environ hands you. Entries with no value are kept around with an empty string.
func Parse(environ []string) Environment {
	env := make(Environment, len(environ))
	for _, entry := range environ {
		key, value, _ := strings.Cut(entry, "=")
		env[key] = value
	}
	return env
}

// Get returns the value for key, or "" when it's not in the map.
func (e Environment) Get(key string) string { return e[key] }

// Lookup returns the value for key plus whether it was actually in there.
func (e Environment) Lookup(key string) (string, bool) {
	value, ok := e[key]
	return value, ok
}

// GetOr returns the value for key, or fallback when it's missing.
func (e Environment) GetOr(key, fallback string) string {
	if value, ok := e[key]; ok {
		return value
	}
	return fallback
}

// First returns the value of the first key that's present, or "".
func (e Environment) First(keys ...string) string {
	value, _ := e.LookupFirst(keys...)
	return value
}

// FirstOr returns the value of the first key that's present, or fallback if none of them are.
func (e Environment) FirstOr(fallback string, keys ...string) string {
	if value, ok := e.LookupFirst(keys...); ok {
		return value
	}
	return fallback
}

// LookupFirst returns the value of the first key that's present, plus whether it found one.
func (e Environment) LookupFirst(keys ...string) (string, bool) {
	for _, key := range keys {
		if value, ok := e[key]; ok {
			return value, true
		}
	}
	return "", false
}

// Set pushes every entry into the process environment and returns e so you can chain.
func (e Environment) Set() Environment {
	for key, value := range e {
		_ = os.Setenv(key, value)
	}
	return e
}

// Unset pulls every entry's key back out of the process environment.
func (e Environment) Unset() {
	for key := range e {
		_ = os.Unsetenv(key)
	}
}

// LoadFromUnmarshaler reads reader and unmarshals whatever it finds into e. Pass
// a nil unmarshaler and it falls back to json.Unmarshal; pass a nil reader and
// you get [ErrNilReader] back.
func (e Environment) LoadFromUnmarshaler(reader io.Reader, unmarshaler func([]byte, any) error) error {
	if reader == nil {
		return ErrNilReader
	}
	if unmarshaler == nil {
		unmarshaler = json.Unmarshal
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return unmarshaler(data, &e)
}

// LoadFromYAML pulls YAML out of reader and into e.
func (e Environment) LoadFromYAML(reader io.Reader) error {
	return e.LoadFromUnmarshaler(reader, yaml.Unmarshal)
}

// LoadFromJSON pulls JSON out of reader and into e.
func (e Environment) LoadFromJSON(reader io.Reader) error {
	return e.LoadFromUnmarshaler(reader, json.Unmarshal)
}
