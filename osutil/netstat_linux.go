// +build linux,!appengine

package osutil

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func init() {
	netstatFn = netstatLinux
}

func netstatLinux() map[string]int64 {
	b, err := ioutil.ReadFile("/proc/net/netstat")
	if err != nil {
		return nil
	}

	headerFound := false
	var fields []string
	var values []int64
	r := make(map[string]int64)
	for _, line := range strings.Split(string(b), "\n") {
		if !strings.HasPrefix(line, "TcpExt:") {
			// we only parse TcpExt section
			continue
		}

		tuples := strings.Fields(line)
		if !headerFound {
			headerFound = true
			for _, field := range tuples[1:] {
				fields = append(fields, field)
			}
		} else {
			for _, val := range tuples[1:] {
				value, _ := strconv.ParseInt(val, 10, 64)
				values = append(values, value)
			}
			break
		}
	}

	if len(fields) != len(values) {
		// should never happen
		return nil
	}

	for i, field := range fields {
		r[field] = values[i]
	}

	return r
}
