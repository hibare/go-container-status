package util

import (
	"log"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Config struct {
	ListenAddr string   `mapstructure:"LISTEN_ADDR"`
	ListenPort string   `mapstructure:"LISTEN_PORT"`
	APIKeys    []string `mapstructure:"API_KEYS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	token, err := uuid.NewRandom()

	viper.SetDefault("LISTEN_ADDR", "127.0.0.1")
	viper.SetDefault("LISTEN_PORT", "5000")
	viper.SetDefault("API_KEYS", token)

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	log.Printf("Auth token(s)=%v", config.APIKeys)
	return
}
