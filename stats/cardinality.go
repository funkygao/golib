package stats

import (
	"encoding/gob"
	"errors"
	h "github.com/funkygao/hyperloglog"
	"hash/fnv"
	"math"
	"os"
	"time"
)

var ErrUnkownType = errors.New("unknown type")

type HllCounter struct {
	Hll       *h.HyperLogLog
	StartedAt time.Time
}

func (this *HllCounter) reset() {
	this.Hll.Reset()
	this.StartedAt = time.Now()
}

func (this *HllCounter) count() uint64 {
	return this.Hll.Count()
}

type CardinalityCounter struct {
	M  uint
	Hc map[string]*HllCounter
}

func NewCardinalityCounter() *CardinalityCounter {
	c := &CardinalityCounter{M: uint(math.Pow(2, float64(18)))} // 2^8KB mem space
	c.Hc = make(map[string]*HllCounter)
	return c
}

func (this *CardinalityCounter) Dump(fn string) error {
	file, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(*this)
	return err
}

func (this *CardinalityCounter) Load(fn string) error {
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	decoder.Decode(this)
	return nil
}

func (this *CardinalityCounter) Add(key string, data interface{}) (err error) {
	if _, ok := this.Hc[key]; !ok {
		var hll *h.HyperLogLog
		hll, err = h.New(this.M)
		if err != nil {
			return
		}

		this.Hc[key] = &HllCounter{StartedAt: time.Now(), Hll: hll}
	}

	switch data.(type) {
	case string:
		hash := fnv.New32()
		hash.Write([]byte(data.(string)))
		this.Hc[key].Hll.Add(hash.Sum32())
	case int:
		this.Hc[key].Hll.Add(uint32(data.(int)))
	case int16:
		this.Hc[key].Hll.Add(uint32(data.(int16)))
	case int32:
		this.Hc[key].Hll.Add(uint32(data.(int32)))
	case int64:
		this.Hc[key].Hll.Add(uint32(data.(int64)))
	case uint:
		this.Hc[key].Hll.Add(uint32(data.(uint)))
	case uint16:
		this.Hc[key].Hll.Add(uint32(data.(uint16)))
	case uint32:
		this.Hc[key].Hll.Add(data.(uint32))
	case uint64:
		this.Hc[key].Hll.Add(uint32(data.(uint64)))
	default:
		err = ErrUnkownType
	}

	return
}

func (this *CardinalityCounter) Reset(key string) {
	if _, ok := this.Hc[key]; !ok {
		return
	}

	this.Hc[key].reset()
}

func (this *CardinalityCounter) ResetAll() {
	for key, _ := range this.Hc {
		this.Reset(key)
	}
}

func (this *CardinalityCounter) Count(key string) uint64 {
	if _, ok := this.Hc[key]; !ok {
		return 0
	}

	return this.Hc[key].count()
}

func (this *CardinalityCounter) StartedAt(key string) time.Time {
	if _, ok := this.Hc[key]; !ok {
		return time.Time{}
	}

	return this.Hc[key].StartedAt
}

func (this *CardinalityCounter) Categories() []string {
	val := make([]string, 0, 10)
	for key, _ := range this.Hc {
		val = append(val, key)
	}

	return val
}
