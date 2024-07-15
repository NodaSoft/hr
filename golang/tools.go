package main

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// GenerateUuid generate a uuid using google/uuid package.
func GenerateUuid() uuid.UUID {
	return uuid.New()
}

// PanicOnError panic if err is not nil.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// NewLogger create an instance of *zap.Logger.
func NewLogger(level int8) *zap.Logger {
	encodeCfg := zap.NewProductionEncoderConfig()
	encodeCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encodeCfg),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zapcore.Level(level)),
	))
}
