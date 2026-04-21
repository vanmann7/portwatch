// Package webhook provides an HTTP webhook notification channel that
// POSTs port change events to a configured URL.
package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Config holds the configuration for a webhook channel.
type Config struct {
	// URL is the endpoint to POST events to.
	URL string
	// Secret is an optional value sent in the X-Portwatch-Secret header.
	Secret string
	// Timeout is the per-request HTTP timeout. Defaults to 5 seconds.
	Timeout time.Duration
}

// Payload is the JSON body sent to the webhook endpoint.
type Payload struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"` // "opened" or "closed"
	Port      int       `json:"port"`
	Service   string    `json:"service,omitempty"`
}

// Channel sends port change events to an HTTP webhook endpoint.
type Channel struct {
	cfg    Config
	client *http.Client
}

// New creates a new webhook Channel. If cfg.Timeout is zero it defaults to 5s.
func New(cfg Config) *Channel {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	return &Channel{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

// Send serialises p as JSON and POSTs it to the configured URL.
// The context controls the lifetime of the HTTP request.
func (c *Channel) Send(ctx context.Context, p Payload) error {
	if p.Timestamp.IsZero() {
		p.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.cfg.Secret != "" {
		req.Header.Set("X-Portwatch-Secret", c.cfg.Secret)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, c.cfg.URL)
	}
	return nil
}
