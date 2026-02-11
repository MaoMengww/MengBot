package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.SugaredLogger

func InitLogger() {

	//配置Encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 日志等级大写
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // 时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zap.DebugLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zap.InfoLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zap.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zap.ErrorLevel
	})
	fatalLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zap.FatalLevel
	})


	//多核心，不同地方输出日志
	cores := [...]zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, infoLevel),
		zapcore.NewCore(
			encoder,
			getWriteSyncer("./data/logs/info.log"),
			infoLevel,
		),
		zapcore.NewCore(
			encoder,
			getWriteSyncer("./data/logs/debug.log"),
			debugLevel,
		),
		zapcore.NewCore(
			encoder,
			getWriteSyncer("./data/logs/warn.log"),
			warnLevel,
		),
		zapcore.NewCore(
			encoder,
			getWriteSyncer("./data/logs/error.log"),
			errorLevel,
		),
		zapcore.NewCore(
			encoder,
			getWriteSyncer("./data/logs/fatal.log"),	
			fatalLevel,
		),
	}

	//合并核心
	logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

// 创建日志轮转器，写入文件/进行轮转
func getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    100,  // MB
		MaxBackups: 30,   // 最大备份数
		MaxAge:     30,   // 保存天数
		Compress:   true, // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 封装日志方法
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
