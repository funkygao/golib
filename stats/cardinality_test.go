package stats

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestCardinalityCounter(t *testing.T) {
	c := NewCardinalityCounter()
	c.Add("dau", 34343434)
	c.Add("dau", 45454)
	c.Add("dau", 888)
	assert.Equal(t, uint64(3), c.Count("dau"))

	c.Reset("msg")
	c.Add("msg", "we are in China")
	c.Add("msg", "where are you")
	assert.Equal(t, uint64(2), c.Count("msg"))
}
