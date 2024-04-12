package config

import (
	"log"

	"github.com/google/uuid"
	"github.com/hibare/GoCommon/v2/pkg/env"
	commonLogger "github.com/hibare/GoCommon/v2/pkg/logger"
	"github.com/hibare/go-container-status/internal/constants"
)

type APIConfig struct {
	ListenAddr string
	ListenPort int
	APIKeys    []string
}

type LoggerConfig struct {
	Level string
	Mode  string
}

type Config struct {
	API    APIConfig
	Logger LoggerConfig
}

var Current *Config

func Load() {
	commonLogger.InitDefaultLogger()

	token := []string{
		uuid.New().String(),
	}

	env.Load()

	Current = &Config{
		API: APIConfig{
			ListenAddr: env.MustString("LISTEN_ADDR", constants.DefaultAPIListenAddr),
			ListenPort: env.MustInt("LISTEN_PORT", constants.DefaultAPIListenPort),
			APIKeys:    env.MustStringSlice("API_KEYS", token),
		},

		Logger: LoggerConfig{
			Level: env.MustString("LOG_LEVEL", constants.DefaultLoggerLevel),
			Mode:  env.MustString("LOG_MODE", constants.DefaultLoggerMode),
		},
	}

	if !commonLogger.IsValidLogLevel(Current.Logger.Level) {
		log.Fatal("Error invalid logger level")
	}

	if !commonLogger.IsValidLogMode(Current.Logger.Mode) {
		log.Fatal("Error invalid logger mode")
	}

	commonLogger.InitLogger(&Current.Logger.Level, &Current.Logger.Mode)
}
