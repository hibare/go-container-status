FROM golang:1.21.6-alpine AS base

# Build main app
FROM base AS build

# Install healthcheck cmd
RUN apk update \
    && apk add curl \
    && apk add cosign \
    && curl -sfL https://raw.githubusercontent.com/hibare/go-docker-healthcheck/main/install.sh | sh -s -- -d -v -b /usr/local/bin

WORKDIR /src/

COPY . /src/

RUN CGO_ENABLED=0 go build -o /bin/go_container_status ./cmd/go-container-status/main.go

# Generate final image
FROM scratch

COPY --from=build /bin/go_container_status /bin/go_container_status

COPY --from=healthcheck /bin/healthcheck /bin/healthcheck

HEALTHCHECK \
    --interval=30s \
    --timeout=3s \
    CMD ["healthcheck", "--url", "http://localhost:5000/ping/"]

EXPOSE 5000

ENTRYPOINT ["/bin/go_container_status"]