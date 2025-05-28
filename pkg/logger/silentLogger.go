package logger

type SilentLogger struct{}

func (l SilentLogger) Info(msg string, fields ...Field)  {}
func (l SilentLogger) Error(msg string, fields ...Field) {}
func (l SilentLogger) With(fields ...Field) Logger       { return l }
func (l SilentLogger) Sync() error                       { return nil }
func (l SilentLogger) Warn(msg string, fields ...Field)  {}
func (l SilentLogger) Debug(msg string, fields ...Field) {}
