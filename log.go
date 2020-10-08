package thinkgo

import (
	"fmt"
	"github.com/sahara-go/thinkgo/log"
)

//默认创建的日志对象,比init先执行
var defaultLogger = newDefaultLogger()

// newDefaultLogger 创建默认日志对象
func newDefaultLogger() log.Logger {
	var adapters []log.Adapter

	//日志输出到控制台
	if AppConfig.Log.Console.Enable {
		adapter := &log.ConsoleAdapter{
			Level: AppConfig.Log.Console.Level,
		}
		adapters = append(adapters, adapter)
	}

	if AppConfig.Log.Kafka.Enable {
		adapter := &log.KafkaAdapter{
			Topic: AppConfig.Log.Kafka.Topic,
			Level: AppConfig.Log.Kafka.Level,
			Addr:  AppConfig.Log.Kafka.Addr,
		}

		adapters = append(adapters, adapter)
	}

	//日志输出到文件
	if len(adapters) < 1 || AppConfig.Log.File.Enable {
		adapter := &log.FileAdapter{
			File:       AppConfig.Log.File.Path + "/appConfigS.log",
			MaxSize:    AppConfig.Log.File.MaxSize,
			MaxBackups: AppConfig.Log.File.MaxBackups,
			MaxAge:     AppConfig.Log.File.MaxAge,
			Compress:   false,
			Level:      AppConfig.Log.File.Level,
		}

		adapters = append(adapters, adapter)
	}

	defaultLogger, err := log.NewLogger(adapters...)
	if err != nil {
		panic(fmt.Errorf("构建默认日志实例发生错误:%s", err))
	}

	return defaultLogger
}

// GetLogger 获取默认的日志实例
func GetLogger() log.Logger {
	return defaultLogger
}
