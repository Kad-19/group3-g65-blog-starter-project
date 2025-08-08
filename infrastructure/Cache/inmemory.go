package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Ensure inMemoryCache implements the Service interface at compile time.
var _ Service = (*inMemoryCache)(nil)

// inMemoryCache is the struct that holds the actual cache instance.
// It is kept private to enforce usage through the Service interface.
type inMemoryCache struct {
	client *cache.Cache
}

// NewInMemoryCache is the constructor function for our in-memory cache service.
// It initializes a new cache with a default expiration and cleanup interval.
//
// defaultExpiration: The default duration for which an item should be kept in the cache.
// cleanupInterval: The interval at which the cache should purge expired items.
func NewInMemoryCache(defaultExpiration, cleanupInterval time.Duration) Service {
	return &inMemoryCache{
		client: cache.New(defaultExpiration, cleanupInterval),
	}
}

// Set adds an item to the cache, replacing any existing item.
func (c *inMemoryCache) Set(key string, value interface{}, duration time.Duration) {
	c.client.Set(key, value, duration)
}

// Get retrieves an item from the cache.
// It returns the item or nil, and a boolean indicating whether the key was found.
func (c *inMemoryCache) Get(key string) (interface{}, bool) {
	return c.client.Get(key)
}

// Delete removes an item from the cache.
// Does nothing if the key is not in the cache.
func (c *inMemoryCache) Delete(key string) {
	c.client.Delete(key)
}