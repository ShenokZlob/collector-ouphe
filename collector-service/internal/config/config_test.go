// go test ./internal/config -v
package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchConfigFile_FromFlag(t *testing.T) {
	// Сохраняем оригинальные аргументы и флаги
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Подменяем аргументы командной строки
	os.Args = []string{"cmd", "-config=test_config"}

	// Сбрасываем флаги, чтобы избежать "flag redefined" ошибки
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	filename := fetchConfigFile()
	assert.Equal(t, "test_config", filename)
}

func TestFetchConfigFile_FromEnv(t *testing.T) {
	// err := godotenv.Load()
	// require.NoError(t, err)

	// Сохраняем и восстанавливаем переменную окружения
	oldEnv := os.Getenv("APP_CONFIG")
	defer os.Setenv("APP_CONFIG", oldEnv)

	os.Setenv("APP_CONFIG", "env_config")

	// Сбрасываем флаги
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	filename := fetchConfigFile()
	assert.Equal(t, "env_config", filename)
}

func TestInitConfigByFilename_ValidFile(t *testing.T) {
	// Создаем временный конфиг-файл
	tmpFile, err := os.CreateTemp("", "test_config_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(`
key: value
nested:
  field: 123
`)
	require.NoError(t, err)
	tmpFile.Close()

	config := InitConfigByFilename(tmpFile.Name())

	assert.Equal(t, "value", config.GetString("key"))
	assert.Equal(t, 123, config.GetInt("nested.field"))
}

func TestInitConfigByFilename_InvalidFile(t *testing.T) {
	// Проверяем, что функция паникует при несуществующем файле
	assert.Panics(t, func() {
		InitConfigByFilename("non_existent_file.yaml")
	})
}
