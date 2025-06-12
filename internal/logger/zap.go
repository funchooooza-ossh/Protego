package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var baseLogger *zap.Logger

func Init(prod bool) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Encoding = "console"
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.DisableStacktrace = false

	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	cfg.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
	}

	var err error
	baseLogger, err = cfg.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic(err)
	}
}

func L() *zap.Logger {
	if baseLogger == nil {
		Init(false)
	}
	return baseLogger
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return L()
	}
	rid := GetRequestID(ctx)
	if rid != "" {
		return L().With(zap.String("rid", shortRequestID(rid)))
	}
	return L()
}

// Log is public, intended for general use (caller skip = 1)
func Log(ctx context.Context, level zapcore.Level, msg string, fields ...zap.Field) {
	logWrite(ctx, level, msg, 1, fields...)
}

// LogFromWrapper is exported ONLY for use by wrap-log packages.
// Not for general use — use logger.Log() instead.
func LogFromWrapper(ctx context.Context, level zapcore.Level, msg string, fields ...zap.Field) {
	logWrite(ctx, level, msg, 2, fields...)
}

func logWrite(ctx context.Context, level zapcore.Level, msg string, skip int, fields ...zap.Field) {
	log := FromContext(ctx).WithOptions(zap.AddCallerSkip(skip))
	switch level {
	case zap.DebugLevel:
		log.Debug(msg, fields...)
	case zap.InfoLevel:
		log.Info(msg, fields...)
	case zap.WarnLevel:
		log.Warn(msg, fields...)
	case zap.ErrorLevel:
		log.Error(msg, fields...)
	case zap.DPanicLevel:
		log.DPanic(msg, fields...)
	case zap.PanicLevel:
		log.Panic(msg, fields...)
	case zap.FatalLevel:
		log.Fatal(msg, fields...)
	}
}

func shortRequestID(full string) string {
	if len(full) >= 8 {
		return full[:8]
	}
	return full
}
