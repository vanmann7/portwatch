package scanner

import (
	"net"
	"testing"
)

// startTestListener opens a TCP listener on an OS-assigned port and returns it.
func startTestListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestScanFindsOpenPort(t *testing.T) {
	ln, port := startTestListener(t)
	defer ln.Close()

	s := New("127.0.0.1")
	ports, err := s.Scan(port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(ports))
	}
	if ports[0].Number != port {
		t.Errorf("expected port %d, got %d", port, ports[0].Number)
	}
}

func TestScanInvalidRange(t *testing.T) {
	s := New("127.0.0.1")
	_, err := s.Scan(100, 50)
	if err == nil {
		t.Error("expected error for invalid range, got nil")
	}
}

func TestScanClosedPort(t *testing.T) {
	// Bind then immediately close to get a port that is very likely closed.
	ln, port := startTestListener(t)
	ln.Close()

	s := New("127.0.0.1")
	ports, err := s.Scan(port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected 0 open ports, got %d", len(ports))
	}
}
