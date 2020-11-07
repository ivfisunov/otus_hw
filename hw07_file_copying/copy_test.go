package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Place your code here
	t.Run("testing offset more than file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/temp.txt", 100500, 0)

		require.Error(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("error opening file", func(t *testing.T){
		err := Copy("dummy", "testdata/temp.txt", 0, 0)

		require.Error(t, err, "error opening file")
	})

	t.Run("copying 10 bytes", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/out_offset0_limit10.txt", 0, 10)

		require.NoError(t, err, nil)
		f, err := os.Open("testdata/out_offset0_limit10.txt")
		defer f.Close()

		require.NoError(t, err, nil)
		stat, err := f.Stat()
		require.NoError(t, err, nil)
		require.Equal(t, int64(10), stat.Size(), nil)

	})
}
