package ringbuffer

import (
	"errors"
)

var (
	ErrInvalidQueueSize = errors.New("queueSize must be power of 2")
)
