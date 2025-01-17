package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func New(logLevel string) *zap.Logger {
	logLevelMap := map[string]zapcore.Level{
		"info":  zap.InfoLevel,
		"debug": zap.DebugLevel,
		"error": zap.ErrorLevel,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "msg",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(os.Stdout),
		logLevelMap[logLevel],
	)
	logger := zap.New(core)
	return logger

}
