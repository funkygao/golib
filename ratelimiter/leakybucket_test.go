package ratelimiter

import (
	"testing"
	"time"
)

func TestLeakyBucketBasic(t *testing.T) {
	l := NewLeakyBucket(5, time.Second)
	t1 := time.Now()
	for i := 1; i < 20; i++ {
		if l.Pour(1) {
			t.Logf("%2d ok", i)
		} else {
			t.Logf("%2d fail", i)
			if !l.Pour(0) {
				t.Logf("%2d check ok", i)
			} else {
				t.Logf("%2d check fail", i)
			}
		}
	}
	t.Logf("%s", time.Since(t1))
}
