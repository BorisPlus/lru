package lru_test

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/BorisPlus/lru"
)

// go test -v list.go list_stringer.go cache.go cache_stringer.go cache_test_data.go  cache_test.go > cache_test.txt

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := lru.NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("small-cache test", func(t *testing.T) {
		zeroCache := lru.NewCache(0)
		_, okZero := zeroCache.Get("0")
		require.False(t, okZero)
		okZero = zeroCache.Set("0", "zero")
		require.False(t, okZero)

		cache := lru.NewCache(2)

		okOne := cache.Set("1", "one")
		require.False(t, okOne)

		okFirst := cache.Set("1", "first")
		require.True(t, okFirst)

		_, okOne = cache.Get("1")
		require.True(t, okOne)

		okSecond := cache.Set("2", "second")
		require.False(t, okSecond)

		okThird := cache.Set("3", "third")
		require.False(t, okThird)

		_, okOne = cache.Get("1")
		require.False(t, okOne)
	})

	t.Run("clear", func(t *testing.T) {
		cache := lru.NewCache(3)
		cache.Set("1", "first")
		cache.Set("2", "second")
		cache.Clear()

		_, okOne := cache.Get("1")
		require.False(t, okOne)
	})

	t.Run("simple", func(t *testing.T) {
		c := lru.NewCache(5)

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

	t.Run("purge logic", func(t *testing.T) {
		// Write me
	})
}

func TestGoroutinedCache(_ *testing.T) {

	c := lru.NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(lru.Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(lru.Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
