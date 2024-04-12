package main

import (
	"log"

	"github.com/hibare/go-container-status/internal/api"
	"github.com/hibare/go-container-status/internal/config"
	"github.com/hibare/go-container-status/internal/containers"
)

func init() {
	config.Load()
}

func main() {
	if err := containers.PlatformPreChecks(); err != nil {
		log.Fatal(err)
	}

	api.Serve()
}
