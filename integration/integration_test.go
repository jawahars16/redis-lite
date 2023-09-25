package main

import (
	"context"
	"log"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()
	container := bootstrap(ctx)
	m.Run()
	container.Terminate(ctx)
}

func Test_RedisLitePing(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	cmd := client.Ping(ctx)
	assert.Equal(t, "PONG", cmd.Val())
}

func Test_RedisLiteSet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	client.Set(ctx, "key", 1, 0)
	cmd := client.Get(ctx, "key")
	assert.Equal(t, "1", cmd.Val())
}

func Test_RedisLiteIncr(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	client.Set(ctx, "key", 1, 0)
	client.Incr(ctx, "key")
	cmd := client.Get(ctx, "key")
	assert.Equal(t, "2", cmd.Val())
}

// persistence is disabled.
// following test will ensure that there is no snapshotting configured
func Test_RedisLiteConfigGetSave(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cmd := client.ConfigGet(ctx, "save")
	data := cmd.Val()
	assert.Equal(t, "save", data[0])
	assert.Equal(t, "", data[1])
}

// persistence is disabled.
// following test will ensure that there is no AOF configured
func Test_RedisLiteConfigGetAppendOnly(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cmd := client.ConfigGet(ctx, "appendonly")
	assert.Equal(t, "appendonly", cmd.Val()[0])
	assert.Equal(t, "no", cmd.Val()[1])
}

func bootstrap(ctx context.Context) testcontainers.Container {
	request := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Dockerfile: "./Dockerfile",
			Context:    "../.",
		},
		ExposedPorts: []string{"6379:6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		log.Panic(err)
	}

	_, err = container.Host(ctx)
	if err != nil {
		log.Panic(err)
	}
	return container
}
