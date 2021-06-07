package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	listenAddr string   `mapstructure:"LISTEN_ADDR"`
	listenPort string   `mapstructure:"LISTEN_PORT"`
	apiKeys    []string `mapstructure:"API_KEYS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
