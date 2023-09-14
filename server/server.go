package server

import (
	"net"

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
func (r *RedisLite) Listen(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	r.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		dataType, data, err := resp.Deserialize(conn)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		switch dataType {
		case resp.SimpleStrings:
			r.handlers[data.(string)]()
		case resp.Arrays:
			items := data.([]resp.ArrayItem)
			if len(items) < 1 {
				slog.Error(err.Error())
				continue
			}
			command := items[0].Value.(string)
			args := []any{}
			for _, a := range items[1:] {
				args = append(args, a.Value)
			}
			bytes, err := r.handlers[command](args...)
			if err != nil {
				slog.Error(err.Error())
			}
			conn.Write(bytes)
		default:
			data, err := resp.Serialize(resp.SimpleErrors, resp.ErrUnrecognizedType.Error())
			if err != nil {
				slog.Error(err.Error())
			}
			conn.Write(data)
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
