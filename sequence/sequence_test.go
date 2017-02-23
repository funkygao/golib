package sequence

import (
	"testing"

	"github.com/funkygao/assert"
)

func TestSequence(t *testing.T) {
	s := New()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}

	min, max, loss := s.Summary()
	assert.Equal(t, 0, min)
	assert.Equal(t, 9, max)
	assert.Equal(t, 0, len(loss))

	s.Add(12)
	min, max, loss = s.Summary()
	assert.Equal(t, 0, min)
	assert.Equal(t, 12, max)
	assert.Equal(t, 1, len(loss))
	assert.Equal(t, 10, loss[0]) // 0 1 2 3 4 5 6 7 8 9 12

	s.Add(13)
	min, max, loss = s.Summary()
	assert.Equal(t, 0, min)
	assert.Equal(t, 13, max)
	assert.Equal(t, 1, len(loss))
	assert.Equal(t, 10, loss[0]) // 0 1 2 3 4 5 6 7 8 9 12 13

	s.Add(13)
	min, max, loss = s.Summary()
	t.Logf("%+v", loss)
	assert.Equal(t, 0, min)
	assert.Equal(t, 13, max)
	assert.Equal(t, 1, len(loss))
	assert.Equal(t, 10, loss[0]) // 0 1 2 3 4 5 6 7 8 9 12 13 13

	s.Add(15)
	_, _, loss = s.Summary()
	t.Logf("%+v", loss)
	assert.Equal(t, 2, len(loss))
	assert.Equal(t, 10, loss[0]) // 0 1 2 3 4 5 6 7 8 9 12 13 13 15
	assert.Equal(t, 14, loss[1])
}
