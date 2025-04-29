package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger
}

type Field = zap.Field

type zapLogger struct {
	l *zap.Logger
}

func NewZapLogger(prod bool) Logger {
	var l *zap.Logger
	if prod {
		l, _ = zap.NewProduction()
	} else {
		l, _ = zap.NewDevelopment()
	}
	return &zapLogger{l: l}
}

func (z *zapLogger) Info(msg string, fields ...Field) {
	z.l.Info(msg, fields...)
}

func (z *zapLogger) Error(msg string, fields ...Field) {
	z.l.Error(msg, fields...)
}

func (z *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{l: z.l.With(fields...)}
}
