package peer

import (
	"testing"

	"github.com/funkygao/assert"
)

func TestLamportClock(t *testing.T) {
	l := &LamportClock{}

	assert.Equal(t, LamportTime(0), l.Time())
	assert.Equal(t, LamportTime(1), l.Inc())
	assert.Equal(t, LamportTime(1), l.Time())

	l.Witness(41)
	assert.Equal(t, LamportTime(42), l.Time())

	l.Witness(41)
	assert.Equal(t, LamportTime(42), l.Time())

	l.Witness(32)
	assert.Equal(t, LamportTime(42), l.Time())
}
