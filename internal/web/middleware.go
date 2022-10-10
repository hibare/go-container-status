package web

import (
	"log"
	"net/http"

	"github.com/hibare/go-container-status/internal/config"
	"github.com/hibare/go-container-status/internal/util"
)

var (
	apiKeys []string
)

func init() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	apiKeys = config.APIKeys

	log.Printf("Found %v API keys", len(config.APIKeys))
}

func authMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Client: [%s] %s\n", r.RemoteAddr, r.RequestURI)

		apiKey := r.Header.Get("Authorization")

		if util.StringInSlice(apiKey, apiKeys) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}
