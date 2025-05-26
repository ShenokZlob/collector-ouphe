package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Warning(msg string, fields ...Field)
	With(fields ...Field) Logger
	Sync() error
}

type zapLogger struct {
	l *zap.Logger
}

func NewZapLogger(isProduction bool) (Logger, error) {
	var log *zap.Logger
	var err error

	if isProduction {
		log, err = zap.NewProduction()
	} else {
		log, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}
	return &zapLogger{l: log}, nil
}

func (zl *zapLogger) Info(msg string, fields ...Field) {
	zl.l.Info(msg, toZapFields(fields)...)
}

func (zl *zapLogger) Error(msg string, fields ...Field) {
	zl.l.Error(msg, toZapFields(fields)...)
}

func (zl *zapLogger) Warning(msg string, fields ...Field) {
	zl.l.Warn(msg, toZapFields(fields)...)
}

func (zl *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{l: zl.l.With(toZapFields(fields)...)}
}

func (zl *zapLogger) Sync() error {
	return zl.l.Sync()
}

type Field struct {
	zapField zap.Field
}

func String(key, value string) Field {
	return Field{zapField: zap.String(key, value)}
}

func Error(err error) Field {
	return Field{zapField: zap.Error(err)}
}

func Int(key string, value int) Field {
	return Field{zapField: zap.Int(key, value)}
}

func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.zapField
	}
	return zapFields
}
