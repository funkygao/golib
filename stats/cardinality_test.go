package stats

import (
	"github.com/funkygao/assert"
	"os"
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

	assert.Equal(t, []string{"dau", "msg"}, c.Categories())
}

func TestMAU(t *testing.T) {
	c := NewCardinalityCounter()
	c.Add("RS.uid.month", 6092491)
	c.Add("RS.uid.month", 12356497)
	for _, k := range c.Categories() {
		t.Logf("%s %d", k, c.Count(k))
	}
}

func TestDumpAndLoad(t *testing.T) {
	const FN = "c.gob"
	c := NewCardinalityCounter()
	c.Add("RS.uid.month", 6092491)
	c.Add("RS.uid.month", 12356497)
	c.Dump(FN)
	d := NewCardinalityCounter()
	d.Load(FN)
	assert.Equal(t, uint64(2), d.Count("RS.uid.month"))
	os.Remove(FN)
}
