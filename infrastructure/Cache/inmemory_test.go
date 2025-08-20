package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCache(t *testing.T) {
	cache := NewInMemoryCache(5*time.Minute, 10*time.Minute)

	key := "testKey"
	value := "testValue"

	// Test Set and Get
	cache.Set(key, value, 1*time.Minute)
	retrievedValue, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, value, retrievedValue)

	// Test Delete
	cache.Delete(key)
	_, found = cache.Get(key)
	assert.False(t, found)

	// Test expiration
	cache.Set(key, value, 1*time.Millisecond)
	time.Sleep(2 * time.Millisecond)
	_, found = cache.Get(key)
	assert.False(t, found)
}
