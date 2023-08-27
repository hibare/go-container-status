package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hibare/go-container-status/internal/containers"
)

func ContainerStatus(w http.ResponseWriter, r *http.Request) {
	var httpStatus int

	containerName := chi.URLParam(r, "container")

	foundContainers, err := containers.ContainerStatus(containerName)

	switch err {
	case nil:
		httpStatus = http.StatusOK
	case containers.ErrUnhealthyContainers:
		httpStatus = http.StatusInternalServerError
	case containers.ErrNoContainersFound:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(foundContainers); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ContainerStatusAll(w http.ResponseWriter, r *http.Request) {
	var httpStatus int

	foundContainers, err := containers.ContainerStatusAll()

	switch err {
	case nil:
		httpStatus = http.StatusOK
	default:
		httpStatus = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpStatus)

	if err := json.NewEncoder(w).Encode(foundContainers); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
