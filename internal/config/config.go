package config

import (
	"github.com/google/uuid"
	"github.com/hibare/go-container-status/internal/env"
)

type Config struct {
	ListenAddr string
	ListenPort int
	APIKeys    []string
}

var Current *Config

func Load() {
	token := []string{
		uuid.New().String(),
	}

	env.Load()

	Current = &Config{
		ListenAddr: env.MustString("LISTEN_ADDR", "0.0.0.0"),
		ListenPort: env.MustInt("LISTEN_PORT", 5000),
		APIKeys:    env.MustStringSlice("API_KEYS", token),
	}
}
