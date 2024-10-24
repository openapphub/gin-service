package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"openapphub/internal/util"
	"openapphub/pkg/cache"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

// 用于规范化 JSON 的辅助结构
type normalizedMap map[string]interface{}

func (m normalizedMap) MarshalJSON() ([]byte, error) {
	// 获取所有键并排序
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 按排序后的键顺序构建新的 map
	orderedMap := make([]struct {
		Key   string
		Value interface{}
	}, len(keys))

	for i, k := range keys {
		orderedMap[i].Key = k
		orderedMap[i].Value = m[k]
	}

	return json.Marshal(orderedMap)
}

// 规范化 JSON 字符串
func normalizeJSON(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return input, nil
	}

	// 解析 JSON 到 map
	var data interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	// 规范化处理
	normalized := normalizeJSONValue(data)

	// 重新序列化
	return json.Marshal(normalized)
}

// 递归处理 JSON 值
func normalizeJSONValue(v interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		normalized := make(normalizedMap)
		for k, val := range v {
			normalized[k] = normalizeJSONValue(val)
		}
		return normalized
	case []interface{}:
		normalized := make([]interface{}, len(v))
		for i, val := range v {
			normalized[i] = normalizeJSONValue(val)
		}
		return normalized
	default:
		return v
	}
}

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
		key := GenerateCacheKey(c)

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
	util.Log().Info("InvalidateCache: %s", key)
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

func GenerateCacheKey(c *gin.Context) string {
	method := c.Request.Method
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	body := getRequestBody(c)

	key := generateCacheKeyInternal(method, path, query, body)
	util.Log().Info(fmt.Sprintf("Generated cache key for request - Method: %s, Path: %s, Key: %s", method, path, key))
	return key
}

func GenerateCacheKeyFromParams(method, path string, query string, body []byte) string {
	key := generateCacheKeyInternal(method, path, query, body)
	util.Log().Info(fmt.Sprintf("Generated cache key from params - Method: %s, Path: %s, Key: %s", method, path, key))
	return key
}

func generateCacheKeyInternal(method, path, query string, body []byte) string {
	version := "v1:"
	key := version + path

	if method == "GET" {
		if query != "" {
			key = key + "?" + query
		}
	} else if method == "POST" {
		if len(body) > 0 {
			// 规范化 JSON
			normalizedBody, err := normalizeJSON(body)
			if err != nil {
				util.Log().Error(fmt.Sprintf("Failed to normalize JSON body: %s", err.Error()))
				normalizedBody = body // 如果规范化失败，使用原始 body
			}

			hash := sha256.Sum256(normalizedBody)
			key = key + ":" + hex.EncodeToString(hash[:])
			util.Log().Info(fmt.Sprintf("Normalized body: %s", string(normalizedBody)))
			util.Log().Info(fmt.Sprintf("POST request body hash: %s", hex.EncodeToString(hash[:])))
		}
	}

	util.Log().Info(fmt.Sprintf("Final cache key: %s", key))
	return key
}

func getRequestBody(c *gin.Context) []byte {
	if c.Request.Method != "POST" {
		return nil
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		util.Log().Error("Failed to read request body: " + err.Error())
		return nil
	}
	// Restore the body for later use
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return body
}
