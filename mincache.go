package mincache

import (
	"container/heap"
	"sync"
	"time"
)

type ItemHeap []*Item

type Item struct {
	Key        string
	Value      any
	Expiration int64
	Index      int
}

type CacheInstance struct {
	store  map[string]*Item
	queue  ItemHeap
	mu     sync.RWMutex
	signal chan struct{}
	close  chan struct{}
}

func New() *CacheInstance {
	cache := CacheInstance{
		store:  make(map[string]*Item),
		queue:  make(ItemHeap, 0),
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

func (c *CacheInstance) Close() {
	close(c.close)
}

func (c *CacheInstance) Get(key string) (value any, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.store[key]
	if !ok || (time.Now().UnixNano() > item.Expiration && item.Expiration != 0) {
		delete(c.store, key)
		return nil, false
	}
	return item.Value, true
}

func (c *CacheInstance) Set(key string, value any, duration time.Duration) {
	expiration := time.Now().Add(duration).UnixNano()
	if duration == 0 {
		expiration = 0
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.store[key]; ok {
		item.Value = value
		item.Expiration = expiration
		heap.Fix(&c.queue, item.Index)
	} else {
		item = &Item{
			Key:        key,
			Value:      value,
			Expiration: expiration,
		}
		heap.Push(&c.queue, item)
		c.store[key] = item
	}

	c.signal <- struct{}{}
}

func (c *CacheInstance) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, ok := c.store[key]; ok {
		delete(c.store, key)
		heap.Remove(&c.queue, item.Index)
	}
}

func (c *CacheInstance) cleanup() {
	for {
		c.mu.Lock()
		for len(c.queue) > 0 {
			item := heap.Pop(&c.queue).(*Item)
			if time.Now().UnixNano() > item.Expiration && item.Expiration != 0 {
				delete(c.store, item.Key)
			} else {
				heap.Push(&c.queue, item)
				break
			}
		}
		c.mu.Unlock()

		timer := time.After(time.Hour * 24)
		if len(c.queue) > 0 {
			timer = time.After(time.Until(time.Unix(0, c.queue[0].Expiration)))
		}

		select {
		case <-c.signal:
		case <-timer:
		case <-c.close:
			return
		}
	}
}
