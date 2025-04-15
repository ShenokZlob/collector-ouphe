package config

import (
	"flag"
	"os"

	"github.com/spf13/viper"
)

func InitConfig() *viper.Viper {
	filename := fetchConfigFile()
	return InitConfigByFilename(filename)
}

func InitConfigByFilename(filename string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(filename)
	config.AddConfigPath("$HOME")
	config.AddConfigPath(".")

	if err := config.ReadInConfig(); err != nil {
		panic("Error reading config file: " + err.Error())
	}

	return config
}

func fetchConfigFile() string {
	var filename string
	flag.StringVar(&filename, "config", "", "for configuration collector server")
	flag.Parse()
	if filename == "" {
		filename = os.Getenv("APP_CONFIG")
	}
	return filename
}
