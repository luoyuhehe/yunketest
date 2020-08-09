package log

import (
	"fmt"
	"os"
)

//默认创建的日志对象,比init先执行
var defaultLogger = newDefaultLogger()

// newDefaultLogger 创建默认日志对象
func newDefaultLogger() Logger {
	var level = os.Getenv("THINKGO_LOG_LEVEL")
	if level == "" {
		//设置默认日志级别
		level = "info"
	}

	defaultLogger, err := NewLogger(&ConsoleAdapter{
		Level: level,
	})
	if err != nil {
		panic(fmt.Errorf("构建默认日志实例发生错误:%s", err))
	}
	return defaultLogger
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	defaultLogger.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	defaultLogger.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	defaultLogger.Errorf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	defaultLogger.Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	defaultLogger.Fatalf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Sync() {
	if err := defaultLogger.Sync(); err != nil {
		panic(err)
	}
}
