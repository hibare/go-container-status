package main

import (
	"github.com/hibare/go-container-status/internal/api"
	"github.com/hibare/go-container-status/internal/config"
)

func init() {
	config.Load()
}

func main() {
	api.Serve()
}
