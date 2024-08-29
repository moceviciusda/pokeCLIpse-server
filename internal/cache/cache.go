package cache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mutex   *sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{mutex: &sync.RWMutex{}, entries: make(map[string]cacheEntry)}

	go c.reapLoop(interval)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = cacheEntry{time.Now(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, ok := c.entries[key]

	return entry.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for now := range ticker.C {
		c.mutex.Lock()

		for key, value := range c.entries {
			if value.createdAt.Add(interval).Before(now) {
				delete(c.entries, key)
			}
		}

		c.mutex.Unlock()
	}
}
