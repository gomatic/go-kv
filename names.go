package kv

import (
	"maps"
	"os"
	"slices"
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
// It never mutates the receiver — each pass works on a fresh clone, so callers
// keep their original map untouched and get the expanded result back.
func (names Names) Expand() Names {
	if len(names) == 0 {
		return names
	}
	expanded := maps.Clone(names)
	for limit, changed := expandLimit, true; limit > 0 && changed; limit-- {
		expanded, changed = expanded.expandPass()
	}
	return expanded
}

// expandPass performs one expansion pass over a sorted snapshot of the keys,
// resolving ${VAR} references via the map (then the process environment). It
// writes into a fresh map and returns it alongside whether anything changed.
func (names Names) expandPass() (Names, bool) {
	changed := false
	next := make(Names, len(names))
	for _, key := range slices.Sorted(maps.Keys(names)) {
		value := names[key]
		x := os.Expand(key, names.Get)
		y := os.Expand(value, names.Get)
		changed = changed || key != x || value != y
		next[orSelf(x, key)] = orSelf(y, value)
	}
	return next, changed
}

// orSelf returns expanded, falling back to original when expansion collapsed a
// non-empty original to the empty string (so a key/value is never lost).
func orSelf(expanded, original string) string {
	if expanded == "" && original != "" {
		return original
	}
	return expanded
}

// Replace swaps out ${VAR}-style references in each of values using Names. It
// returns a fresh slice and never mutates the caller's backing array.
func (names Names) Replace(values ...string) []string {
	replaced := slices.Clone(values)
	for i, value := range replaced {
		replaced[i] = os.Expand(value, names.Get)
	}
	return replaced
}
