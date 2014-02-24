package slab

import (
	"errors"
)

var (
	ErrTooBig            = errors.New("buf size too big")
	ErrNoChunkMem        = errors.New("no mem chunk left")
	ErrOutsideArena      = errors.New("buf not from this area")
	ErrInvalidRefCount   = errors.New("unexpected ref-count")
	ErrPushChunkRefCount = errors.New("pushFreeChunk() when non-zero refs")
)
