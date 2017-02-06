package peer

// memberInfo is the gossiped data.
type memberInfo struct {
	Tags    []string `json:"tags"`
	APIAddr string   `json:"api_addr"`
	APIPort int      `json:"api_port"`
}
