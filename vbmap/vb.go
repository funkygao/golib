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

func (this *VBucketMap) Node(key int64) string {
	slot := key % int64(this.slots)
	return this.nodes[this.bucketMap[slot]]
}
