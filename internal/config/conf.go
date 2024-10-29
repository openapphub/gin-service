package config

import (
	"fmt"
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
	// 输出一下环境变量
	fmt.Println(os.Getenv("GIN_MODE"))
	fmt.Println(os.Getenv("PORT"))
	fmt.Println(os.Getenv("MYSQL_DSN"))
	fmt.Println(os.Getenv("REDIS_DSN"))

	// 初始化 zap logger
	middleware.InitLogger()

	// 使用 middleware 中的 zapLogger 初始化 util.Logger
	util.BuildLogger(middleware.GetZapLogger())

	// 读取翻译文件
	localesPath := findLocalesFile()
	if err := LoadLocales(localesPath); err != nil {
		util.Log().Panic("翻译文件加载失败")
	}

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()
}

// findLocalesFile 查找翻译文件
func findLocalesFile() string {
	possiblePaths := []string{
		"locales/zh-cn.yaml",
		"/app/locales/zh-cn.yaml",
		"internal/config/locales/zh-cn.yaml",
		"/app/internal/config/locales/zh-cn.yaml",
	}

	for _, path := range possiblePaths {
		absPath, _ := filepath.Abs(path)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	return ""
}
