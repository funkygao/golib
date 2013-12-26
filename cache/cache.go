package cache

// A Key may be any value that is comparable.
// See http://golang.org/ref/spec#Comparison_operators
type Key interface{}

type entry struct {
	key   Key
	value interface{}
}

type Cacheable interface {
	Set(key Key, value interface{})
	Get(key Key) (value interface{}, ok bool)
	Del(key Key)
}

type HasLength interface {
	Len() int
}
