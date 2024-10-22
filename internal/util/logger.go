package util

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger 日志
type Logger struct {
	zapLogger *zap.Logger
}

var logger *Logger

// BuildLogger 构建logger
func BuildLogger(zapLogger *zap.Logger) {
	logger = &Logger{
		zapLogger: zapLogger,
	}
}

// Log 返回日志对象
func Log() *Logger {
	if logger == nil {
		panic("Logger not initialized. Call BuildLogger first.")
	}
	return logger
}

// Println 打印
func (ll *Logger) Println(msg string) {
	ll.zapLogger.Info(msg)
}

// Panic 极端错误
func (ll *Logger) Panic(format string, v ...interface{}) {
	ll.zapLogger.Panic(fmt.Sprintf(format, v...))
}

// Error 错误
func (ll *Logger) Error(format string, v ...interface{}) {
	ll.zapLogger.Error(fmt.Sprintf(format, v...))
}

// Warning 警告
func (ll *Logger) Warning(format string, v ...interface{}) {
	ll.zapLogger.Warn(fmt.Sprintf(format, v...))
}

// Info 信息
func (ll *Logger) Info(format string, v ...interface{}) {
	ll.zapLogger.Info(fmt.Sprintf(format, v...))
}

// Debug 校验
func (ll *Logger) Debug(format string, v ...interface{}) {
	ll.zapLogger.Debug(fmt.Sprintf(format, v...))
}
