package server_test

import (
	"net"
	"testing"

	"github.com/jawahars16/redis-lite/server"
	"github.com/stretchr/testify/assert"
)

func Test_ServerCommand(t *testing.T) {
	pingHandlerInvoked := make(chan bool)
	redisLite := server.New()
	redisLite.Handle("PING", func(args ...any) ([]byte, error) {
		pingHandlerInvoked <- true
		return nil, nil
	})
	go redisLite.Listen("localhost", 3000)
	defer redisLite.Close()
	send([]byte("+PING\r\n"), "localhost:3000")

	assert.True(t, <-pingHandlerInvoked)
}

func Test_ServerCommandWithArgs(t *testing.T) {
	pingHandlerInvoked := make(chan bool)
	var arg1 any
	var arg2 any
	redisLite := server.New()
	redisLite.Handle("SET", func(args ...any) ([]byte, error) {
		arg1 = args[0]
		arg2 = args[1]
		pingHandlerInvoked <- true
		return nil, nil
	})
	go redisLite.Listen("localhost", 3000)
	defer redisLite.Close()
	send([]byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"), "localhost:3000")

	assert.True(t, <-pingHandlerInvoked)
	assert.Equal(t, "key", arg1)
	assert.Equal(t, "value", arg2)
}

func send(data []byte, addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	conn.Write(data)
	return nil
}
