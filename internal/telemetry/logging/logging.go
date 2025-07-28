package logging

import (
	"os"

	"github.com/obsidian-project-plugin/auth-service/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *zap.SugaredLogger
)

func Init(cfg *config.Config) {
	if cfg.Stage.IsDev {
		initDevelopment()
		return
	}
	initProduction(cfg.Stage.LogFilePath)
}

func initDevelopment() {
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
	Logger = l.Sugar()
}

func initProduction(logFilePath string) {
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
	Logger = l.Sugar()
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
	_ = Logger.Sync()
	os.Exit(1)
}

func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
	_ = Logger.Sync()
	os.Exit(1)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}
