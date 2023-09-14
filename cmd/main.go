package main

import (
	"github.com/jawahars16/redis-lite/core"
	"github.com/jawahars16/redis-lite/server"
	"golang.org/x/exp/slog"
)

func main() {
	redisLite := server.New()
	redisLite.Handle("PING", core.HandlePing)
	redisLite.Handle("COMMAND", core.HandleCommand)
	if err := redisLite.Listen("127.0.0.1:6379"); err != nil {
		slog.Error(err.Error())
	}
}
