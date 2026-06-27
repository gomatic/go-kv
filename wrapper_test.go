package kv

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapCalls(t *testing.T) {
	const key = "KV_TEST_WRAP"
	_ = os.Unsetenv(key)
	t.Cleanup(func() { _ = os.Unsetenv(key) })
	try := assert.New(t)

	var observed string
	err := WrapCalls(Environment{key: "wrapped"}, func() error {
		observed = os.Getenv(key)
		return nil
	})
	try.NoError(err)
	try.Equal("wrapped", observed)

	_, ok := os.LookupEnv(key)
	try.False(ok, "environment is restored after the calls")
}

func TestWrapper_Call(t *testing.T) {
	try := assert.New(t)

	t.Run("nil functions are skipped", func(_ *testing.T) {
		try.NoError(Wrapper{}.Call(nil, nil))
	})

	t.Run("stops at the first error", func(_ *testing.T) {
		ran := false
		err := Wrapper{}.Call(
			func() error { return errTest },
			func() error { ran = true; return nil },
		)
		try.ErrorIs(err, errTest)
		try.False(ran, "functions after the first error do not run")
	})
}

func TestNewWrapper(t *testing.T) {
	const key = "KV_TEST_WRAPPER"
	t.Setenv(key, "original")
	try := assert.New(t)

	w := NewWrapper(Environment{key: "modified"})
	try.Equal("modified", os.Getenv(key))

	w.Restore()
	try.Equal("original", os.Getenv(key))
}

func ExampleWrapCalls() {
	_ = WrapCalls(Environment{"KV_EXAMPLE": "scoped"}, func() error {
		fmt.Println(Get("KV_EXAMPLE"))
		return nil
	})
	fmt.Printf("%q\n", Get("KV_EXAMPLE"))
	// Output:
	// scoped
	// ""
}
