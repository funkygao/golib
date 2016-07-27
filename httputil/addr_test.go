package httputil

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/funkygao/assert"
)

func TestHttpRealIpBehindProxy(t *testing.T) {
	httpReqRaw := strings.TrimSpace(fmt.Sprintf(`
POST /topics/foobar/v1 HTTP/1.1
Host: localhost:9191
User-Agent: Go-http-client/1.1
Content-Length: %d
Content-Type: application/x-www-form-urlencoded
Appid: myappid
Pubkey: mypubkey
X-Forwarded-For: 1.1.1.12
Accept-Encoding: gzip`, 100)) + "\r\n\r\n"
	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpReqRaw)))
	if err != nil {
		t.Fatal(err.Error())
	}

	r.RemoteAddr = "13.55.21.11:19098"
	t.Logf("%s", r.RemoteAddr)
	realIp := HttpRealIp(r)
	assert.Equal(t, "1.1.1.12", realIp)
}

func TestHttpRealIpWithoutProxy(t *testing.T) {
	httpReqRaw := strings.TrimSpace(fmt.Sprintf(`
POST /topics/foobar/v1 HTTP/1.1
Host: localhost:9191
User-Agent: Go-http-client/1.1
Content-Length: %d
Content-Type: application/x-www-form-urlencoded
Appid: myappid
Pubkey: mypubkey
Accept-Encoding: gzip`, 100)) + "\r\n\r\n"
	r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(httpReqRaw)))
	if err != nil {
		t.Fatal(err.Error())
	}

	r.RemoteAddr = "13.55.21.11:19098" // ipv4
	t.Logf("%s", r.RemoteAddr)
	realIp := HttpRealIp(r)
	assert.Equal(t, "13.55.21.11", realIp)

	r.RemoteAddr = "::ffff:10.213.17.181:11640" // ipv6
	realIp = HttpRealIp(r)
	assert.Equal(t, "::ffff:10.213.17.181", realIp)

}
