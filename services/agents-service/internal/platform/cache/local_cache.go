package cache

import (
	"sync"
	"time"
)

type entry struct {
	serverID uint
	expiry   time.Time
}

type LocalCache struct {
	data    map[string]entry
	mu      sync.RWMutex
	ttl     time.Duration
	maxSize int
}

func NewLocalCache(ttl time.Duration, maxSize int) *LocalCache {
	c := &LocalCache{
		data:    make(map[string]entry),
		ttl:     ttl,
		maxSize: maxSize,
	}

	// background cleanup (optional but good)
	go c.cleanup()

	return c
}

func (c *LocalCache) Get(key string) (uint, bool) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return 0, false
	}

	// check expiry
	if time.Now().After(e.expiry) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return 0, false
	}

	return e.serverID, true
}

func (c *LocalCache) Set(key string, serverID uint) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// simple eviction if full
	if len(c.data) >= c.maxSize {
		// delete random (cheap) OR oldest (better but more code)
		for k := range c.data {
			delete(c.data, k)
			break
		}
	}

	c.data[key] = entry{
		serverID: serverID,
		expiry:   time.Now().Add(c.ttl),
	}
}

func (c *LocalCache) Delete(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

// cleanup expired entries every minute
func (c *LocalCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		c.mu.Lock()
		for k, v := range c.data {
			if now.After(v.expiry) {
				delete(c.data, k)
			}
		}
		c.mu.Unlock()
	}
}
