package fingerprint_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/fingerprint"
)

// startBannerServer opens a TCP listener that writes banner on each connection.
func startBannerServer(t *testing.T, banner string) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_, _ = conn.Write([]byte(banner))
			conn.Close()
		}
	}()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port
}

func TestProbeReturnsBanner(t *testing.T) {
	port := startBannerServer(t, "SSH-2.0-OpenSSH_9.0")
	f := fingerprint.New(time.Second, 64)
	res := f.Probe(port)

	if res.Port != port {
		t.Errorf("port: got %d, want %d", res.Port, port)
	}
	if res.Banner == "" {
		t.Error("expected non-empty banner")
	}
	if res.Protocol != "ssh" {
		t.Errorf("protocol: got %q, want \"ssh\"", res.Protocol)
	}
}

func TestProbeClosedPortReturnsUnknown(t *testing.T) {
	// Port 1 is almost certainly closed in test environments.
	f := fingerprint.New(200*time.Millisecond, 64)
	res := f.Probe(1)

	if res.Protocol != "unknown" {
		t.Errorf("expected unknown protocol for closed port, got %q", res.Protocol)
	}
	if res.Banner != "" {
		t.Errorf("expected empty banner for closed port, got %q", res.Banner)
	}
}

func TestProbeHTTPBanner(t *testing.T) {
	port := startBannerServer(t, "HTTP/1.1 200 OK")
	f := fingerprint.New(time.Second, 128)
	res := f.Probe(port)

	if res.Protocol != "http" {
		t.Errorf("protocol: got %q, want \"http\"", res.Protocol)
	}
}

func TestProbeFTPBanner(t *testing.T) {
	port := startBannerServer(t, "220 FTP Server ready")
	f := fingerprint.New(time.Second, 128)
	res := f.Probe(port)

	if res.Protocol != "ftp" {
		t.Errorf("protocol: got %q, want \"ftp\"", res.Protocol)
	}
	if res.Banner == "" {
		t.Error("expected non-empty banner for FTP server")
	}
}

func TestDefaultsApplied(t *testing.T) {
	// Zero values should not panic and should use sensible defaults.
	f := fingerprint.New(0, 0)
	if f == nil {
		t.Fatal("expected non-nil Fingerprinter")
	}
}
