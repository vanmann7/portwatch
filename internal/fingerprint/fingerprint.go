// Package fingerprint provides port fingerprinting by capturing a lightweight
// signature (protocol hint + banner snippet) for an open port.
package fingerprint

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Result holds the fingerprint data collected for a single port.
type Result struct {
	Port    int
	Banner  string
	Protocol string
}

// Fingerprinter probes open ports and returns a Result.
type Fingerprinter struct {
	timeout time.Duration
	maxRead int
}

// New returns a Fingerprinter with the given dial timeout and maximum banner
// bytes to read. Sensible defaults are used when zero values are supplied.
func New(timeout time.Duration, maxRead int) *Fingerprinter {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	if maxRead <= 0 {
		maxRead = 128
	}
	return &Fingerprinter{timeout: timeout, maxRead: maxRead}
}

// Probe connects to the given TCP port on localhost and attempts to read a
// banner. It returns a Result regardless of whether a banner was obtained.
func (f *Fingerprinter) Probe(port int) Result {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", addr, f.timeout)
	if err != nil {
		return Result{Port: port, Protocol: "unknown"}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(f.timeout))
	buf := make([]byte, f.maxRead)
	n, _ := conn.Read(buf)

	banner := strings.TrimSpace(string(buf[:n]))
	protocol := guessProtocol(port, banner)

	return Result{
		Port:     port,
		Banner:   banner,
		Protocol: protocol,
	}
}

// guessProtocol applies simple heuristics to label a protocol.
func guessProtocol(port int, banner string) string {
	switch port {
	case 22:
		return "ssh"
	case 80, 8080, 8000:
		return "http"
	case 443, 8443:
		return "https"
	case 21:
		return "ftp"
	case 25, 587:
		return "smtp"
	case 3306:
		return "mysql"
	case 5432:
		return "postgres"
	}
	b := strings.ToLower(banner)
	switch {
	case strings.HasPrefix(b, "ssh"):
		return "ssh"
	case strings.Contains(b, "http"):
		return "http"
	case strings.HasPrefix(b, "220"):
		return "smtp"
	}
	return "unknown"
}
