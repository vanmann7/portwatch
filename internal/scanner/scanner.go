package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents an open port with its protocol and state.
type Port struct {
	Protocol string
	Number   int
	Address  string
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// New creates a Scanner with sensible defaults.
func New(host string) *Scanner {
	return &Scanner{
		Host:    host,
		Timeout: 500 * time.Millisecond,
	}
}

// Scan checks the given port range and returns open ports.
func (s *Scanner) Scan(startPort, endPort int) ([]Port, error) {
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", startPort, endPort)
	}

	var open []Port
	for port := startPort; port <= endPort; port++ {
		addr := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err == nil {
			conn.Close()
			open = append(open, Port{
				Protocol: "tcp",
				Number:   port,
				Address:  addr,
			})
		}
	}
	return open, nil
}
