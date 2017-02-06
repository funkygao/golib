package peer

// peerInfo is the gossiped data.
type peerInfo struct {
	Tags    []string `json:"tags"`
	APIAddr string   `json:"api_addr"`
	APIPort int      `json:"api_port"`
}
