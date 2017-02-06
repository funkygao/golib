package peer

import (
	"encoding/json"
	"sync"

	"github.com/hashicorp/memberlist"
)

var (
	_ memberlist.Delegate      = &delegate{}
	_ memberlist.EventDelegate = &delegate{}
)

// delegate manages gossiped data: the set of peers, their type and API port.
type delegate struct {
	mu sync.RWMutex

	broadcast *memberlist.TransmitLimitedQueue
	data      map[string]peerInfo
	myName    string
}

func newDelegate() *delegate {
	return &delegate{
		broadcast: nil,
		data:      map[string]peerInfo{},
	}
}

func (d *delegate) init(myName string, myTags []string, apiAddr string, apiPort int, numNodes func() int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.myName = myName
	d.broadcast = &memberlist.TransmitLimitedQueue{
		NumNodes:       numNodes,
		RetransmitMult: 3,
	}
	d.data[myName] = peerInfo{
		Tags:    myTags,
		APIAddr: apiAddr,
		APIPort: apiPort,
	}
}

func (d *delegate) state() map[string]peerInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// snapshot clone
	r := map[string]peerInfo{}
	for k, v := range d.data {
		r[k] = v
	}
	return r
}

// Implements memberlist.Delegate.
func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{} // no meta-data
}

// Implements memberlist.Delegate.
func (d *delegate) NotifyMsg(b []byte) {
	if len(b) == 0 {
		return
	}

	var data map[string]peerInfo
	if err := json.Unmarshal(b, &data); err != nil {
		// TODO log err
		return
	}

	d.mu.Lock()
	for k, v := range data {
		// Removing data is handled by NotifyLeave
		d.data[k] = v
	}
	d.mu.Unlock()
}

// Implements memberlist.Delegate.
func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	d.mu.RLock()
	d.mu.RUnlock()

	return d.broadcast.GetBroadcasts(overhead, limit)
}

// Implements memberlist.Delegate.
func (d *delegate) LocalState(join bool) []byte {
	d.mu.RLock()
	defer d.mu.RUnlock()

	b, _ := json.Marshal(d.data)
	return b
}

// Implements memberlist.Delegate.
func (d *delegate) MergeRemoteState(b []byte, join bool) {
	d.NotifyMsg(b)
}

// Implements memberlist.EventDelegate.
func (d *delegate) NotifyJoin(n *memberlist.Node) {}

// Implements memberlist.EventDelegate.
func (d *delegate) NotifyLeave(n *memberlist.Node) {
	d.mu.Lock()
	delete(d.data, n.Name)
	d.mu.Unlock()
}

// Implements memberlist.EventDelegate.
func (d *delegate) NotifyUpdate(n *memberlist.Node) {}
