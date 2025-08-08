package middleware

import (
	"bytes"
	"g3-g65-bsp/infrastructure/cache"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// responseCache is the struct that will be stored in the cache.
// It holds all the necessary information to reconstruct an HTTP response.
type responseCache struct {
	Status int
	Header http.Header
	Body   []byte
}

// cachedWriter is a custom gin.ResponseWriter that captures the response body and status.
type cachedWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

// Write captures the response body.
func (w *cachedWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteHeader captures the response status code.
func (w *cachedWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// CachePage returns a Gin middleware that implements cache-aside logic.
// It depends on the cache.Service interface, making it decoupled and testable.
//
// service: The cache service that fulfills the cache.Service interface.
// ttl: The time-to-live for a cache entry.
func CachePage(service cache.Service, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// We only cache GET requests.
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Use the request URL as the cache key.
		key := c.Request.URL.String()

		// --- 1. CHECK CACHE (Cache Hit) ---
		// Attempt to retrieve the response from the cache.
		cachedData, found := service.Get(key)
		if found {
			// If found, we need to type-assert it back to our responseCache struct.
			if response, ok := cachedData.(responseCache); ok {
				// Replay the cached response.
				c.Writer.WriteHeader(response.Status)
				for h, v := range response.Header {
					c.Writer.Header()[h] = v
				}
				c.Writer.Write(response.Body)
				c.Abort() // Stop processing further handlers.
				return
			}
		}

		// --- 2. PROCESS REQUEST (Cache Miss) ---
		// If the data is not in the cache, we proceed with the request.
		// We replace the original ResponseWriter with our custom cachedWriter to capture the response.
		writer := &cachedWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer([]byte{}),
			status:         http.StatusOK, // Default to 200 OK
		}
		c.Writer = writer

		c.Next() // Process the request by calling the actual handler.

		// --- 3. CACHE RESPONSE ---
		// After the handler has run, we cache the response if it was successful (200 OK).
		if writer.status == http.StatusOK {
			responseToCache := responseCache{
				Status: writer.status,
				Header: writer.Header(),
				Body:   writer.body.Bytes(),
			}
			service.Set(key, responseToCache, ttl)
		}
	}
}

// RevalidateCache returns a Gin middleware that deletes the cached value for the current request URL.
// Use this on endpoints that modify data and need to invalidate the cache for the affected resource.
// service: The cache service that fulfills the cache.Service interface.
func RevalidateCache(service cache.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.URL.String()
		service.Delete(key)
		c.Next()
	}
}
