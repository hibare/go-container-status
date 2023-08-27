package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		fmt.Println(err)
		return
	}

	cli.NegotiateAPIVersion(ctx)
	info, _ := cli.Info(ctx)
	fmt.Println(info.Swarm.ControlAvailable)
	swarmInfo, err := cli.SwarmInspect(ctx)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Swarm is enabled")
		fmt.Printf("Swarm ID: %s\n", swarmInfo.ID)
	}

}
