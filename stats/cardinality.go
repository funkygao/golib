package stats

import (
	"errors"
	h "github.com/funkygao/hyperloglog"
	"hash/fnv"
	"math"
	"time"
)

var ErrUnkownType = errors.New("unknown type")

type hllCounter struct {
	hll       *h.HyperLogLog
	startedAt time.Time
}

func (this *hllCounter) reset() {
	this.hll.Reset()
	this.startedAt = time.Now()
}

func (this *hllCounter) count() uint64 {
	return this.hll.Count()
}

type CardinalityCounter struct {
	m  uint
	hc map[string]*hllCounter
}

func NewCardinalityCounter() *CardinalityCounter {
	c := &CardinalityCounter{m: uint(math.Pow(2, float64(18)))} // 2^8KB mem space
	c.hc = make(map[string]*hllCounter)
	return c
}

func (this *CardinalityCounter) Add(key string, data interface{}) (err error) {
	if _, ok := this.hc[key]; !ok {
		hc := &hllCounter{startedAt: time.Now()}
		this.hc[key] = hc
		this.hc[key].hll, err = h.New(this.m)
		if err != nil {
			return
		}
	}

	switch data.(type) {
	case string:
		hash := fnv.New32()
		hash.Write([]byte(data.(string)))
		this.hc[key].hll.Add(hash.Sum32())
	case int:
		this.hc[key].hll.Add(uint32(data.(int)))
	case int16:
		this.hc[key].hll.Add(uint32(data.(int16)))
	case int32:
		this.hc[key].hll.Add(uint32(data.(int32)))
	case int64:
		this.hc[key].hll.Add(uint32(data.(int64)))
	case uint:
		this.hc[key].hll.Add(uint32(data.(uint)))
	case uint16:
		this.hc[key].hll.Add(uint32(data.(uint16)))
	case uint32:
		this.hc[key].hll.Add(data.(uint32))
	case uint64:
		this.hc[key].hll.Add(uint32(data.(uint64)))
	default:
		err = ErrUnkownType
	}

	return
}

func (this *CardinalityCounter) Reset(key string) {
	if _, ok := this.hc[key]; !ok {
		return
	}

	this.hc[key].reset()
}

func (this *CardinalityCounter) ResetAll() {
	for key, _ := range this.hc {
		this.Reset(key)
	}
}

func (this *CardinalityCounter) Count(key string) uint64 {
	if _, ok := this.hc[key]; !ok {
		return 0
	}

	return this.hc[key].count()
}

func (this *CardinalityCounter) StartedAt(key string) time.Time {
	if _, ok := this.hc[key]; !ok {
		return time.Time{}
	}

	return this.hc[key].startedAt
}

func (this *CardinalityCounter) Categories() []string {
	val := make([]string, 0, 10)
	for key, _ := range this.hc {
		val = append(val, key)
	}

	return val
}
