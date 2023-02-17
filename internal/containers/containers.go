package containers

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/hibare/go-container-status/internal/utils"
)

var (
	ErrUnhealthyContainers = "unhealthy containers found"
	ErrNoContainersFound   = "no containers found"
)

type Container struct {
	Name   []string
	State  string
	Status string
	Image  string
}

func ContainerStatus(containerName string) ([]Container, error) {
	ctx := context.Background()
	foundContainers := []Container{}
	containerFavorableConditions := []string{"running", "healthy"}

	containerName = strings.Replace(containerName, "\n", "", -1)
	containerName = strings.Replace(containerName, "\r", "", -1)

	log.Printf("Checking status for container %s", containerName)

	cli, err := client.NewClientWithOpts()
	if err != nil {
		log.Fatal(err)
		return foundContainers, err
	}

	cli.NegotiateAPIVersion(ctx)

	containerFilter := filters.NewArgs()
	containerFilter.Add("name", containerName)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: containerFilter})
	if err != nil {
		log.Fatal(err)
		return foundContainers, err
	}

	unhealthyContainers := false

	if len(containers) > 0 {
		for _, container := range containers {
			foundContainers = append(foundContainers, Container{
				Name:   container.Names,
				State:  container.State,
				Status: container.Status,
				Image:  container.Image,
			})

			if !utils.SliceContains(container.State, containerFavorableConditions) {
				unhealthyContainers = true
			}
		}

		if unhealthyContainers {
			return foundContainers, errors.New(ErrUnhealthyContainers)
		} else {
			return foundContainers, nil
		}
	} else {
		return foundContainers, errors.New(ErrNoContainersFound)
	}
}
