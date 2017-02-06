package peer

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hashicorp/memberlist"
)

// Peer represents this node in the gossip cluster.
type Peer struct {
	m *memberlist.Memberlist
	d *delegate
}

// New creates or joins a gossip cluster with the seed nodes.
// We will listen for cluster communications on the given addr:port.
// We advertise a HTTP API, reachable on apiPort.
func New(addr string, port int, tags []string, seeds []string, apiPort int) (*Peer, error) {
	d := newDelegate()

	cf := memberlist.DefaultLANConfig()
	host, _ := os.Hostname()
	cf.Name = fmt.Sprintf("%s-%d", host, os.Getpid())
	cf.BindAddr = addr
	cf.BindPort = port
	cf.LogOutput = ioutil.Discard
	cf.Delegate = d
	cf.Events = d

	ml, err := memberlist.Create(cf)
	if err != nil {
		return nil, err
	}

	d.init(cf.Name, tags, ml.LocalNode().Addr.String(), apiPort, ml.NumMembers)
	ml.Join(seeds)
	if len(seeds) > 0 {
		go func(d time.Duration) {
			for range time.Tick(d) {
				if n := ml.NumMembers(); n <= 1 {
					// warning TODO
				}
			}
		}(time.Second * 10)
	}

	return &Peer{
		m: ml,
		d: d,
	}, nil
}

// Name returns the uniq ID of this node in the cluster.
func (p *Peer) Name() string {
	return p.m.LocalNode().Name
}

// Leave the cluster, waiting up to timeout.
func (p *Peer) Leave(timeout time.Duration) error {
	return p.m.Leave(timeout)
}

// ClusterSize returns the total size of the cluster from this node's perspective.
func (p *Peer) ClusterSize() int {
	return p.m.NumMembers()
}

// State returns a JSON-serializable dump of cluster state.
// Useful for debugging.
func (p *Peer) State() map[string]interface{} {
	return map[string]interface{}{
		"self":     p.m.LocalNode(),
		"members":  p.m.Members(),
		"size":     p.ClusterSize(),
		"delegate": p.d.state(),
	}
}
