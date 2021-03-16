package kv

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetFromFile(t *testing.T) {
	const (
		target  = "KV_TEST_FROM_FILE"
		fileVar = "KV_TEST_FROM_FILE_PATH"
	)

	t.Run("already-set target is untouched", func(t *testing.T) {
		t.Setenv(target, "existing")
		t.Setenv(fileVar, "testdata/successfully-reading-from-file.txt")
		try := assert.New(t)
		try.NoError(SetFromFile(target, fileVar))
		try.Equal("existing", os.Getenv(target))
	})

	t.Run("unset fileVar is a no-op", func(t *testing.T) {
		os.Unsetenv(target)
		os.Unsetenv(fileVar)
		t.Cleanup(func() { os.Unsetenv(target) })
		try := assert.New(t)
		try.NoError(SetFromFile(target, fileVar))
		_, ok := os.LookupEnv(target)
		try.False(ok)
	})

	t.Run("loads trimmed file content", func(t *testing.T) {
		os.Unsetenv(target)
		t.Cleanup(func() { os.Unsetenv(target) })
		t.Setenv(fileVar, "testdata/whitespace.txt")
		try := assert.New(t)
		try.NoError(SetFromFile(target, fileVar))
		try.Equal("cats", os.Getenv(target))
	})

	t.Run("missing file returns ErrFileLoad", func(t *testing.T) {
		os.Unsetenv(target)
		t.Cleanup(func() { os.Unsetenv(target) })
		t.Setenv(fileVar, "testdata/does-not-exist.txt")
		try := assert.New(t)
		try.ErrorIs(SetFromFile(target, fileVar), ErrFileLoad)
	})
}

func TestLoadWithUnmarshaler(t *testing.T) {
	try := assert.New(t)
	got, err := LoadWithUnmarshaler(nil, json.Unmarshal)
	try.ErrorIs(err, ErrNilReader)
	try.Equal(Environment{}, got)
}

func TestLoadFromJSON(t *testing.T) {
	try := assert.New(t)
	got, err := LoadFromJSON(strings.NewReader(`{"TEST_KEY":"test-value"}`))
	try.NoError(err)
	try.Equal(Environment{"TEST_KEY": "test-value"}, got)
}

func TestLoadFromYAML(t *testing.T) {
	try := assert.New(t)
	got, err := LoadFromYAML(strings.NewReader("TEST_KEY: test-value\n"))
	try.NoError(err)
	try.Equal(Environment{"TEST_KEY": "test-value"}, got)
}

func TestLoadFileWithUnmarshaler(t *testing.T) {
	try := assert.New(t)
	_, err := LoadFileWithUnmarshaler("testdata/does-not-exist.json", json.Unmarshal)
	try.Error(err)
}

func TestLoadFromJSONFile(t *testing.T) {
	try := assert.New(t)
	got, err := LoadFromJSONFile("testdata/happy.json")
	try.NoError(err)
	try.Equal(Environment{"TEST_KEY": "test-value"}, got)
}

func TestLoadFromYAMLFile(t *testing.T) {
	try := assert.New(t)
	got, err := LoadFromYAMLFile("testdata/happy.yaml")
	try.NoError(err)
	try.Equal(Environment{"TEST_KEY": "test-value"}, got)
}
