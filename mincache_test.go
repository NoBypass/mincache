package mincache

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test(t *testing.T) {
	c := New()

	t.Run("Test normal set and get", func(t *testing.T) {
		c.Set("key", "value", 0)
		value, ok := c.Get("key")
		assert.True(t, ok)
		assert.Equal(t, "value", value)
	})
	t.Run("Test expired key", func(t *testing.T) {
		c.Set("key", "value", -1)
		_, ok := c.Get("key")
		assert.False(t, ok)
	})
	t.Run("Test not yet expired key", func(t *testing.T) {
		c.Set("key", "value", 1*time.Second)
		_, ok := c.Get("key")
		time.Sleep(1 * time.Second)
		_, ok2 := c.Get("key")
		assert.True(t, ok)
		assert.False(t, ok2)
	})
	t.Run("Test Delete", func(t *testing.T) {
		c.Set("key", "value", 0)
		c.Delete("key")
		_, ok := c.Get("key")
		assert.False(t, ok)
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
			c.Set(key, key, 0)
			val, _ := c.Get(key)
			t.Log(val)
			c.Delete(key)
		}()
	}
	wg.Wait()
}
