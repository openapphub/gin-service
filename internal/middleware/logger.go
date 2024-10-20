package middleware

import (
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLogger *zap.Logger

// InitLogger 初始化zap日志
func InitLogger() {
	// 设置日志输出格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志轮转
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	})

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer), zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	)

	// 创建logger
	zapLogger = zap.New(core, zap.AddCaller())

	zapLogger.Info("Logger initialized")
}

// Logger 返回一个Gin的中间件，用于记录API请求
func Logger() gin.HandlerFunc {
	return ginzap.Ginzap(zapLogger, time.RFC3339, true)
}

// RecoveryWithZap 返回一个Gin的中间件，用于恢复panic并记录
func RecoveryWithZap() gin.HandlerFunc {
	return ginzap.RecoveryWithZap(zapLogger, true)
}

// GetZapLogger 返回zap logger实例，以便在其他地方使用
func GetZapLogger() *zap.Logger {
	return zapLogger
}
