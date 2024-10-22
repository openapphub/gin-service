package middleware

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"openapphub/pkg/cache"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

type CachedResponse struct {
	Status int
	Header gin.H
	Data   []byte
}

var (
	group singleflight.Group
)

func CacheMiddleware(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow caching for GET and POST requests
		if c.Request.Method != "GET" && c.Request.Method != "POST" {
			c.Next()
			return
		}

		// Check if we should bypass the cache
		if c.GetHeader("X-Bypass-Cache") == "true" {
			c.Next()
			return
		}

		// Generate cache key
		key := generateCacheKey(c)

		var cachedResponse *CachedResponse
		var fromCache bool

		// Use singleflight to handle concurrent requests
		resp, err, _ := group.Do(key, func() (interface{}, error) {
			// Try to get the cached response
			cr, err := getCachedResponse(c, key)
			if err == nil {
				fromCache = true
				return cr, nil
			}

			// If not in cache, create a custom ResponseWriter
			w := &responseWriter{
				ResponseWriter: c.Writer,
				body:           &bytes.Buffer{},
			}
			c.Writer = w

			// Process the request
			c.Next()

			// Create the response
			response := &CachedResponse{
				Status: w.Status(),
				Header: make(gin.H),
				Data:   w.body.Bytes(),
			}

			// Copy headers
			for k, v := range w.Header() {
				response.Header[k] = v[0]
			}

			// Cache the response if it's successful
			if w.Status() >= 200 && w.Status() < 300 {
				go cacheResponse(c.Copy(), key, response, duration) // Cache asynchronously
			}

			return response, nil
		})

		if err != nil {
			c.Next() // If there's an error, just continue without caching
			return
		}

		cachedResponse, ok := resp.(*CachedResponse)
		if !ok || cachedResponse == nil {
			// This should not happen, but if it does, we've already written the response
			return
		}

		// If the response is from cache, write it to the client
		if fromCache {
			// Set headers and write response
			for k, v := range cachedResponse.Header {
				c.Header(k, fmt.Sprint(v))
			}
			c.Header("X-From-Cache", "true")
			c.Data(cachedResponse.Status, c.Writer.Header().Get("Content-Type"), cachedResponse.Data)
		}
		// If it's not from cache, the response has already been written by c.Next()
		c.Abort() // Prevent further handlers from being called
	}
}

func generateCacheKey(c *gin.Context) string {
	// 添加版本前缀，方便未来升级时废弃旧缓存
	version := "v1:"
	if c.Request.Method == "GET" {
		return version + c.Request.URL.String()
	}
	// POST 请求可以考虑加入请求体的哈希
	return version + c.Request.URL.String()
}

func getCachedResponse(c *gin.Context, key string) (*CachedResponse, error) {
	var response CachedResponse
	data, err := cache.Get(c, key)
	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(bytes.NewReader([]byte(data)))
	err = decoder.Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func cacheResponse(c *gin.Context, key string, response *CachedResponse, duration time.Duration) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(response)
	if err != nil {
		return err
	}

	return cache.Set(c, key, buf.String(), duration)
}

type responseWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Body() []byte {
	return w.body.Bytes()
}

func InvalidateCache(c *gin.Context, key string) error {
	return cache.Del(c, key)
}

func RefreshCache(c *gin.Context, key string, duration time.Duration) error {
	// Get the current cached response
	cachedResponse, err := getCachedResponse(c, key)
	if err != nil {
		return err
	}

	// Re-cache the response with a new duration
	return cacheResponse(c, key, cachedResponse, duration)
}

func ClearCacheByPrefix(c *gin.Context, prefix string) error {
	return cache.DelByPrefix(c, prefix)
}

func GenerateCacheKey(method, path string, body []byte) string {
	if method == "GET" {
		return path
	} else if method == "POST" {
		return path + string(body)
	}
	return ""
}
