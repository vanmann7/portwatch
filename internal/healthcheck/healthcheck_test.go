package healthcheck_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
	"github.com/user/portwatch/internal/metrics"
)

func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freeAddr: %v", err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

func TestHealthReturnsOK(t *testing.T) {
	tracker := metrics.New()
	tracker.Record(5, 1)

	addr := freeAddr(t)
	srv := healthcheck.New(addr, tracker)
	if err := srv.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() { _ = srv.Stop() })

	url := fmt.Sprintf("http://%s/healthz", addr)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body healthcheck.Response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Status != "ok" {
		t.Errorf("status = %q, want ok", body.Status)
	}
	if body.OpenPorts != 5 {
		t.Errorf("open_ports = %d, want 5", body.OpenPorts)
	}
	if body.TotalScans != 1 {
		t.Errorf("total_scans = %d, want 1", body.TotalScans)
	}
}

func TestHealthLastScanAt(t *testing.T) {
	tracker := metrics.New()
	before := time.Now()
	tracker.Record(0, 0)

	addr := freeAddr(t)
	srv := healthcheck.New(addr, tracker)
	if err := srv.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() { _ = srv.Stop() })

	url := fmt.Sprintf("http://%s/healthz", addr)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	var body healthcheck.Response
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body.LastScanAt.Before(before) {
		t.Errorf("last_scan_at %v is before test start %v", body.LastScanAt, before)
	}
}
