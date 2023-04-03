package model

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SugaredLogger *zap.SugaredLogger

func determineLogLevel() zap.AtomicLevel {
	lvl := zap.NewAtomicLevel()
	lvl.SetLevel(zap.DebugLevel)
	return lvl
}

func CreateSugaredLogger() *zap.SugaredLogger {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		determineLogLevel(),
	))
	SugaredLogger = logger.Sugar()
	return SugaredLogger
}
