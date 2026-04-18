// Package healthcheck provides a simple HTTP health endpoint
// for the portwatch daemon, exposing last scan time and metrics.
package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

// Server exposes a lightweight HTTP health endpoint.
type Server struct {
	tracker *metrics.Tracker
	addr    string
	server  *http.Server
}

// Response is the JSON payload returned by the health endpoint.
type Response struct {
	Status      string    `json:"status"`
	LastScanAt  time.Time `json:"last_scan_at"`
	TotalScans  int       `json:"total_scans"`
	OpenPorts   int       `json:"open_ports"`
	AlertsSent  int       `json:"alerts_sent"`
}

// New creates a new health check Server bound to addr.
func New(addr string, tracker *metrics.Tracker) *Server {
	s := &Server{addr: addr, tracker: tracker}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return s
}

// Start begins listening in a background goroutine.
func (s *Server) Start() error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		return fmt.Errorf("healthcheck: %w", err)
	case <-time.After(50 * time.Millisecond):
		return nil
	}
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop() error {
	return s.server.Close()
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	snap := s.tracker.Snapshot()
	resp := Response{
		Status:     "ok",
		LastScanAt: snap.LastScanAt,
		TotalScans: snap.TotalScans,
		OpenPorts:  snap.OpenPorts,
		AlertsSent: snap.AlertsSent,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
