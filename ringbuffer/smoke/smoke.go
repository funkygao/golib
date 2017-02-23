// A smoke test for ringbuffer package.
package main

import (
	"log"
	"sync"
	"time"

	"github.com/funkygao/gafka/diagnostics/agent"
	"github.com/funkygao/golib/color"
	"github.com/funkygao/golib/ringbuffer"
	"github.com/funkygao/golib/sampling"
	"github.com/funkygao/golib/sequence"
)

const (
	msgN     = 1 << 7
	ringSize = 32
	batch    = 5
)

type producer struct {
	rb  *ringbuffer.RingBuffer
	wg  sync.WaitGroup
	seq *sequence.Sequence

	out chan int
}

func newProducer() *producer {
	rb, _ := ringbuffer.New(ringSize)
	p := &producer{
		rb:  rb,
		seq: sequence.New(),
		out: make(chan int, 25),
	}

	p.wg.Add(1)
	go p.senderWorker()

	p.wg.Add(1)
	go p.sendToKafka()

	return p
}

func (p *producer) send(i int) {
	p.rb.Write(i)
}

func (p *producer) senderWorker() {
	defer p.wg.Done()

	for {
		if msg, ok := p.rb.ReadTimeout(time.Second); ok {
			i := msg.(int)
			p.out <- i
		} else {
			log.Println("bye from sender worker!")
			close(p.out)
			break
		}
	}
}

func (p *producer) sendToKafka() {
	defer p.wg.Done()

	b := 0
	ints := make([]int, batch)
	for i := range p.out {
		ints[b] = i

		b++
		if b == batch {
			if sampling.SampleRateSatisfied(1200) {
				// fails
				log.Println(color.Red("%+v", ints))

				for j := 0; j < batch; j++ {
					p.rb.Rewind()
				}
			} else {
				// success
				log.Println(color.Green("%+v", ints))

				for j := 0; j < batch; j++ {
					p.rb.Advance()
					p.seq.Add(ints[j])
				}
			}

			b = 0
			ints = ints[0:]
		}
	}
}

func (p *producer) close() {
	log.Println("closing...")
	p.wg.Wait()

	log.Println(p.seq.Length())
	min, max, loss := p.seq.Summary()
	log.Println(p.seq)
	log.Printf("%d-%d", min, max)
	if len(loss) == 0 {
		log.Println(color.Green("ok"))
	} else {
		log.Println(color.Red("lost %d", len(loss)))
		log.Println(loss)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	agent.Start()

	inChan := make(chan int)
	go func() {
		for i := 0; i < msgN; i++ {
			inChan <- i
		}

		close(inChan)
	}()

	p := newProducer()
	for i := range inChan {
		p.send(i)
	}

	p.close()
}
