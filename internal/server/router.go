package server

import (
	"fmt"
	"openapphub/internal/api"
	"openapphub/internal/middleware"
	"os"
	"time"

	_ "openapphub/docs" // This line is important

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Cors())
	// 使用安全中间件
	r.Use(middleware.SecureMiddleware())
	// 使用日志中间件
	r.Use(middleware.Logger())
	r.Use(middleware.RecoveryWithZap())
	// 使用限流中间件
	r.Use(middleware.RateLimiter(middleware.RateLimiterConfig{
		RateString:  "100-H",
		LimitByUser: false,
	}))
	// 使用 gzip
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// 根据认证模式选择中间件
	authMode := os.Getenv("AUTH_MODE")

	if authMode != "jwt" {
		r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	}
	r.Use(middleware.CurrentUser())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// API 路由
	apiVersion := "v1" // 可以轻松更改 API 版本
	// v1 := r.Group("/api/v1")
	v1 := r.Group(fmt.Sprintf("/api/%s", apiVersion))
	{
		// 公开路由
		v1.GET("ping", middleware.CacheMiddleware(5*time.Minute), api.Ping)
		// 缓存 ping
		v1.POST("ping", middleware.CacheMiddleware(5*time.Minute), api.Ping)
		// 用户登录
		v1.POST("user/register", api.UserRegister)
		// 用户登录
		v1.POST("user/login", middleware.CacheMiddleware(5*time.Minute), api.UserLogin)
		// 刷新用户token
		v1.POST("user/refresh", api.RefreshToken)

		// 缓存管理, 常规情况下只允许管理员执行
		v1.POST("cache/clear", api.ClearCacheByPrefix)
		v1.POST("cache/refresh", api.RefreshCache)
		v1.POST("cache/invalidate", api.InvalidateCache)

		// 需要认证的路由
		auth := v1.Group("")
		auth.Use(middleware.AuthRequired())
		{
			// User Routing
			auth.GET("user/me", api.UserMe)
			auth.DELETE("user/logout", api.UserLogout)
			auth.POST("user/logout/all", api.UserLogoutAll)
			auth.POST("user/logout/:device_id", api.UserLogoutDevice)
			auth.GET("user/devices", api.UserDevices)
		}
	}
	return r
}
