package ip

import (
	"net"
	"strings"
)

// 127.0.0.1 excluded
func LocalIpAddrs() (ips []string) {
	info, _ := net.InterfaceAddrs()
	ips = make([]string, 0)
	for _, addr := range info {
		ip := strings.Split(addr.String(), "/")[0]
		if !strings.HasPrefix(ip, "127.0") {
			ips = append(ips, ip)
		}
	}

	return
}
