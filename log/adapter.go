package log

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// 控制台日志适配器
type Adapter interface {
	getEncoder() zapcore.Encoder
	getWriter() zapcore.WriteSyncer
	getPriority() zap.LevelEnablerFunc
}

// 控制台日志适配器
type ConsoleAdapter struct {
	Level string //日志级别
}

// 文件日志适配器
type FileAdapter struct {
	File       string //文件（包含路径）
	MaxSize    int    //在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups int    //保留旧文件的最大个数
	MaxAge     int    //保留旧文件的最大天数
	Compress   bool   //是否压缩
	Level      string //日志级别
}

// kafka日志适配器
type KafkaAdapter struct {
	Level string   //日志级别
	Topic string   //主题
	Addr  []string //服务地址
}

// getEncoder 获取编码器
func (f *FileAdapter) getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

// getEncoder 获取控制台编码器
func (c *ConsoleAdapter) getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getEncoder 获取控制台编码器
func (k *KafkaAdapter) getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getFileWriter
func (c *ConsoleAdapter) getWriter() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stdout)
}

// getFileWriter
func (f *FileAdapter) getWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   f.File,
		MaxSize:    f.MaxSize,    //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: f.MaxBackups, //保留旧文件的最大个数
		MaxAge:     f.MaxAge,     //保留旧文件的最大天数
		Compress:   f.Compress,
		LocalTime:  true,
	}

	return zapcore.AddSync(lumberJackLogger)
}

// getFileWriter
func (k *KafkaAdapter) getWriter() zapcore.WriteSyncer {
	var kafkaLog KafkaLog
	var err error
	kafkaLog.Topic = k.Topic
	// 设置日志输入到Kafka的配置
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机的分区类型
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	kafkaLog.Producer, err = sarama.NewSyncProducer(k.Addr, config)
	if err != nil {
		panic(fmt.Errorf("connect kafka failed: %+v\n", err))
	}

	return zapcore.AddSync(&kafkaLog)
}

// getPriority
func (f *FileAdapter) getPriority() zap.LevelEnablerFunc {
	// 文件日志级别
	filePriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= getLogLevel(f.Level)
	})
	return filePriority
}

// getPriority
func (c *ConsoleAdapter) getPriority() zap.LevelEnablerFunc {
	// 文件日志级别
	filePriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= getLogLevel(c.Level)
	})
	return filePriority
}

// getPriority
func (k *KafkaAdapter) getPriority() zap.LevelEnablerFunc {
	// 文件日志级别
	filePriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= getLogLevel(k.Level)
	})
	return filePriority
}

// setLogLevel 设置日志级别
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.WarnLevel
	}
}
