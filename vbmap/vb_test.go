package vbmap

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestAll(t *testing.T) {
	vb := New(1024).SetNodes([]string{"1.1.1.1:9001", "1.1.1.2:9001", "1.1.1.1:9001"})
	t.Logf("%#v", *vb)
	assert.Equal(t, "1.1.1.1:9001", vb.Node(2323233))

	vb.SetNodes([]string{"1.1.1.1:9001", "1.1.1.2:9001"})
	assert.Equal(t, "1.1.1.2:9001", vb.Node(2323233))
}
