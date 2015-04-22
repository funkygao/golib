package stress

import (
	"log"
	"sync"
	"time"
)

type waitGroupWrapper struct {
	wg sync.WaitGroup
}

func (this *waitGroupWrapper) Wrap(seq int, cb func(seq int)) {
	this.wg.Add(1)
	go func() {
		cb(seq)
		this.wg.Done()
	}()
}

func (this *waitGroupWrapper) Wait() {
	this.wg.Wait()
}

func RunStress(cb func(seq int)) {
	var waitGroup waitGroupWrapper
	var t0 = time.Now()
	for c := flags.c1; c <= flags.c2; c += flags.step {
		for r := 0; r < flags.round; r++ {
			log.Printf("concurrency: %6d started", c)
			t1 := time.Now()
			for i := 0; i < c; i++ {
				waitGroup.Wrap(i, cb)
			}
			waitGroup.Wait()
			log.Printf("concurrency: %6d elapsed: %s", c, time.Since(t1))
		}

	}

	log.Printf("All done in %s", time.Since(t0))
}
