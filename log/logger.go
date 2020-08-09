package log

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 日志对象接口
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Sync() error
}

type logger struct {
	*zap.SugaredLogger
}

// NewLogger 创建新的日志实例
func NewLogger(adapters ...Adapter) (Logger, error) {
	var allCore []zapcore.Core

	if len(adapters) == 0 {
		return nil, errors.New("参数不能为空")
	}

	for _, adapter := range adapters {
		writer := adapter.getWriter()
		encoder := adapter.getEncoder()
		priority := adapter.getPriority()
		allCore = append(allCore, zapcore.NewCore(encoder, writer, priority))
	}

	core := zapcore.NewTee(allCore...)
	zapLogger := zap.New(core)
	return &logger{zapLogger.Sugar()}, nil
}
