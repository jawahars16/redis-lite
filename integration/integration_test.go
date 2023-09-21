package main

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Test_RedisLitePing(t *testing.T) {
	ctx := context.Background()
	container := bootstrap(ctx, t)
	defer container.Terminate(ctx)

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	cmd := client.Ping(ctx)
	assert.Equal(t, "PONG", cmd.Val())
}

func Test_RedisLiteSet(t *testing.T) {
	ctx := context.Background()
	container := bootstrap(ctx, t)
	defer container.Terminate(ctx)

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	client.Set(ctx, "key", 1, 0)
	cmd := client.Get(ctx, "key")
	assert.Equal(t, "1", cmd.Val())
}

func Test_RedisLiteIncr(t *testing.T) {
	ctx := context.Background()
	container := bootstrap(ctx, t)
	defer container.Terminate(ctx)

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
	ctx := context.Background()
	container := bootstrap(ctx, t)
	defer container.Terminate(ctx)

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cmd := client.ConfigGet(ctx, "save")
	assert.Equal(t, "save", cmd.Val()[0])
	assert.Equal(t, "\"\"", cmd.Val()[1])
}

// persistence is disabled.
// following test will ensure that there is no AOF configured
func Test_RedisLiteConfigGetAppendOnly(t *testing.T) {
	ctx := context.Background()
	container := bootstrap(ctx, t)
	defer container.Terminate(ctx)

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cmd := client.ConfigGet(ctx, "appendonly")
	client.Pipeline()
	assert.Equal(t, "appendonly", cmd.Val()[0])
	assert.Equal(t, "no", cmd.Val()[1])
}

func bootstrap(ctx context.Context, t *testing.T) testcontainers.Container {
	request := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Dockerfile: "./Dockerfile",
			Context:    "../.",
		},
		ExposedPorts: []string{"6379:6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	container, error := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if error != nil {
		t.Fatal(error)
	}

	_, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return container
}
