package golib

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestUUID(t *testing.T) {
	uuid, _ := UUID()
	t.Logf("%v\n", uuid)
	assert.Equal(t, 32, len(uuid))
}
