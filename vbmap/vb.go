// VBucketMap
// learned from couchbase
package vbmap

type VBucketMap struct {
	slots     int      // fixed number
	nodes     []string // server list
	bucketMap []int
}

// If slots=0, use default slots num
func New(slots int) (this *VBucketMap) {
	if slots <= 0 {
		slots = SLOTS_DEFAULT
	}

	this = &VBucketMap{slots: slots}
	this.bucketMap = make([]int, slots)

	return
}

// VBHash finds the vbucket for the given key.
func (this *VBucketMap) Hash(key string) uint32 {
	crc := uint32(0xffffffff)
	for x := 0; x < len(key); x++ {
		crc = (crc >> 8) ^ crc32tab[(uint64(crc)^uint64(key[x]))&0xff]
	}
	return ((^crc) >> 16) & 0x7fff & (uint32(this.slots) - 1)
}

// Set server list, each node being 'host:port' alike string
func (this *VBucketMap) SetNodes(nodes []string) *VBucketMap {
	this.nodes = nodes
	n := len(nodes)
	// distribute nodes evenly across all slots
	for i := 0; i < this.slots; i++ {
		this.bucketMap[i] = i % n
	}

	return this
}

func (this *VBucketMap) Node(key string) string {
	hash := int64(this.Hash(key))
	slot := hash % int64(this.slots)
	return this.nodes[this.bucketMap[slot]]
}
