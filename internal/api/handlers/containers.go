package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hibare/go-container-status/internal/containers"
)

func ContainerStatus(w http.ResponseWriter, r *http.Request) {

	containerName := chi.URLParam(r, "container")

	foundContainers, err := containers.ContainerStatus(containerName)

	httpStatus := http.StatusOK

	if err != nil {
		switch err.Error() {
		case containers.ErrUnhealthyContainers:
			httpStatus = http.StatusInternalServerError
		case containers.ErrNoContainersFound:
			httpStatus = http.StatusNotFound
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(foundContainers); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
