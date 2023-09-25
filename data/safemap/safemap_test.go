package safemap_test

import (
	"testing"
	"time"

	"github.com/jawahars16/redis-lite/data/safemap"
	"github.com/jawahars16/redis-lite/data/safemap/option"
	"github.com/stretchr/testify/assert"
)

func Test_safeMapSetGet(t *testing.T) {
	m := safemap.New()
	m.Set("key", "value", nil)
	value, ok := m.Get("key")

	assert.Equal(t, true, ok)
	assert.Equal(t, "value", value)
}

func Test_safeMapEX(t *testing.T) {
	m := safemap.New()
	m.Set("key", "value", option.WithEX(3))
	value, ok := m.Get("key")

	assert.True(t, ok)
	assert.Equal(t, "value", value)

	<-time.After(time.Second * 3)

	value, ok = m.Get("key")
	assert.False(t, ok)
	assert.Equal(t, nil, value)
}

func Test_safeMapPX(t *testing.T) {
	m := safemap.New()
	m.Set("key", "value", option.WithPX(3000))
	value, ok := m.Get("key")

	assert.True(t, ok)
	assert.Equal(t, "value", value)

	<-time.After(time.Second * 3)

	value, ok = m.Get("key")
	assert.False(t, ok)
	assert.Equal(t, nil, value)
}
