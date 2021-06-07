# go-container-status

A REST API designed in go to check for a container status

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