// A smoke test for ringbuffer package.
package main

import (
	"log"
	"sync"
	"time"

	"github.com/funkygao/golib/color"
	"github.com/funkygao/golib/ringbuffer"
	"github.com/funkygao/golib/sampling"
	"github.com/funkygao/golib/sequence"
)

var (
	rb     *ringbuffer.RingBuffer
	inChan = make(chan int)
	seq    = sequence.New()
)

func init() {
	rb, _ = ringbuffer.New(32)
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 1<<10; i++ {
			inChan <- i
		}

		close(inChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if msg, ok := rb.ReadTimeout(time.Second); ok {
				i := msg.(int)
				if sampling.SampleRateSatisfied(100) {
					log.Printf("rewind for %d", i)
					rb.Rewind()
				} else {
					rb.Advance()
					seq.Add(i)
				}
			} else {
				log.Println("bye from sender worker")
				break
			}
		}
	}()

	for i := range inChan {
		rb.Write(i)
	}

	wg.Wait()

	min, max, loss := seq.Summary()
	log.Printf("%d-%d", min, max)
	if len(loss) == 0 {
		log.Println(color.Green("ok"))
	} else {
		log.Println(color.Red("lost %d", len(loss)))
		log.Println(loss)
	}

}
