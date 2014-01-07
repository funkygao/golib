package stats

import (
	"errors"
	h "github.com/funkygao/hyperloglog"
	"hash/fnv"
	"math"
)

type CardinalityCounter struct {
	m   uint
	hll map[string]*h.HyperLogLog
}

func NewCardinalityCounter() *CardinalityCounter {
	c := &CardinalityCounter{m: uint(math.Pow(2, float64(18)))} // 2^8KB mem space
	c.hll = make(map[string]*h.HyperLogLog)
	return c
}

func (this *CardinalityCounter) Add(key string, data interface{}) (err error) {
	if _, ok := this.hll[key]; !ok {
		this.hll[key], err = h.New(this.m)
		if err != nil {
			return
		}
	}

	switch data.(type) {
	case string:
		hash := fnv.New32()
		hash.Write([]byte(data.(string)))
		this.hll[key].Add(hash.Sum32())
	case int:
		this.hll[key].Add(uint32(data.(int)))
	case int16:
		this.hll[key].Add(uint32(data.(int16)))
	case int32:
		this.hll[key].Add(uint32(data.(int32)))
	case int64:
		this.hll[key].Add(uint32(data.(int64)))
	case uint:
		this.hll[key].Add(uint32(data.(uint)))
	case uint16:
		this.hll[key].Add(uint32(data.(uint16)))
	case uint32:
		this.hll[key].Add(data.(uint32))
	case uint64:
		this.hll[key].Add(uint32(data.(uint64)))
	default:
		err = errors.New("unknown type")
	}

	return
}

func (this *CardinalityCounter) Reset(key string) {
	if _, ok := this.hll[key]; !ok {
		return
	}

	this.hll[key].Reset()
}

func (this *CardinalityCounter) ResetAll() {
	for key, _ := range this.hll {
		this.Reset(key)
	}
}

func (this *CardinalityCounter) Count(key string) uint64 {
	if _, ok := this.hll[key]; !ok {
		return 0
	}

	return this.hll[key].Count()
}

func (this *CardinalityCounter) Categories() []string {
	val := make([]string, 0, 10)
	for key, _ := range this.hll {
		val = append(val, key)
	}

	return val
}

func (this *CardinalityCounter) Counters() map[string]*h.HyperLogLog {
	return this.hll
}
