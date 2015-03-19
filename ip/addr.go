package ip

import (
	"fmt"
	"net"
)

// LocalIpv4Addrs scan all ip addresses with loopback excluded.
func LocalIpv4Addrs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for _, addr := range addrs {
		var ip net.IP
		switch x := addr.(type) {
		case *net.IPNet:
			ip = x.IP
		case *net.IPAddr:
			ip = x.IP
		default:
			err = fmt.Errorf("unknown interface address type for: %+v", x)
			return nil, err
		}

		if ip.IsLoopback() || ip.To4() == nil {
			// loopback excluded, ipv6 excluded
			continue
		}

		ips = append(ips, ip.String())
	}

	return ips, nil
}
