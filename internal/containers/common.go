package containers

import (
	"github.com/docker/docker/client"
)

func getClient() (*client.Client, error) {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}

	return cli, nil
}
