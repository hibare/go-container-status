package middlewares

import (
	"log"
	"net/http"

	"github.com/hibare/go-container-status/internal/config"
	"github.com/hibare/go-container-status/internal/utils"
)

func TokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Client: [%s] %s\n", r.RemoteAddr, r.RequestURI)

		apiKey := r.Header.Get("Authorization")

		if utils.SliceContains(apiKey, config.Current.APIKeys) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}
