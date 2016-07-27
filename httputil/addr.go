package httputil

import (
	"net/http"
	"strings"
)

// HttpRealIp returns ip only, without remote port.
func HttpRealIp(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For") // client_ip,proxy_ip,proxy_ip,...
	if forwardedFor == "" {
		// directly connected
		idx := strings.LastIndex(r.RemoteAddr, ":")
		if idx == -1 {
			return r.RemoteAddr
		}

		return r.RemoteAddr[:idx]
	}

	return forwardedFor // FIXME forwardedFor might be comma seperated ip list, but here for performance ignore it
}
