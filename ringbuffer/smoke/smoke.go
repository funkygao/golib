// A smoke test for ringbuffer package.
package main

import (
	"fmt"
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
	msgN     = 1 << 6
	ringSize = 32
	batch    = 5
)

type producer struct {
	rb          *ringbuffer.RingBuffer
	wg          sync.WaitGroup
	kafkaSentOk *sequence.Sequence

	out   chan int
	rbOut []int
}

func newProducer() *producer {
	rb, _ := ringbuffer.New(ringSize)
	p := &producer{
		rb:          rb,
		kafkaSentOk: sequence.New(),
		out:         make(chan int, 25),
		rbOut:       []int{},
	}

	p.wg.Add(1)
	go p.senderWorker()

	p.wg.Add(1)
	go p.sendToKafka()

	return p
}

func (p *producer) brokenIdx() []int {
	r := make([]int, 0)
	var last = p.rbOut[0]
	for x, i := range p.rbOut[1:] {
		if i != last+1 {
			r = append(r, x)
		}

		last = i
	}

	return r
}

func (p *producer) spans() string {
	var s = "["
	spans := p.brokenIdx()
	j := 0
	for i, v := range p.rbOut {
		if i == spans[j] {
			s += color.Red("%d ", v)
		} else {
			s += fmt.Sprintf("%d ", v)
		}
	}

	return s[:len(s)-1] + "]"
}

func (p *producer) send(i int) {
	p.rb.Write(i)
}

func (p *producer) senderWorker() {
	defer p.wg.Done()

	for {
		if msg, ok := p.rb.ReadTimeout(time.Second); ok {
			i := msg.(int)
			p.rbOut = append(p.rbOut, i)
			p.out <- i
		} else {
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
			if sampling.SampleRateSatisfied(1800) {
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
					p.kafkaSentOk.Add(ints[j])
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

	log.Printf("msg=%d, kafka sent=%d", msgN, p.kafkaSentOk.Length())
	min, max, loss := p.kafkaSentOk.Summary()
	log.Println("kafka sent:", p.kafkaSentOk)
	log.Printf("%d-%d", min, max)
	if len(loss) == 0 {
		log.Println(color.Green("ok"))
	} else {
		log.Println(color.Red("lost %d", len(loss)))
		log.Println(loss)
	}

	log.Println(p.spans())
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
