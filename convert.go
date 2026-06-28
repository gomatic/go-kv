package kv

import "os"

// Convert turns an environment value into a T. Handy bit: standard-library
// parsers like strconv.Atoi, strconv.ParseBool, and time.ParseDuration already
// match this signature, so you can hand them straight over — no adapter needed.
type Convert[T any] func(string) (T, error)

// LookupAs looks up key and runs its value through convert. ok tells you whether
// key was set at all. err is only non-nil when key was set but its value wouldn't
// convert, and in that case value comes back as the zero value of T.
func LookupAs[T any](key Name, convert Convert[T]) (value T, ok bool, err error) {
	raw, ok := os.LookupEnv(string(key))
	if !ok {
		return value, false, nil
	}
	value, err = convert(raw)
	if err != nil {
		var zero T
		return zero, true, err
	}
	return value, true, nil
}

// GetOrAs gives you the converted value of key, or fallback if key is unset or
// its value won't convert. Reach for [LookupAs] instead when you need to tell a
// conversion failure apart from a key that was never set.
func GetOrAs[T any](key Name, convert Convert[T], fallback T) T {
	if value, ok, err := LookupAs(key, convert); ok && err == nil {
		return value
	}
	return fallback
}
