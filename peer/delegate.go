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

// delegate manages gossiped data: the set of members, their type and API port.
type delegate struct {
	mu sync.RWMutex

	broadcast *memberlist.TransmitLimitedQueue
	members   map[string]memberInfo
	myName    string
}

func newDelegate() *delegate {
	return &delegate{
		broadcast: nil,
		members:   map[string]memberInfo{},
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
	d.members[myName] = memberInfo{
		Tags:    myTags,
		APIAddr: apiAddr,
		APIPort: apiPort,
	}
}

func (d *delegate) state() map[string]memberInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// snapshot clone
	r := map[string]memberInfo{}
	for k, v := range d.members {
		r[k] = v
	}
	return r
}

// Implements memberlist.Delegate.
func (d *delegate) NodeMeta(limit int) []byte {
	// memberlist.Node.Meta
	return []byte{} // TODO encode my tags
}

// Implements memberlist.Delegate.
func (d *delegate) NotifyMsg(b []byte) {
	if len(b) == 0 {
		return
	}

	t := messageType(b[0])
	switch t {
	case messageTypeLeave:
	case messageTypeJoin:
		// TODO decodeMessage b[1:]
	case messageTypePushPull:
	case messageTypeUserEvent:
	default:
	}

	var members map[string]memberInfo
	if err := json.Unmarshal(b, &members); err != nil {
		// TODO log err
		return
	}

	d.mu.Lock()
	for k, v := range members {
		// Removing members is handled by NotifyLeave
		d.members[k] = v
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

	//encodeMessage(messageTypePushPull, messagePushPull)

	b, _ := json.Marshal(d.members)
	return b
}

// Implements memberlist.Delegate.
func (d *delegate) MergeRemoteState(b []byte, join bool) {
	if len(b) == 0 {
		return
	}

	if false && messageType(b[0]) != messageTypePushPull {
		return
	}

	d.NotifyMsg(b)
}

// Implements memberlist.EventDelegate.
func (d *delegate) NotifyJoin(n *memberlist.Node) {}

// Implements memberlist.EventDelegate.
func (d *delegate) NotifyLeave(n *memberlist.Node) {
	d.mu.Lock()
	delete(d.members, n.Name)
	d.mu.Unlock()
}

// Implements memberlist.EventDelegate.
func (d *delegate) NotifyUpdate(n *memberlist.Node) {}
