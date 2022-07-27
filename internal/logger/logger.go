package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	zap *zap.Logger
}

var log Logger

func init() {
	z, _ := zap.NewProduction()
	log = Logger{zap: z}
}

func Panic(msg string, tags ...zap.Field) {
	log.zap.Panic(msg, tags...)
}

func Debug(msg string, tags ...zap.Field) {
	log.zap.Debug(msg, tags...)
}

func Info(msg string, tags ...zap.Field) {
	log.zap.Info(msg, tags...)
}

func Error(msg string, err error, tags ...zap.Field) {
	log.zap.Error(msg, tags...)
	log.zap.Sync()
}
