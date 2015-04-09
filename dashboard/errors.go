package dashboard

import (
	"errors"
)

var (
	ErrEmptyHttpAddr   = errors.New("http listen addr empty")
	ErrEmptyDataSource = errors.New(("empty data source"))
)
