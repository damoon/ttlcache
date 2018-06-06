package ttlcache

import (
	"sync"
	"time"
)

// Cache is a synchronised map of items that auto-expire once stale
type Cache struct {
	mutex     sync.RWMutex
	ttl       time.Duration
	items     map[K]*Item
	zeroValue V
}

// K is the type of the keys
type K interface{}

// V is the type of the values
type V interface{}

// Set is a thread-safe way to add new items to the map
func (cache *Cache) Set(key K, data V) {
	cache.mutex.Lock()
	item := &Item{data: data}
	item.touch(cache.ttl)
	cache.items[key] = item
	cache.mutex.Unlock()
}

// Get is a thread-safe way to lookup items
// Every lookup, also touches the item, hence extending it's life
func (cache *Cache) Get(key K) (data V, found bool) {
	cache.mutex.Lock()
	item, exists := cache.items[key]
	if !exists || item.expired() {
		data = cache.zeroValue
		found = false
	} else {
		item.touch(cache.ttl)
		data = item.data
		found = true
	}
	cache.mutex.Unlock()
	return
}

// Count returns the number of items in the cache
// (helpful for tracking memory leaks)
func (cache *Cache) Count() int {
	cache.mutex.RLock()
	count := len(cache.items)
	cache.mutex.RUnlock()
	return count
}

func (cache *Cache) cleanup() {
	cache.mutex.Lock()
	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
	cache.mutex.Unlock()
}

// StartCleanupTimer creates a goroutine to cleanup old entries in the background
func (cache *Cache) StartCleanupTimer(interval time.Duration) *time.Ticker {
	ticker := time.NewTicker(interval)
	go (func() {
		for {
			select {
			case <-ticker.C:
				cache.cleanup()
			}
		}
	})()
	return ticker
}

// NewCache is a helper to create instance of the Cache struct
func NewCache(duration time.Duration, z V) *Cache {
	cache := &Cache{
		ttl:       duration,
		items:     map[K]*Item{},
		zeroValue: z,
	}
	return cache
}
