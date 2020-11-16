package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	t.Run("success code", func(t *testing.T) {
		env := make(Environment)
		code := RunCmd([]string{"ls", "-l"}, env)

		require.Equal(t, 0, code)
	})

	t.Run("error code", func(t *testing.T) {
		env := make(Environment)
		code := RunCmd([]string{"ls", "foo"}, env)

		require.NotEqual(t, 0, code)
	})
}
