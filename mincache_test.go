package mincache

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test(t *testing.T) {
	t.Run("Test normal set and get", func(t *testing.T) {
		t.Parallel()
		c := New()
		c.Set("key", "value")
		value, ok := c.Get("key")
		assert.True(t, ok)
		assert.Equal(t, "value", value)
	})
	t.Run("Test Delete", func(t *testing.T) {
		t.Parallel()
		c := New()
		c.Set("key", "value")
		_, ok := c.Get("key")
		assert.True(t, ok)
		c.Delete("key")
		_, ok = c.Get("key")
		assert.False(t, ok)
	})
	t.Run("Test expiration duration", func(t *testing.T) {
		t.Parallel()
		c := New()
		c.Set("key", "value", ExpireIn(1*time.Second))
		time.Sleep(time.Second / 2)
		_, ok := c.Get("key")
		assert.True(t, ok)
		time.Sleep(1 * time.Second)
		_, ok = c.Get("key")
		assert.False(t, ok)
	})
	t.Run("Test expiration time", func(t *testing.T) {
		t.Parallel()
		c := New()
		c.Set("key", "value", ExpireAt(time.Now().Add(1*time.Second)))
		time.Sleep(time.Second / 2)
		_, ok := c.Get("key")
		assert.True(t, ok)
		time.Sleep(1 * time.Second)
		_, ok = c.Get("key")
		assert.False(t, ok)
	})
	t.Run("Test expiration order", func(t *testing.T) {
		t.Parallel()
		c := New()
		c.Set("key1", "value1", ExpireIn(1*time.Second))
		c.Set("key2", "value2", ExpireIn(2*time.Second))
		c.Set("key3", "value3", ExpireIn(3*time.Second))
		c.Set("key4", "value4")
		time.Sleep(time.Second / 2)
		_, ok := c.Get("key1")
		assert.True(t, ok)
		_, ok = c.Get("key2")
		assert.True(t, ok)
		_, ok = c.Get("key3")
		assert.True(t, ok)
		_, ok = c.Get("key4")
		assert.True(t, ok)

		time.Sleep(1 * time.Second)
		_, ok = c.Get("key1")
		assert.False(t, ok)
		_, ok = c.Get("key2")
		assert.True(t, ok)
		_, ok = c.Get("key3")
		assert.True(t, ok)
		_, ok = c.Get("key4")
		assert.True(t, ok)

		time.Sleep(1 * time.Second)
		_, ok = c.Get("key1")
		assert.False(t, ok)
		_, ok = c.Get("key2")
		assert.False(t, ok)
		_, ok = c.Get("key3")
		assert.True(t, ok)
		_, ok = c.Get("key4")
		assert.True(t, ok)

		time.Sleep(1 * time.Second)
		_, ok = c.Get("key1")
		assert.False(t, ok)
		_, ok = c.Get("key2")
		assert.False(t, ok)
		_, ok = c.Get("key3")
		assert.False(t, ok)
		_, ok = c.Get("key4")
		assert.True(t, ok)
	})
	t.Run("Test dureation if value was overwritten", func(t *testing.T) {
		t.Parallel()
		c := New()
		c.Set("key", "value", ExpireIn(1*time.Second))
		c.Set("key", "value2", ExpireIn(3*time.Second))
		time.Sleep(2 * time.Second)
		v, ok := c.Get("key")
		assert.True(t, ok)
		assert.Equal(t, "value2", v)
		time.Sleep(2 * time.Second)
		_, ok = c.Get("key")
		assert.False(t, ok)
	})
}

func TestSafe(t *testing.T) {
	c := NewSafe[string, int]()

	t.Run("Test normal set and get", func(t *testing.T) {
		c.Set("key", 1)
		value, ok := c.Get("key")
		assert.True(t, ok)
		assert.Equal(t, 1, value)
	})
}

func TestConcurrent(t *testing.T) {
	c := New()

	var wg sync.WaitGroup
	wg.Add(10)
	for i := range 10 {
		go func() {
			defer wg.Done()
			key := strconv.Itoa(i)
			c.Set(key, key)
			val, _ := c.Get(key)
			t.Log(val)
			c.Delete(key)
		}()
	}
	wg.Wait()
}
