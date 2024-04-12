package containers

import (
	"context"
	"strings"

	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/hibare/GoCommon/v2/pkg/slice"
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

	slog.Info("Checking status for container", "container", containerName)

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Error creating docker client", "error", err)
		return foundContainers, err
	}
	defer cli.Close()

	options := types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("name", containerName),
		),
	}

	containers, err := cli.ContainerList(ctx, options)
	if err != nil {
		slog.Error("Error listing containers", "error", err)
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

			if !slice.SliceContains(container.State, containerFavorableConditions) {
				unhealthyContainers = true
			}
		}

		if unhealthyContainers {
			return foundContainers, ErrUnhealthyContainers
		} else {
			return foundContainers, nil
		}
	} else {
		return foundContainers, ErrNoContainersFound
	}
}

func ContainerStatusAll() ([]Container, error) {
	ctx := context.Background()
	foundContainers := []Container{}

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Error creating docker client", "error", err)
		return foundContainers, err
	}
	defer cli.Close()

	options := types.ContainerListOptions{
		All: true,
	}

	containers, err := cli.ContainerList(ctx, options)
	if err != nil {
		slog.Error("Error listing containers", "error", err)
		return foundContainers, err
	}

	for _, container := range containers {
		foundContainers = append(foundContainers, Container{
			Name:   container.Names,
			State:  container.State,
			Status: container.Status,
			Image:  container.Image,
		})
	}

	return foundContainers, nil
}
