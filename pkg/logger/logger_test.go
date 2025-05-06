// go test -v ./pkg/logger
package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewZapLoggerDevelopment(t *testing.T) {
	log, err := NewZapLogger(false)
	require.NoError(t, err)
	require.NotNil(t, log)

	log.Info("test info", String("key", "value"))
	log.Error("test error", Int("code", 500))
}

func TestNewZapLoggerProduction(t *testing.T) {
	log, err := NewZapLogger(true)
	require.NoError(t, err)
	require.NotNil(t, log)

	log.Info("test info")
}

func TestLoggerWith(t *testing.T) {
	log, err := NewZapLogger(false)
	require.NoError(t, err)

	newLog := log.With(String("key", "value"))
	require.NotNil(t, newLog)

	newLog.Info("with logger")
}
