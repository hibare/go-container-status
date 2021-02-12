package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Container struct {
	Name string `json:"name"`
	Status string `json:"status"`
}
type Result struct {
	Status bool `json:"status"`
	Message string `json:"message"`
	Result []Container `json:containers`
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Good to see you")
}

func containerStatus(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	fmt.Fprintf(w, "Container: %v", params["container"])
	log.Printf("Checking status for container %s", params["container"])

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}

}

func ping(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "pong")
}

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/container/{container}/status/", containerStatus).Methods("GET")
	r.HandleFunc("/ping/", ping).Methods("GET")

	log.Fatal(http.ListenAndServe("0.0.0.0:5000", r))
}

func main() {
	fmt.Println("Starting")
	handleRequests()
}