package api

import (
	"encoding/json"
	"fmt"
	"openapphub/internal/config"
	"openapphub/internal/model"
	"openapphub/internal/util"
	"openapphub/pkg/cache"
	"openapphub/pkg/serializer"

	"openapphub/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

// const cachePrefix = "cache_"

// Ping godoc
// @Summary Ping test
// @Description do ping
// @Tags ping
// @Accept json
// @Produce json
// @Success 200 {string} string "pong"
// @Router /ping [post]
func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Pong",
	})
}

// CurrentUser 获取当前用户
func CurrentUser(c *gin.Context) *model.User {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*model.User); ok {
			return u
		}
	}
	return nil
}

// ErrorResponse 返回错误消息
func ErrorResponse(err error) serializer.Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := config.T(fmt.Sprintf("Field.%s", e.Field()))
			tag := config.T(fmt.Sprintf("Tag.Valid.%s", e.Tag()))
			return serializer.ParamErr(
				fmt.Sprintf("%s%s", field, tag),
				err,
			)
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.ParamErr("JSON类型不匹配", err)
	}

	return serializer.ParamErr("参数错误", err)
}

// ClearCacheByPrefix godoc
// @Summary Clear cache by prefix
// @Description Clear all cached items with a specific prefix
// @Tags cache
// @Accept json
// @Produce json
// @Param prefix body string true "Cache key prefix"
// @Success 200 {object} serializer.Response "Cache cleared successfully"
// @Failure 400 {object} serializer.Response "Bad request"
// @Failure 500 {object} serializer.Response "Internal server error"
// @Router /cache/clear [post]
func ClearCacheByPrefix(c *gin.Context) {
	var input struct {
		Prefix string `json:"prefix" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, ErrorResponse(err))
		return
	}

	err := middleware.ClearCacheByPrefix(c, input.Prefix)
	if err != nil {
		if err.Error() == "no keys found matching prefix: "+input.Prefix {
			c.JSON(200, serializer.Response{
				Code: 0,
				Msg:  "No cache entries found with the given prefix",
			})
		} else {
			c.JSON(500, serializer.Err(500, "Failed to clear cache", err))
		}
		return
	}

	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Cache cleared successfully",
	})
}

// RefreshCache godoc
// @Summary Refresh cache for a specific key
// @Description Refresh the cache for a specific key with a new duration
// @Tags cache
// @Accept json
// @Produce json
// @Param refresh_info body RefreshCacheInput true "Refresh Cache Info"
// @Success 200 {object} serializer.Response "Cache refreshed successfully"
// @Failure 400 {object} serializer.Response "Bad request"
// @Failure 500 {object} serializer.Response "Internal server error"
// @Router /cache/refresh [post]
func RefreshCache(c *gin.Context) {
	var input RefreshCacheInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, ErrorResponse(err))
		return
	}

	// 使用输入参数生成缓存键
	var body []byte
	if input.Body != "" {
		body = []byte(input.Body)
	}
	key := middleware.GenerateCacheKeyFromParams(input.Method, input.Path, "", body)

	util.Log().Info("Attempting to refresh cache with key: " + key)
	exists, err := cache.Exists(c, key)
	if err != nil {
		c.JSON(500, serializer.Err(500, "Failed to check cache existence", err))
		return
	}

	if !exists {
		c.JSON(200, serializer.Response{
			Code: 0,
			Msg:  "Cache key not found",
		})
		return
	}

	err = cache.Expire(c, key, time.Duration(input.Duration)*time.Second)
	if err != nil {
		c.JSON(500, serializer.Err(500, "Failed to refresh cache", err))
		return
	}

	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Cache refreshed successfully",
	})
}

// InvalidateCache godoc
// @Summary Invalidate cache for a specific key
// @Description Remove a specific key from the cache
// @Tags cache
// @Accept json
// @Produce json
// @Param invalidate_info body InvalidateCacheInput true "Invalidate Cache Info"
// @Success 200 {object} serializer.Response "Cache invalidated successfully"
// @Failure 400 {object} serializer.Response "Bad request"
// @Failure 500 {object} serializer.Response "Internal server error"
// @Router /cache/invalidate [post]
func InvalidateCache(c *gin.Context) {
	var input InvalidateCacheInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, ErrorResponse(err))
		return
	}

	var body []byte
	if input.Body != "" {
		body = []byte(input.Body)
		util.Log().Info(fmt.Sprintf("Invalidate cache input body: %s", input.Body))
	}

	key := middleware.GenerateCacheKeyFromParams(input.Method, input.Path, "", body)
	util.Log().Info(fmt.Sprintf("Attempting to invalidate cache with key: %s", key))

	// 尝试删除缓存
	err := cache.Del(c, key)
	if err != nil {
		if err == redis.Nil {
			util.Log().Info(fmt.Sprintf("Cache key not found: %s", key))
			c.JSON(200, serializer.Response{
				Code: 0,
				Msg:  "Cache key not found",
			})
		} else {
			util.Log().Error(fmt.Sprintf("Failed to invalidate cache: %s, error: %s", key, err.Error()))
			c.JSON(500, serializer.Err(500, "Failed to invalidate cache", err))
		}
		return
	}

	util.Log().Info(fmt.Sprintf("Successfully invalidated cache key: %s", key))
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Cache invalidated successfully",
	})
}

type RefreshCacheInput struct {
	Method   string `json:"method" binding:"required,oneof=GET POST"`
	Path     string `json:"path" binding:"required"`
	Body     string `json:"body"`
	Duration int    `json:"duration" binding:"required,min=1"`
}

type InvalidateCacheInput struct {
	Method string `json:"method" binding:"required,oneof=GET POST"`
	Path   string `json:"path" binding:"required"`
	Body   string `json:"body"`
}
