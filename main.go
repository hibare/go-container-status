package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/hibare/util"
)

var (
	config struct{}
)

func init() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	log.Printf("Found %v API keys", len(config.apiKeys))
}

// Container data holder
type Container struct {
	Name   []string
	State  string
	Status string
	Image  string
}

// FavorableConditions is container favorable conditions
type FavorableConditions []string

// Has is array has an element check implementation
func (items FavorableConditions) Has(i string) bool {

	for _, item := range items {
		if item == i {
			return true
		}
	}
	return false
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Good to see you")
}

func containerStatus(w http.ResponseWriter, r *http.Request) {
	containerFavorableConditions := FavorableConditions{"running", "healthy"}
	params := mux.Vars(r)

	containerName := params["container"]

	log.Printf("Checking status for container %s", containerName)

	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	cli.NegotiateAPIVersion(ctx)

	containerFilter := filters.NewArgs()
	containerFilter.Add("name", containerName)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: containerFilter})
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	healthFlag := false
	if len(containers) > 0 {
		foundContainers := []Container{}

		for _, container := range containers {

			state := container.State
			status := container.Status
			name := container.Names
			image := container.Image

			_container := Container{
				Name:   name,
				State:  state,
				Status: status,
				Image:  image,
			}

			foundContainers = append(foundContainers, _container)

			log.Print(foundContainers)

			if !containerFavorableConditions.Has(state) {
				healthFlag = true
			}
		}

		if healthFlag {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(foundContainers); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(foundContainers); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	} else {
		log.Printf("Container %s not found", containerName)
		http.Error(w, "Not found", http.StatusNotFound)
	}

}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", authMiddleware(http.HandlerFunc(home))).Methods("GET")
	r.HandleFunc("/container/{container}/status/", authMiddleware(http.HandlerFunc(containerStatus))).Methods("GET")
	r.HandleFunc("/ping/", ping).Methods("GET")

	log.Fatal(http.ListenAndServe(listenAddr+":"+listenPort, r))
}

func authMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Client: [%s] %s\n", r.RemoteAddr, r.RequestURI)

		apiKey := r.Header.Get("Authorization")

		if stringInSlice(apiKey, apiKeys) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}

func main() {
	log.Printf("Listening for address %s on port %s\n", listenAddr, listenPort)
	handleRequests()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
