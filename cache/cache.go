package cache

import (
	"sync"
	"time"
)

type Cache struct {
	Mu    *sync.Mutex
	Cache map[string]CacheEntry
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		Cache: make(map[string]CacheEntry),
		Mu:    &sync.Mutex{},
	}
	go c.DeleteLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Cache[key] = CacheEntry{
		createdAt: time.Now().UTC(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	val, ok := c.Cache[key]
	return val.val, ok
}

func (c *Cache) DeleteLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.Delete(time.Now().UTC(), interval)
	}
}

func (c *Cache) Delete(now time.Time, last time.Duration) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	for k, v := range c.Cache {
		if v.createdAt.Before(now.Add(-last)) {
			delete(c.Cache, k)
		}
	}
}
