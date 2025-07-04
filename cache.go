package cache

import (
	"slices"
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	ttl  time.Duration
	mu   sync.RWMutex
	data map[K]EntryWithTimeout[V]

	maxSize           int
	chronologicalKeys []K
}

type EntryWithTimeout[V any] struct {
	value   V
	expires time.Time
}

func New[K comparable, V any](maxSize int, ttl time.Duration) Cache[K, V] {
	return Cache[K, V]{
		ttl:               ttl,
		data:              make(map[K]EntryWithTimeout[V]),
		maxSize:           maxSize,
		chronologicalKeys: make([]K, 0, maxSize),
	}
}

func (c *Cache[K, V]) Read(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var zeroV V
	e, ok := c.data[key]

	switch {
	case !ok:
		return zeroV, false
	case e.expires.Before(time.Now()):
		c.deleteKeyValue(key)
		return zeroV, false
	default:
		return e.value, true
	}
}

func (c *Cache[K, V]) Upsert(key K, val V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, alreadyPresent := c.data[key]
	switch {
	case alreadyPresent:
		c.deleteKeyValue(key)
	case len(c.data) == c.maxSize:
		c.deleteKeyValue(c.chronologicalKeys[0])
	}
	c.addKeyValue(key, val)
	return nil
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deleteKeyValue(key)
}

// addKeyValue inserts a key and its value into the cache.
func (c *Cache[K, V]) addKeyValue(key K, value V) {
	c.data[key] = EntryWithTimeout[V]{
		value:   value,
		expires: time.Now().Add(c.ttl),
	}
	c.chronologicalKeys = append(c.chronologicalKeys, key)
}

// deleteKeyValue removes a key and its associated value from the cache.
func (c *Cache[K, V]) deleteKeyValue(key K) {
	c.chronologicalKeys = slices.DeleteFunc(
		c.chronologicalKeys,
		func(k K) bool { return k == key })
	delete(c.data, key)
}
