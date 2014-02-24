package slab

import (
	"errors"
)

var (
	ErrOutsideArena      = errors.New("buf not from this area")
	ErrInvalidRefCount   = errors.New("unexpected ref-count")
	ErrPushChunkRefCount = errors.New("pushFreeChunk() when non-zero refs")
)
