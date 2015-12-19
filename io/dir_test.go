package io

import (
	"os"
	"testing"

	"github.com/funkygao/assert"
)

func TestDirExists(t *testing.T) {
	assert.Equal(t, true, DirExists("/tmp"))
	assert.Equal(t, false, DirExists("/tmp____xxxx"))
}

func TestDirChildren(t *testing.T) {
	predicate := func(path string, info os.FileInfo, err error) bool {
		if err != nil {
			return false
		}
		if info.IsDir() {
			return false
		}

		return true
	}
	children := DirChildren(".", predicate)
	t.Logf("%+v", children)
}
