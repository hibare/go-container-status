package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hibare/go-container-status/internal/api/handlers"
	"github.com/hibare/go-container-status/internal/api/middlewares"
	"github.com/hibare/go-container-status/internal/config"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Good to see you")
}

func Serve() {
	wait := time.Second * 15
	addr := fmt.Sprintf("%s:%d", config.Current.ListenAddr, config.Current.ListenPort)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.StripSlashes)
	r.Use(middleware.CleanPath)

	r.Get("/", home)
	r.Get("/ping", handlers.HealthCheck)
	r.Route("/container/{container}", func(r chi.Router) {
		r.Use(middlewares.TokenAuth)
		r.Get("/status", handlers.ContainerStatus)
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Listening for address %s on port %d\n", config.Current.ListenAddr, config.Current.ListenPort)

	log.Printf("Token(s): %v", config.Current.APIKeys)
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)
}
