package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/hibare/go-container-status/internal/config"
)

var (
	listenAddr string
	listenPort string
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	listenAddr = config.ListenAddr
	listenPort = config.ListenPort
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

func containerStatus(w http.ResponseWriter, r *http.Request) {
	containerFavorableConditions := FavorableConditions{"running", "healthy"}
	params := mux.Vars(r)

	containerName := strings.Replace(params["container"], "\n", "", -1)
	containerName = strings.Replace(containerName, "\r", "", -1)

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

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Good to see you")
}

// HandleRequests : API request handler
func HandleRequests() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", authMiddleware(http.HandlerFunc(home))).Methods("GET")
	r.HandleFunc("/container/{container}/status/", authMiddleware(http.HandlerFunc(containerStatus))).Methods("GET")
	r.HandleFunc("/ping/", ping).Methods("GET")

	log.Printf("Listening for address %s on port %s\n", listenAddr, listenPort)
	log.Printf("Token(s): %v", apiKeys)
	log.Fatal(http.ListenAndServe(listenAddr+":"+listenPort, r))
}
