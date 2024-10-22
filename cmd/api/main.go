package main

import (
	"fmt"
	"io"
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
	// middleware.InitLogger()

	// 装载路由
	gin.SetMode(os.Getenv("GIN_MODE"))
	// 禁用
	gin.DefaultWriter = io.Discard

	r := server.NewRouter()

	middleware.GetZapLogger().Info("服务器正在启动")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // 默认端口
	}
	fmt.Printf("服务器正在启动，监听端口：%s\n", port)

	if err := r.Run(":" + port); err != nil {
		middleware.GetZapLogger().Error("服务器启动失败", zap.Error(err))
	}
}
