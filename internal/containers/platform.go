package containers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type ContainerPlatform string

const (
	PlatformDockerStandalone = ContainerPlatform("docker-standalone")
	PlatformDockerSwarm      = ContainerPlatform("docker-swarm")
)

func DetermineContainerPlatform(cli *client.Client) (ContainerPlatform, error) {
	ctx := context.Background()

	info, err := cli.Info(ctx)
	if err != nil {
		if client.IsErrConnectionFailed(err) {
			slog.Warn("Failed to connect to docker daemon")
			return "", nil
		}

		return "", errors.WithMessage(err, "failed to fetch docker info")
	}

	if info.Swarm.NodeID == "" {
		return PlatformDockerStandalone, nil
	}

	return PlatformDockerSwarm, nil
}

func PlatformPreChecks() error {
	ctx := context.Background()

	cli, err := getClient()
	if err != nil {
		slog.Error("Error creating docker client", "error", err)
		return ErrClientGet
	}

	defer cli.Close()

	platform, err := DetermineContainerPlatform(cli)
	if err != nil {
		return err
	}

	if platform == PlatformDockerStandalone {
		slog.Info("Docker is running in standalone mode")
		return nil
	}

	if platform == PlatformDockerSwarm {
		info, err := cli.Info(ctx)
		if err != nil {
			slog.Error("Error fetching docker info", "error", err)
			return ErrClientInfo
		}

		node, _, err := cli.NodeInspectWithRaw(ctx, info.Swarm.NodeID)
		if err != nil {
			slog.Error("Error inspecting docker node", "error", err)
			return ErrNodeInspect
		}

		if node.ManagerStatus == nil || !node.ManagerStatus.Leader {
			return ErrNoSwarmManager
		}

		slog.Info("Docker is running in swarm mode")
		return nil
	}

	return ErrUnsupportedPlatform
}

func GetSwarmNodes(cli *client.Client) {
	ctx := context.Background()
	nodes, err := cli.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		panic(err)
	}

	for _, node := range nodes {
		nodeInfo, _, err := cli.NodeInspectWithRaw(ctx, node.ID)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Node ID: %s, Hostname: %s, Status: %s\n", nodeInfo.ID, nodeInfo.Description.Hostname, node.Status.State)
	}
}
