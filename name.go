package kv

// Name is a typed environment-variable name that comes with its own value accessors.
type Name string

// Value returns the environment variable's value, or "" if it isn't set.
func (n Name) Value() string { return Get(n) }

// ValueOr returns the environment variable's value, or fallback if it isn't set.
func (n Name) ValueOr(fallback string) string { return GetOr(n, Value(fallback)) }

// Lookup returns the environment variable's value along with whether it's set.
func (n Name) Lookup() (string, bool) { return Lookup(n) }
