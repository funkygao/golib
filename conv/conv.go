package conv

import (
	"errors"
)

const (
	ascii_0 = '0'
	ascii_9 = '9'
)

var (
	ErrInvalidFormat = errors.New("invalid format")
)

func ParseInt(b []byte) (n int, err error) {
	if len(b) == 0 {
		return -1, ErrInvalidFormat
	}

	for _, v := range b {
		if v < ascii_0 || v > ascii_9 {
			return -1, ErrInvalidFormat
		}

		n = n*10 + int(v) - ascii_0
	}

	return
}
