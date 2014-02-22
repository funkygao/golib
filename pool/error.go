package pool

import (
	"errors"
)

var (
	CLOSED_ERR = errors.New("ResourcePool is closed")
)
