package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

type Cache struct {
	mu       sync.RWMutex
	cache    map[string]cacheEntry
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		cache:    make(map[string]cacheEntry),
		interval: interval,
	}

	go newCache.ReapLoop()

	return newCache
}

func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		val:       value,
		createdAt: time.Now(),
	}

}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, found := c.cache[key]
	if !found {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) ReapLoop() {
	tick := time.NewTicker(c.interval)
	defer tick.Stop()

	for range tick.C {
		c.mu.Lock()
		for key, entry := range c.cache {
			now := time.Now()
			if c.interval < now.Sub(entry.createdAt) {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}
