package containers

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
			log.Warn(err)
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
		log.Error(err)
		return ErrClientGet
	}

	defer cli.Close()

	platform, err := DetermineContainerPlatform(cli)
	if err != nil {
		return err
	}

	if platform == PlatformDockerStandalone {
		log.Info("Docker is running in standalone mode")
		return nil
	}

	if platform == PlatformDockerSwarm {
		info, err := cli.Info(ctx)
		if err != nil {
			log.Error(err)
			return ErrClientInfo
		}

		node, _, err := cli.NodeInspectWithRaw(ctx, info.Swarm.NodeID)
		if err != nil {
			log.Error(err)
			return ErrNodeInspect
		}

		if node.ManagerStatus == nil || !node.ManagerStatus.Leader {
			return ErrNoSwarmManager
		}

		log.Info("Docker is running in swarm mode")
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
