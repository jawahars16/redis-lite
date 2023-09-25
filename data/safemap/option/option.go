package option

import (
	"time"
)

type ExpiryOption struct {
	Duration time.Duration
}

func WithEX(seconds int) *ExpiryOption {
	return &ExpiryOption{
		Duration: time.Second * time.Duration(seconds),
	}
}

func WithPX(milliseconds int) *ExpiryOption {
	return &ExpiryOption{
		Duration: time.Millisecond * time.Duration(milliseconds),
	}
}
