package cache

import "time"

// Cache is the interface for a generic cache service.
// It allows for easy swapping of caching implementations (e.g., in-memory, Redis).
type Service interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}