package stats

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestSummaryMean(t *testing.T) {
	var vals = []float64{1, 11, 13, 56}
	s := Summary{}
	s.AddValues(vals)
	var sum float64
	for _, v := range vals {
		sum += v
	}

	assert.Equal(t, sum/float64(len(vals)), s.Mean)
	assert.Equal(t, 56., s.Max)
	assert.Equal(t, 24.404576073624664, s.Sd())
	assert.Equal(t, sum, s.Sum)
}
