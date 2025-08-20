package middleware

import (
	"g3-g65-bsp/infrastructure/cache"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCachePage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cacheService := cache.NewInMemoryCache(5*time.Minute, 10*time.Minute)
	ttl := 1 * time.Minute

	router := gin.New()
	router.Use(CachePage(cacheService, ttl))
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "this is a test")
	})

	// First request, should be a cache miss
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "this is a test", w1.Body.String())

	// Second request, should be a cache hit
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "this is a test", w2.Body.String())
}

func TestRevalidateCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cacheService := cache.NewInMemoryCache(5*time.Minute, 10*time.Minute)
	key := "/test"
	cacheService.Set(key, responseCache{Status: 200, Body: []byte("cached")}, 1*time.Minute)

	router := gin.New()
	router.Use(RevalidateCache(cacheService))
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	_, found := cacheService.Get(key)
	assert.False(t, found)
}
