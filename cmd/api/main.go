package main

import (
	"openapphub/internal/config"
	"openapphub/internal/middleware"
	"openapphub/internal/server"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title           openapphub API
// @version         1.0
// @description     This is a sample server for openapphub.
// @host            localhost:3000
// @BasePath        /api/v1
func main() {
	// 从配置文件读取配置
	config.Init()

	// 初始化日志
	middleware.InitLogger()

	// 装载路由
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := server.NewRouter()

	// // 使用日志中间件
	// r.Use(middleware.Logger())
	// r.Use(middleware.RecoveryWithZap())

	middleware.GetZapLogger().Info("服务器正在启动，监听端口 :3000")
	if err := r.Run(":3000"); err != nil {
		middleware.GetZapLogger().Error("服务器启动失败", zap.Error(err))
	}
}
