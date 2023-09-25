# Redis Lite
[![Continuous Integration](https://github.com/jawahars16/redis-lite/actions/workflows/test.yml/badge.svg)](https://github.com/jawahars16/redis-lite/actions/workflows/test.yml)

This repo contains the lite version of redis server built with Golang. It supports very few commands (PING, SET, GET). The idea inspired from John Cricket's coding challenge. 

https://codingchallenges.fyi/

## How to run

```bash
make run
```

This runs the server on port 6379. The port can be changed by passing the port number as an argument.

```bash
go run cmd/redis-lite/main.go -p 6379
```

The lite server can be accessed through redis CLI or redis-benchmark.

```bash
redis-cli PING
```

This lite version currently supports concurrent requests as well. It can be verified by below command,

```
redis-benchmark -t SET,GET
```

## How to test

```bash
make test
```

This includes unit and integration tests. Integration tests are executed against the server running on port 6379. The server is running using a docker container.
