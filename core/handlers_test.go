package core_test

import (
	"strings"
	"testing"

	"github.com/jawahars16/redis-lite/core"
	"github.com/jawahars16/redis-lite/resp"
	"github.com/stretchr/testify/assert"
)

func Test_handlePing(t *testing.T) {
	bytes, err := core.HandlePing()
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
