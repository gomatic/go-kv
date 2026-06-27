package kv

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test error")

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errTest }

var _ io.Reader = errReader{}

func TestNew(t *testing.T) {
	t.Setenv("KV_TEST_NEW", "value")
	try := assert.New(t)
	try.Equal("value", New()["KV_TEST_NEW"])
}

func TestParse(t *testing.T) {
	try := assert.New(t)
	got := Parse([]string{"A=1", "B=two=2", "C="})
	try.Equal(Environment{"A": "1", "B": "two=2", "C": ""}, got)
}

func TestEnvironment_Get(t *testing.T) {
	e := Environment{"A": "1"}
	try := assert.New(t)
	try.Equal("1", e.Get("A"))
	try.Equal("", e.Get("missing"))
}

func TestEnvironment_Lookup(t *testing.T) {
	e := Environment{"A": "1"}
	try := assert.New(t)

	value, ok := e.Lookup("A")
	try.True(ok)
	try.Equal("1", value)

	value, ok = e.Lookup("missing")
	try.False(ok)
	try.Equal("", value)
}

func TestEnvironment_GetOr(t *testing.T) {
	e := Environment{"A": "1"}
	try := assert.New(t)
	try.Equal("1", e.GetOr("A", "fallback"))
	try.Equal("fallback", e.GetOr("missing", "fallback"))
}

func TestEnvironment_First(t *testing.T) {
	e := Environment{"B": "2"}
	try := assert.New(t)
	try.Equal("2", e.First("A", "B"))
	try.Equal("", e.First("A"))
}

func TestEnvironment_FirstOr(t *testing.T) {
	e := Environment{"B": "2"}
	try := assert.New(t)
	try.Equal("2", e.FirstOr("fallback", "A", "B"))
	try.Equal("fallback", e.FirstOr("fallback", "A"))
}

func TestEnvironment_LookupFirst(t *testing.T) {
	e := Environment{"B": "2"}
	try := assert.New(t)

	value, ok := e.LookupFirst("A", "B")
	try.True(ok)
	try.Equal("2", value)

	value, ok = e.LookupFirst("A")
	try.False(ok)
	try.Equal("", value)
}

func TestEnvironment_Set(t *testing.T) {
	const key = "KV_TEST_SET"
	_ = os.Unsetenv(key)
	t.Cleanup(func() { _ = os.Unsetenv(key) })
	try := assert.New(t)

	e := Environment{key: "value"}
	try.Equal(e, e.Set())
	try.Equal("value", os.Getenv(key))
}

func TestEnvironment_Unset(t *testing.T) {
	const key = "KV_TEST_UNSET"
	t.Setenv(key, "value")
	try := assert.New(t)

	Environment{key: "value"}.Unset()
	_, ok := os.LookupEnv(key)
	try.False(ok)
}

func TestEnvironment_LoadFromUnmarshaler(t *testing.T) {
	try := assert.New(t)

	t.Run("nil reader", func(_ *testing.T) {
		e := Environment{}
		err := e.LoadFromUnmarshaler(nil, nil)
		try.ErrorIs(err, ErrNilReader)
	})

	t.Run("read error", func(_ *testing.T) {
		e := Environment{}
		err := e.LoadFromUnmarshaler(errReader{}, nil)
		try.ErrorIs(err, errTest)
	})

	t.Run("default unmarshaler is JSON", func(_ *testing.T) {
		e := Environment{}
		err := e.LoadFromUnmarshaler(strings.NewReader(`{"A":"1"}`), nil)
		try.NoError(err)
		try.Equal(Environment{"A": "1"}, e)
	})
}

func TestEnvironment_LoadJSONInto(t *testing.T) {
	try := assert.New(t)
	e := Environment{}
	try.NoError(e.LoadJSONInto(strings.NewReader(`{"A":"1"}`)))
	try.Equal(Environment{"A": "1"}, e)
}

func TestEnvironment_LoadYAMLInto(t *testing.T) {
	try := assert.New(t)
	e := Environment{}
	try.NoError(e.LoadYAMLInto(strings.NewReader("A: \"1\"\n")))
	try.Equal(Environment{"A": "1"}, e)
}
