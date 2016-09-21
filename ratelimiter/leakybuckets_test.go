package ratelimiter

import (
	"testing"
	"time"

	"github.com/funkygao/assert"
)

func TestLeakyBucketsBasic(t *testing.T) {
	l := NewLeakyBuckets(5, time.Second)
	t1 := time.Now()
	for i := 1; i < 20; i++ {
		ok := l.Pour("key", 1)
		if i > 5 {
			assert.Equal(t, false, ok)
		}

		if !ok {
			ok1 := l.Pour("key", 0)
			assert.Equal(t, false, ok1)
		}
	}

	t.Logf("%s", time.Since(t1))
}
