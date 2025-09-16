package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]cacheEntry),
	}
	go cache.readLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, ok
	}

	valCopy := make([]byte, len(entry.val))
	copy(valCopy, entry.val)
	return valCopy, ok
}

func (c *Cache) readLoop(interval time.Duration) {
	for tick := range time.Tick(interval) {
		c.mu.Lock()
		for k, v := range c.entries {
			// Reap the entry
			if tick.Sub(v.createdAt) > interval {
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}
