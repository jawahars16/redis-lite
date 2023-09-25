package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/jawahars16/redis-lite/resp"
	"golang.org/x/exp/slog"
)

type Handler func(args ...any) ([]byte, error)

type RedisLite struct {
	handlers map[string]Handler
	listener net.Listener
}

func New() *RedisLite {
	return &RedisLite{
		handlers: make(map[string]Handler),
	}
}

// Listen on the given address
func (r *RedisLite) Listen(ip string, port int) error {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})

	if err != nil {
		return err
	}
	defer listener.Close()
	r.listener = listener

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*3))
		go func() {
			r.handleConnection(ctx, conn)
		}()
		go func() {
			<-ctx.Done()
			conn.Close()
			cancel()
		}()
	}
}

func (r *RedisLite) handleConnection(ctx context.Context, conn *net.TCPConn) {
	for {
		dataType, data, err := resp.Deserialize(conn)
		if err != nil {
			if err == io.EOF {
				// no data. keep reading
				conn.Close()
			}
			continue
		}
		switch dataType {
		case resp.SimpleStrings:
			data, err := r.handlers[data.(string)]()
			if err != nil {
				errorData, _ := resp.Serialize(resp.SimpleErrors, err.Error())
				conn.Write(errorData)
				continue
			}
			conn.Write(data)
		case resp.Arrays:
			items := data.([]resp.ArrayItem)
			if len(items) < 1 {
				slog.Error(err.Error())
				writeError(err, conn)
				continue
			}
			command := items[0].Value.(string)
			args := []any{}
			for _, a := range items[1:] {
				args = append(args, a.Value)
			}
			handler, ok := r.handlers[strings.ToUpper(command)]
			if !ok {
				writeError(resp.ErrUnrecognizedType, conn)
				continue
			}
			bytes, err := handler(args...)
			if err != nil {
				data, err := resp.Serialize(resp.SimpleErrors, fmt.Sprintf("ERR %s", err.Error()))
				if err != nil {
					slog.Error(err.Error())
					writeError(err, conn)
				}
				conn.Write(data)
				continue
			}
			conn.Write(bytes)
		default:
			writeError(resp.ErrUnrecognizedType, conn)
		}
	}
}

func (r *RedisLite) Close() {
	if r == nil {
		panic("attempt to close before listen")
	}
	r.listener.Close()
}

// Register handler for the given command.
// Only one handler is possible for the command.
// In case of multiple handlers, last handler registered will considered.
func (r *RedisLite) Handle(command string, handler Handler) {
	r.handlers[command] = handler
}

func writeError(err error, conn io.Writer) {
	data, err := resp.Serialize(resp.SimpleErrors, err.Error())
	if err != nil {
		slog.Error(err.Error())
	}
	conn.Write(data)
}
