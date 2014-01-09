package bjtime

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestNowBj(t *testing.T) {
	nowBj := NowBj()
	t.Log(nowBj)
}

func TestTsToString(t *testing.T) {
	assert.Equal(t, "01-09 10:14:59", TsToString(1389233699))
}
