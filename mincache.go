package mincache

import (
	"container/heap"
	"sync"
	"time"
)

type itemHeap []*cacheItem

type cacheItem struct {
	Key        any
	Value      any
	Expiration *options
	Index      int
}

type Cache struct {
	store  map[any]*cacheItem
	queue  itemHeap
	mu     sync.RWMutex
	signal chan struct{}
	close  chan struct{}
}

func New() *Cache {
	cache := Cache{
		store:  make(map[any]*cacheItem),
		queue:  make(itemHeap, 0),
		signal: make(chan struct{}),
	}
	heap.Init(&cache.queue)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		cache.cleanup()
	}()
	wg.Wait()
	return &cache
}

func (c *Cache) Close() {
	close(c.close)
}

func (c *Cache) Get(key any) (value any, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.store[key]
	if !ok {
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) Set(key any, value any, opts ...Option) {
	o := apply(opts)
	c.mu.Lock()

	if item, ok := c.store[key]; ok {
		item.Value = value
		item.Expiration = o
		if o.expires {
			heap.Fix(&c.queue, item.Index)
		} else {
			heap.Remove(&c.queue, item.Index)
		}
	} else {
		item = &cacheItem{
			Key:        key,
			Value:      value,
			Expiration: o,
		}
		c.store[key] = item
		if o.expires {
			heap.Push(&c.queue, item)
		}
	}

	c.mu.Unlock()

	if o.expires {
		c.signal <- struct{}{}
	}
}

func (c *Cache) Delete(key any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, ok := c.store[key]; ok {
		delete(c.store, key)
		if item.Expiration.expires {
			heap.Remove(&c.queue, item.Index)
		}
	}
}

func (c *Cache) cleanup() {
	for {
		c.mu.Lock()
		for len(c.queue) > 0 {
			item := heap.Pop(&c.queue).(*cacheItem)
			if item.Expiration.expired() {
				delete(c.store, item.Key)
			} else {
				heap.Push(&c.queue, item)
				break
			}
		}
		c.mu.Unlock()

		timer := time.After(time.Hour * 24)
		if len(c.queue) > 0 {
			timer = time.After(time.Until(c.queue[0].Expiration.expiration()))
		}

		select {
		case <-c.signal:
		case <-timer:
		case <-c.close:
			return
		}
	}
}
