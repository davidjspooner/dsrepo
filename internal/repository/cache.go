package repository

import (
	"sort"
	"sync"
	"time"

	"golang.org/x/exp/maps"
)

type cacheEntry[T any] struct {
	key      string
	value    *T
	modified time.Time
}

type RefreshFunc[T any] func(key string, cached *T, age time.Duration) (cache *T, modified bool)

type Cache[T any] struct {
	MaxCount int
	entries  map[string]cacheEntry[T]
	lock     sync.Mutex
}

func NewCacheMap[T any](maxCount int) *Cache[T] {
	return &Cache[T]{
		MaxCount: maxCount,
		entries:  make(map[string]cacheEntry[T]),
	}
}

func (c *Cache[T]) Use(key string, refreshFunc RefreshFunc[T]) {
	c.lock.Lock()
	defer c.lock.Unlock()

	cacheEntry, hit := c.entries[key]
	var age time.Duration
	if !hit {
		cacheEntry.key = key
	} else {
		age = time.Since(cacheEntry.modified)
	}
	var modified bool
	cacheEntry.value, modified = refreshFunc(key, cacheEntry.value, age)
	if cacheEntry.value == nil {
		delete(c.entries, key)
		return
	}
	if modified {
		cacheEntry.modified = time.Now()
	}
	c.entries[key] = cacheEntry
	if len(c.entries) > c.MaxCount {
		c.prune()
	}
}

func (c *Cache[T]) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.entries)
}

func (c *Cache[T]) Prune() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.prune()
}

func (c *Cache[T]) prune() {

	if len(c.entries) <= c.MaxCount {
		return
	}
	entries := maps.Values(c.entries)

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].modified.Before(entries[j].modified)
	})

	for i := 0; i < len(entries)-c.MaxCount; i++ {
		delete(c.entries, entries[i].key)
	}
}

func (c *Cache[T]) UseAll(refreshFunc RefreshFunc[T]) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var modified bool
	for key, entry := range c.entries {
		age := time.Since(entry.modified)
		entry.value, modified = refreshFunc(key, entry.value, age)
		if entry.value == nil {
			delete(c.entries, key)
			continue
		}
		if modified {
			entry.modified = time.Now()
		}
		c.entries[key] = entry
	}
	if len(c.entries) > c.MaxCount {
		c.prune()
	}
}
