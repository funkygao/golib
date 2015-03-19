package ip

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestLocalIpv4Addrs(t *testing.T) {
	ips, err := LocalIpv4Addrs()
	assert.Equal(t, nil, err)
	t.Logf("%+v", ips)
	assert.Equal(t, true, len(ips) > 0) // assume all hosts has NIC

	// loopback is excluded
	for _, ip := range ips {
		if ip == "127.0.0.1" {
			t.Errorf("loopback not excluded")
		}
	}
}
