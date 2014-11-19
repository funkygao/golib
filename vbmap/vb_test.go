package vbmap

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestAll(t *testing.T) {
	vb := New(1024).SetNodes([]string{"1.1.1.1:9001", "1.1.1.2:9001", "1.1.1.1:9001"})
	t.Logf("%#v", *vb)
	assert.Equal(t, "1.1.1.1:9001", vb.Node("user:1"))

	vb.SetNodes([]string{"1.1.1.1:9001", "1.1.1.2:9001"})
	assert.Equal(t, "1.1.1.2:9001", vb.Node("user:1"))
	assert.Equal(t, "1.1.1.1:9001", vb.Node("user:2"))

	assert.Equal(t, uint32(0x3a5), vb.Hash("user:1"))
	assert.Equal(t, uint32(0x2ac), vb.Hash("user:2"))
}
