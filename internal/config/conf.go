package config

import (
	"openapphub/internal/middleware"
	"openapphub/internal/model"
	"openapphub/internal/util"
	"openapphub/pkg/cache"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	godotenv.Load()

	// 初始化 zap logger
	middleware.InitLogger()

	// 使用 middleware 中的 zapLogger 初始化 util.Logger
	util.BuildLogger(middleware.GetZapLogger())

	// 读取翻译文件
	localesPath := findLocalesFile()
	if err := LoadLocales(localesPath); err != nil {
		util.Log().Panic("翻译文件加载失败", err)
	}

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()
}

// findLocalesFile 查找翻译文件
func findLocalesFile() string {
	possiblePaths := []string{
		"internal/config/locales/zh-cn.yaml",
	}

	for _, path := range possiblePaths {
		absPath, _ := filepath.Abs(path)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	return ""
}
