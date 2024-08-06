package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cacher defines the interface for a cache with Get and Set methods.
type Cacher interface {
	Get(k string) (any, bool)
	Set(k string, x any)
}

// Cache implements the Cacher interface and provides a caching mechanism with a default expiration time.
type Cache struct {
	cache             *cache.Cache
	defaultExpiration time.Duration
}

// New creates a new Cache instance with the specified default expiration time and cleanup interval.
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	return &Cache{
		cache: cache.New(
			defaultExpiration,
			cleanupInterval),
		defaultExpiration: defaultExpiration,
	}
}

// Get retrieves an item from the cache using the provided key.
// It returns the item and a boolean indicating whether the item was found.
func (c *Cache) Get(k string) (interface{}, bool) {
	return c.cache.Get(k)
}

// Set adds an item to the cache with the specified key and the default expiration time.
func (c *Cache) Set(k string, x any) {
	c.cache.Set(k, x, c.defaultExpiration)
}
