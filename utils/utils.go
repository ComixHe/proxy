package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {

	if logger != nil {
		return logger
	}

	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	level := zap.DebugLevel // test for now

	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)

	l := zap.New(core)
	logger = l.Sugar().Named("deepin-proxy")

	return logger
}
