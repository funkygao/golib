package math

import (
	"testing"

	"github.com/funkygao/assert"
)

func TestMinInt(t *testing.T) {
	assert.Equal(t, 3, MinInt(5, 3))
	assert.Equal(t, 3, MinInt(3, 3))
	assert.Equal(t, -1, MinInt(-1, 3))
	assert.Equal(t, 0, MinInt(0, 3))
}

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 5, MaxInt(3, 5))
	assert.Equal(t, 5, MaxInt(5, 5))
}
