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

// Expand chases down ${VAR}-style references in both the keys and values of
// Names, looking them up in Names first and then the process environment, and
// keeps going until the map stops changing or it hits a fixed iteration limit.
func (names Names) Expand() Names {
	if len(names) == 0 {
		return names
	}
	response := names
	for done, limit := false, 10; !done && limit > 0; limit-- {
		done = true
		update := response
		var keys sort.StringSlice
		for key := range response {
			keys = append(keys, key)
		}
		keys.Sort()
		for _, key := range keys {
			value := response[key]
			x := os.Expand(key, update.Get)
			y := os.Expand(value, update.Get)
			done = done && key == x && y == value
			if x == "" && key != "" {
				x = key
			}
			if y == "" && value != "" {
				y = value
			}
			update[x] = y
		}
		response = update
	}
	return response
}

// Replace swaps out ${VAR}-style references in each of values using Names.
func (names Names) Replace(values ...string) []string {
	for i, value := range values {
		values[i] = os.Expand(value, names.Get)
	}
	return values
}
