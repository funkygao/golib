package ip

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestIpConvert(t *testing.T) {
	addr := "192.34.45.4"
	assert.Equal(t, uint32(0xc0222d04), Ip2long(addr))
	assert.Equal(t, "192.34.45.4", Long2ip(0xc0222d04))
}
