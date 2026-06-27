package kv

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SetFromFile fills in the environment variable target from a file, but only
// when target isn't already set and fileVar points at a readable file. It's
// handy for big or secret values that are nicer to drop in a file than to type
// inline. If target is already set it's left alone, and an unset fileVar just
// does nothing. Whatever gets loaded is whitespace-trimmed.
func SetFromFile(target, fileVar string) error {
	if os.Getenv(target) != "" {
		return nil
	}
	filename := os.Getenv(fileVar)
	if filename == "" {
		return nil
	}
	content, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("%w %s=%q into %s: %w", ErrFileLoad, fileVar, filename, target, err)
	}
	return os.Setenv(target, strings.TrimSpace(string(content)))
}

// LoadWithUnmarshaler reads reader and unmarshals it into a brand-new Environment.
func LoadWithUnmarshaler(reader io.Reader, unmarshaler func([]byte, any) error) (Environment, error) {
	env := Environment{}
	return env, env.LoadFromUnmarshaler(reader, unmarshaler)
}

// LoadFromJSON reads JSON from reader into a fresh Environment.
func LoadFromJSON(reader io.Reader) (Environment, error) {
	return LoadWithUnmarshaler(reader, json.Unmarshal)
}

// LoadFromYAML reads YAML from reader into a fresh Environment.
func LoadFromYAML(reader io.Reader) (Environment, error) {
	return LoadWithUnmarshaler(reader, yaml.Unmarshal)
}

// LoadFileWithUnmarshaler opens file and unmarshals what's inside into a fresh Environment.
func LoadFileWithUnmarshaler(file string, unmarshaler func([]byte, any) error) (Environment, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	return LoadWithUnmarshaler(r, unmarshaler)
}

// LoadFromJSONFile opens file and reads its JSON into a fresh Environment.
func LoadFromJSONFile(file string) (Environment, error) {
	return LoadFileWithUnmarshaler(file, json.Unmarshal)
}

// LoadFromYAMLFile opens file and reads its YAML into a fresh Environment.
func LoadFromYAMLFile(file string) (Environment, error) {
	return LoadFileWithUnmarshaler(file, yaml.Unmarshal)
}
