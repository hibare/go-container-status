# go-container-status



[![Go Report Card](https://goreportcard.com/badge/github.com/hibare/go-container-status)](https://goreportcard.com/report/github.com/hibare/go-container-status)
[![Docker Hub](https://img.shields.io/docker/pulls/hibare/go-container-status)](https://hub.docker.com/r/hibare/go-container-status)
[![Docker image size](https://img.shields.io/docker/image-size/hibare/go-container-status/latest)](https://hub.docker.com/r/hibare/go-container-status) 
[![GitHub issues](https://img.shields.io/github/issues/hibare/go-container-status)](https://github.com/hibare/go-container-status/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/hibare/go-container-status)](https://github.com/hibare/go-container-status/pulls)
[![GitHub](https://img.shields.io/github/license/hibare/go-container-status)](https://github.com/hibare/go-container-status/blob/main/LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/hibare/go-container-status)](https://github.com/hibare/go-container-status/releases)


A REST API designed in go to check for a container status

> [!IMPORTANT]
> This repository is no longer maintained and has been superseded by [ArguSwarm](https://github.com/hibare/ArguSwarm)


## Getting Started

go-container-status is packaged as docker container. Docker image is available on [Docker Hub](https://hub.docker.com/r/hibare/go-container-status).

### Docker run

```shell
docker run -p 5000:5000 -v /var/run/docker.sock:/var/run/docker.sock -e LISTEN_ADDR='0.0.0.0' go-container-status
```

### Docker Compose

```yml
version: "3"
services:
  go-container-status:
    image: hibare/go-container-status
    container_name: go-container-status
    hostname: go-container-status
    environment: 
      - LISTEN_ADDR=0.0.0.0
    volumes:
        - /var/run/docker.sock:/var/run/docker.sock:ro
    ports:
      - "5000:5000"
    logging:
      driver: "json-file"
      options:
        max-size: "500k"
        max-file: "5"
```
## Endpoints

1. Check health

```shell
/ping/
```

2. Check a container status

```shell
/container/{container_name}/status/
```

Example:
```shell
> curl -H "Authorization: test" http://127.0.0.1:5000/container/postgres/status/

[{"Name":["/postgres"],"State":"exited","Status":"Exited (0) 2 hours ago","Image":"postgres:latest"}]

> curl -H "Authorization: test" http://127.0.0.1:5000/container/postgrsaes/status/

Not found
```

## Supported Environment Variables
| Variable | Description | Default Value |
| --------- | ----------- | ------------- |
| LISTEN_ADDR | IP address to bind to | 127.0.0.1 |
| LISTEN_PORT | Port to bind to | 5000 |
| API_KEYS | API Keys used for auth | Autogenerated |
