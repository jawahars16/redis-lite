package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jawahars16/redis-lite/resp"
)

type Handler struct {
	dictionary dictionary
}

type dictionary interface {
	Set(key string, value any)
	Get(key string) (any, bool)
}

func NewHandler(dictionary dictionary) *Handler {
	return &Handler{
		dictionary: dictionary,
	}
}

func (h *Handler) Ping(args ...any) ([]byte, error) {
	return resp.Serialize(resp.SimpleStrings, "PONG")
}

func (h *Handler) Set(args ...any) ([]byte, error) {
	if len(args) < 2 {
		return nil, errors.New("wrong number of arguments for 'set' command")
	}
	key := args[0].(string)
	v := args[1]
	number, ok := toInt(v)
	if ok {
		h.dictionary.Set(key, number)
		return resp.Serialize(resp.SimpleStrings, "OK")
	}
	h.dictionary.Set(key, v.(string))
	return resp.Serialize(resp.SimpleStrings, "OK")
}

func (h *Handler) Get(args ...any) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("wrong number of arguments for 'get' command")
	}
	key := args[0].(string)
	v, _ := h.dictionary.Get(key)
	number, ok := toInt(v)
	if ok {
		value := strconv.Itoa(number)
		return resp.Serialize(resp.SimpleStrings, value)
	}
	return resp.Serialize(resp.SimpleStrings, v.(string))
}

func (h *Handler) Incr(args ...any) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("wrong number of arguments for 'incr' command")
	}
	key := args[0].(string)
	v, _ := h.dictionary.Get(key)
	number, ok := toInt(v)
	if ok {
		number = number + 1
		h.dictionary.Set(key, number)
		return resp.Serialize(resp.Integers, number)
	} else {
		return nil, errors.New("value is not an integer or out of range")
	}
}

func (h *Handler) Config(args ...any) ([]byte, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("wrong number of arguments for 'config' command")
	}
	cmd := args[0].(string)
	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of arguments for 'config|%s' command", cmd)
	}
	if strings.ToLower(cmd) == "get" {
		option := args[1].(string)
		switch strings.ToLower(option) {
		case "save":
			return resp.Serialize(resp.Arrays, []resp.ArrayItem{
				{
					Value:    "save",
					DataType: resp.SimpleStrings,
				},
				{
					Value:    "\"\"",
					DataType: resp.SimpleStrings,
				},
			})
		case "appendonly":
			return resp.Serialize(resp.Arrays, []resp.ArrayItem{
				{
					Value:    "appendonly",
					DataType: resp.SimpleStrings,
				},
				{
					Value:    "no",
					DataType: resp.SimpleStrings,
				},
			})
		default:
			return nil, fmt.Errorf("'config|get|%s' not implemented", option)
		}
	}
	return nil, fmt.Errorf("'config|%s' not implemented", cmd)
}

func toInt(value any) (int, bool) {
	number, isNumber := value.(int)
	if isNumber {
		return number, true
	}
	str, isString := value.(string)
	if isString {
		number, err := strconv.Atoi(str)
		if err != nil {
			return 0, false
		}
		return number, true
	}
	return 0, false
}
