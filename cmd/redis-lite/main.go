package main

import (
	"flag"
	"fmt"

	"github.com/jawahars16/redis-lite/core"
	"github.com/jawahars16/redis-lite/data/safemap"
	"github.com/jawahars16/redis-lite/server"
	"golang.org/x/exp/slog"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Host to listen on")
	port := flag.Int("port", 6379, "Port to listen on")
	flag.Parse()

	handler := core.NewHandler(safemap.New())
	redisLite := server.New()

	redisLite.Handle("PING", handler.Ping)
	redisLite.Handle("SET", handler.Set)
	redisLite.Handle("GET", handler.Get)
	redisLite.Handle("INCR", handler.Incr)
	redisLite.Handle("CONFIG", handler.Config)

	slog.Info(fmt.Sprintf("Listening on %s:%d", *host, *port))
	errChan := make(chan error)
	go func() {
		if err := redisLite.Listen(*host, *port); err != nil {
			errChan <- err
		}
	}()
	slog.Info("Ready to accept connections")
	err := <-errChan
	slog.Error(err.Error())
}
