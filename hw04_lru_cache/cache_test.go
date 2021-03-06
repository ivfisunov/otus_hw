package hw04_lru_cache //nolint:golint,stylecheck

import (
	// "fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("pop first item", func(t *testing.T) {
		// Write me
		c := NewCache(3)
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ccc", 333)
		c.Set("ddd", 400)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("pop rarely used etem", func(t *testing.T) {
		c := NewCache(5)
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)
		c.Set("eee", 555)
		val, ok := c.Get("eee")
		require.True(t, ok)
		require.Equal(t, 555, val);
		c.Get("bbb")
		c.Get("bbb")
		c.Set("bbb", 222)
		c.Get("bbb")
		c.Set("ccc", 333)
		c.Get("ccc")
		c.Get("ccc")
		c.Get("ccc")
		c.Get("ddd")
		c.Set("ddd", 404)
		c.Set("ddd", 444)
		c.Get("ddd")
		c.Set("NEW", 999)
		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 222, val);

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove if task with asterisk completed

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
