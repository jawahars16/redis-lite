# Redis Lite
[![Continuous Integration](https://github.com/jawahars16/redis-lite/actions/workflows/test.yml/badge.svg)](https://github.com/jawahars16/redis-lite/actions/workflows/test.yml)

This repo contains the lite version of redis server built with Golang. It supports very few commands (PING, SET, GET). The idea inspired from John Cricket's coding challenge. 

https://codingchallenges.fyi/

## How to run

```bash
make run
```

This runs the server on port 6379.

## How to test

```bash
make test
```

This includes unit and integration tests. Integration tests are executed against the server running on port 6379. The server is running using a docker container.

## Benchmark

This lite version currently supports concurrent requests as well. It can be verified by below command,

```
redis-benchmark -t SET,GET
```