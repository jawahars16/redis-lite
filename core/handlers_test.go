package core_test

import (
	"strings"
	"testing"

	"github.com/jawahars16/redis-lite/core"
	"github.com/jawahars16/redis-lite/resp"
	"github.com/stretchr/testify/assert"
)

type simpleDictionary struct {
	m map[string]interface{}
}

func (d *simpleDictionary) Set(key string, value interface{}) {
	d.m[key] = value
}

func (d *simpleDictionary) Get(key string) (interface{}, bool) {
	v, ok := d.m[key]
	return v, ok
}

func Test_handlePing(t *testing.T) {
	h := core.NewHandler(&simpleDictionary{m: make(map[string]interface{})})
	bytes, err := h.Ping()
	if err != nil {
		t.Error(err)
	}
	dataType, data, err := resp.Deserialize(strings.NewReader(string(bytes)))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.SimpleStrings, dataType)
	assert.Equal(t, "PONG", data)
}

func Test_handleSetWithNumber(t *testing.T) {
	h := core.NewHandler(&simpleDictionary{m: make(map[string]interface{})})
	bytes, err := h.Set("key", 1)
	if err != nil {
		t.Error(err)
	}
	dataType, data, err := resp.Deserialize(strings.NewReader(string(bytes)))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.SimpleStrings, dataType)
	assert.Equal(t, "OK", data)
}

func Test_handleSetWithString(t *testing.T) {
	h := core.NewHandler(&simpleDictionary{m: make(map[string]interface{})})
	bytes, err := h.Set("key", "a")
	if err != nil {
		t.Error(err)
	}
	dataType, data, err := resp.Deserialize(strings.NewReader(string(bytes)))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.SimpleStrings, dataType)
	assert.Equal(t, "OK", data)
}

func Test_handleGetWithString(t *testing.T) {
	h := core.NewHandler(&simpleDictionary{m: make(map[string]interface{})})
	h.Set("key", "1")
	bytes, err := h.Get("key")
	if err != nil {
		t.Error(err)
	}
	dataType, data, err := resp.Deserialize(strings.NewReader(string(bytes)))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.SimpleStrings, dataType)
	assert.Equal(t, "1", data)
}

func Test_handleGetWithNumber(t *testing.T) {
	h := core.NewHandler(&simpleDictionary{m: make(map[string]interface{})})
	h.Set("key", 1)
	bytes, err := h.Get("key")
	if err != nil {
		t.Error(err)
	}
	dataType, data, err := resp.Deserialize(strings.NewReader(string(bytes)))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.SimpleStrings, dataType)
	assert.Equal(t, "1", data)
}

func Test_handleIncr(t *testing.T) {
	h := core.NewHandler(&simpleDictionary{m: make(map[string]interface{})})
	h.Set("key", 1)
	bytes, err := h.Incr("key")
	if err != nil {
		t.Error(err)
	}
	dataType, data, err := resp.Deserialize(strings.NewReader(string(bytes)))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, resp.Integers, dataType)
	assert.Equal(t, 2, data)
}

func Test_handleConfig(t *testing.T) {
	h := core.NewHandler(&simpleDictionary{m: make(map[string]interface{})})
	bytes, err := h.Config("get", "save")
	if err != nil {
		t.Error(err)
	}
	dataType, data, err := resp.Deserialize(strings.NewReader(string(bytes)))
	if err != nil {
		t.Error(err)
	}
	items := data.([]resp.ArrayItem)
	assert.Equal(t, resp.Arrays, dataType)
	assert.Equal(t, "save", items[0].Value)
	assert.Equal(t, resp.SimpleStrings, items[0].DataType)
	assert.Equal(t, "\"\"", items[1].Value)
	assert.Equal(t, resp.SimpleStrings, items[1].DataType)
}
