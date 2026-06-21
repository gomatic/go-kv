package kv

import (
	"os"
	"sort"
)

// Names maps environment-variable names to override values, and quietly falls
// back to the process environment for any name it doesn't define.
type Names map[string]string

// Get returns the override value for name, or whatever the process environment
// has if there's no non-empty override.
func (names Names) Get(name string) string {
	if names == nil || name == "" {
		return ""
	}
	if value := names[name]; value != "" {
		return value
	}
	return os.Getenv(name)
}

// expandLimit caps how many times Expand re-scans the map before giving up on
// reaching a fixed point, bounding the cost of self-referential definitions.
const expandLimit = 10

// Expand chases down ${VAR}-style references in both the keys and values of
// Names, looking them up in Names first and then the process environment, and
// keeps going until the map stops changing or it hits a fixed iteration limit.
func (names Names) Expand() Names {
	if len(names) == 0 {
		return names
	}
	for limit := expandLimit; limit > 0 && names.expandPass(); limit-- {
	}
	return names
}

// expandPass performs one expansion pass over a sorted snapshot of the keys,
// resolving ${VAR} references via the map (then the process environment). It
// mutates the map in place and reports whether anything changed this pass.
func (names Names) expandPass() bool {
	changed := false
	for _, key := range names.sortedKeys() {
		value := names[key]
		x := os.Expand(key, names.Get)
		y := os.Expand(value, names.Get)
		changed = changed || key != x || value != y
		names[orSelf(x, key)] = orSelf(y, value)
	}
	return changed
}

// sortedKeys returns the map's keys in deterministic (ascending) order so each
// pass visits them the same way regardless of map iteration order.
func (names Names) sortedKeys() sort.StringSlice {
	keys := make(sort.StringSlice, 0, len(names))
	for key := range names {
		keys = append(keys, key)
	}
	keys.Sort()
	return keys
}

// orSelf returns expanded, falling back to original when expansion collapsed a
// non-empty original to the empty string (so a key/value is never lost).
func orSelf(expanded, original string) string {
	if expanded == "" && original != "" {
		return original
	}
	return expanded
}

// Replace swaps out ${VAR}-style references in each of values using Names.
func (names Names) Replace(values ...string) []string {
	for i, value := range values {
		values[i] = os.Expand(value, names.Get)
	}
	return values
}
