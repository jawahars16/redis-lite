package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jawahars16/redis-lite/data/safemap/option"
	"github.com/jawahars16/redis-lite/resp"
)

type Handler struct {
	data   dictionary
	config dictionary
}

type dictionary interface {
	Set(key string, value any, expiry *option.ExpiryOption)
	Get(key string) (any, bool)
}

func NewHandler(data dictionary, config dictionary) *Handler {
	return &Handler{
		data:   data,
		config: config,
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
		h.data.Set(key, number, nil)
		return resp.Serialize(resp.SimpleStrings, "OK")
	}
	h.data.Set(key, v.(string), nil)
	return resp.Serialize(resp.SimpleStrings, "OK")
}

func (h *Handler) Get(args ...any) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("wrong number of arguments for 'get' command")
	}
	key := args[0].(string)
	v, exists := h.data.Get(key)
	if !exists {
		return resp.Serialize(resp.BulkStrings, nil)
	}
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
	v, _ := h.data.Get(key)
	number, ok := toInt(v)
	if ok {
		number = number + 1
		h.data.Set(key, number, nil)
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
	fmt.Println(cmd, args)
	if strings.ToLower(cmd) == "get" {
		key := args[1].(string)
		value, _ := h.config.Get(strings.ToLower(key))
		// var valueDataType resp.DataType
		// number, isNumber := toInt(value)
		// if isNumber {
		// 	valueDataType = resp.Integers
		// 	value = number
		// } else {
		// 	valueDataType = resp.SimpleStrings
		// }
		return resp.Serialize(resp.Arrays, []resp.ArrayItem{
			{
				Value:    key,
				DataType: resp.SimpleStrings,
			},
			{
				Value:    value,
				DataType: resp.SimpleStrings,
			},
		})
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
