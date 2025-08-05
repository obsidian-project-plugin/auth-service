// internal/telemetry/logging/logging.go
package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	Logger *zap.SugaredLogger
)

func Init(logFilePath string, isDev bool) {
	var sugarLogger *zap.SugaredLogger
	if isDev {
		sugarLogger = initDevelopment()
	} else {
		sugarLogger = initProduction(logFilePath)
	}

	Logger = sugarLogger
}

func initDevelopment() *zap.SugaredLogger {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.DebugLevel,
	)
	l := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.WarnLevel),
	)
	return l.Sugar()
}

func initProduction(logFilePath string) *zap.SugaredLogger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	})

	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.InfoLevel)
	fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), fileWriter, zap.InfoLevel)

	core := zapcore.NewTee(consoleCore, fileCore)

	l := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	return l.Sugar()
}

func Info(args ...interface{}) {
	if Logger != nil {
		Logger.Info(args...)
	}
}

func Infof(template string, args ...interface{}) {
	if Logger != nil {
		Logger.Infof(template, args...)
	}
}

func Debug(args ...interface{}) {
	if Logger != nil {
		Logger.Debug(args...)
	}
}

func Warn(args ...interface{}) {
	if Logger != nil {
		Logger.Warn(args...)
	}
}

func Error(args ...interface{}) {
	if Logger != nil {
		Logger.Error(args...)
		_ = Logger.Sync()
	}
}

func Errorf(template string, args ...interface{}) {
	if Logger != nil {
		Logger.Errorf(template, args...)
		_ = Logger.Sync()
	}
}

func Fatal(args ...interface{}) {
	if Logger != nil {
		Logger.Fatal(args...)
		_ = Logger.Sync()
	}
}

func Fatalf(template string, args ...interface{}) {
	if Logger != nil {
		Logger.Fatalf(template, args...)
		_ = Logger.Sync()
	}
}
