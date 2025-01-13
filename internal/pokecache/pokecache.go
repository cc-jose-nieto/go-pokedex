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
	mutex   sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{entries: make(map[string]cacheEntry)}
	go cache.readLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = cacheEntry{time.Now().UTC(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entries, ok := c.entries[key]
	return entries.val, ok
}

func (c *Cache) readLoop(interval time.Duration) {
	tick := time.NewTicker(interval)
	for range tick.C {
		c.mutex.Lock()
		c.mutex.Unlock()
		toDeleteBefore := time.Now().UTC().Add(-interval)
		for key, entry := range c.entries {
			if entry.createdAt.Before(toDeleteBefore) {
				delete(c.entries, key)
			}
		}
	}
}
