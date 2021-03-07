FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/go_container_status

FROM scratch
COPY --from=build /bin/go_container_status /bin/go_container_status
ENTRYPOINT ["/bin/go_container_status"]