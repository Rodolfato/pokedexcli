package pokecache

import (
	"fmt"
	"time"
)

func Hello() {
	fmt.Println("Hello from pokecache")
}

const (
	Nanosecond  time.Duration = 1
	Microsecond               = 1000 * Nanosecond
	Millisecond               = 1000 * Microsecond
	Second                    = 1000 * Millisecond
	Minute                    = 60 * Second
	Hour                      = 60 * Minute
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	CacheMap map[string]cacheEntry
	Interval time.Duration
}

func (cache *Cache) Add(entryName string, value []byte) {
	cache.CacheMap[entryName] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
	cache.reapLoop()
}

func (cache *Cache) Get(entryName string) ([]byte, bool) {
	var zeroValue []byte
	entry, ok := cache.CacheMap[entryName]
	if !ok {
		fmt.Printf("\tKey %s in cache not found\n", entryName)
		return zeroValue, false
	} else {
		fmt.Printf("\tFound key %s in cache\n", entryName)
		return entry.val, true
	}

}

func (cache *Cache) reapLoop() {
	for key, entry := range cache.CacheMap {
		duration := time.Since(entry.createdAt)
		if duration > cache.Interval {
			delete(cache.CacheMap, key)
			fmt.Printf("Deleted %s\n", key)
		}
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		CacheMap: make(map[string]cacheEntry),
		Interval: interval,
	}
	return &cache
}
