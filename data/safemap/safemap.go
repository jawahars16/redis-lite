package safemap

import (
	"time"

	"github.com/jawahars16/redis-lite/data/safemap/option"
)

type SafeMap struct {
	m         map[string]interface{}
	semaphore chan struct{}
}

func New() *SafeMap {
	return &SafeMap{
		m:         make(map[string]interface{}),
		semaphore: make(chan struct{}, 1),
	}
}

func (s *SafeMap) Set(key string, value interface{}, expiryOption *option.ExpiryOption) {
	s.set(key, value)
	if expiryOption != nil {
		go func() {
			<-time.After(expiryOption.Duration)
			delete(s.m, key)
		}()
	}
}

func (s *SafeMap) Get(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	val, ok := s.m[key]
	return val, ok
}

func (s *SafeMap) set(key string, value interface{}) {
	if key == "" {
		return
	}

	s.semaphore <- struct{}{} // acquire lock
	s.m[key] = value
	<-s.semaphore
}
