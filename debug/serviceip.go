package debug

import (
	"net"
	"sync"
)

var (
	serviceIP     string
	serviceIPOnce sync.Once
)

// GetServiceIP returns the outbound IP address of this process.
// It dials a UDP connection (no data is sent) to resolve the local address,
// then caches the result for the lifetime of the process.
func GetServiceIP() string {
	serviceIPOnce.Do(func() {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return
		}
		defer conn.Close()
		serviceIP = conn.LocalAddr().(*net.UDPAddr).IP.String()
	})
	return serviceIP
}
