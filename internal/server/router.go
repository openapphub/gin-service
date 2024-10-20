package server

import (
	"openapphub/internal/api"
	"openapphub/internal/middleware"
	"os"

	_ "openapphub/docs" // This line is important

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Cors())
	// 使用日志中间件
	r.Use(middleware.Logger())
	r.Use(middleware.RecoveryWithZap())

	// 根据认证模式选择中间件
	authMode := os.Getenv("AUTH_MODE")

	if authMode != "jwt" {
		r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	}
	r.Use(middleware.CurrentUser())

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("ping", api.Ping)

		// 用户登录
		v1.POST("user/register", api.UserRegister)

		// 用户登录
		v1.POST("user/login", api.UserLogin)
		// 刷新token
		v1.POST("user/refresh", api.RefreshToken)
		// 需要登录保护的
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
