package io

import (
	"testing"

	"github.com/funkygao/assert"
)

func TestDirExists(t *testing.T) {
	assert.Equal(t, true, DirExists("/tmp"))
	assert.Equal(t, false, DirExists("/tmp____xxxx"))
}
