package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	t.Run("no such directory", func(t *testing.T) {
		envs, err := ReadDir("./some/dir")

		require.Nil(t, envs)
		require.Error(t, err)
	})

	t.Run("empty dir", func(t *testing.T) {
		dir, err := ioutil.TempDir("./testdata", "*-tempdir")
		defer os.RemoveAll(dir)

		require.NoError(t, err)
		envs, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, make(Environment), envs)
	})
}
