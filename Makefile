SHELL=/bin/bash

UI := $(shell id -u)
GID := $(shell id -g)
MAKEFLAGS += -s
DOCKER_COMPOSE_PREFIX = HOST_UID=${UID} HOST_GID=${GID} docker-compose -f docker-compose.dev.yml

all: app-up

app-up:
	${DOCKER_COMPOSE_PREFIX} up go-container-status

clean: 
	${DOCKER_COMPOSE_PREFIX} down
	go mod tidy

.PHONY = all clean app-up
