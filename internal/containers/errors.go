package containers

import "errors"

var (
	ErrClientGet           = errors.New("failed to get docker client")
	ErrPlatformPreCheck    = errors.New("platform pre-check failed")
	ErrClientInfo          = errors.New("failed to fetch docker info")
	ErrNodeInspect         = errors.New("failed to inspect node")
	ErrNoSwarmManager      = errors.New("node is not a swarm manager")
	ErrUnsupportedPlatform = errors.New("unsupported platform")
	ErrUnhealthyContainers = errors.New("unhealthy containers found")
	ErrNoContainersFound   = errors.New("no containers found")
)
