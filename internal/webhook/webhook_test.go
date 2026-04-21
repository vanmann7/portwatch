package webhook_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/webhook"
)

func TestSendPostsJSON(t *testing.T) {
	var got webhook.Payload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := webhook.New(webhook.Config{URL: srv.URL})
	p := webhook.Payload{Event: "opened", Port: 8080, Service: "http-alt"}
	if err := c.Send(context.Background(), p); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got.Port != 8080 || got.Event != "opened" || got.Service != "http-alt" {
		t.Errorf("unexpected payload: %+v", got)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSendIncludesSecret(t *testing.T) {
	var secret string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret = r.Header.Get("X-Portwatch-Secret")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := webhook.New(webhook.Config{URL: srv.URL, Secret: "s3cr3t"})
	if err := c.Send(context.Background(), webhook.Payload{Event: "closed", Port: 22}); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if secret != "s3cr3t" {
		t.Errorf("expected secret header, got %q", secret)
	}
}

func TestSendNonOKStatusReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := webhook.New(webhook.Config{URL: srv.URL})
	err := c.Send(context.Background(), webhook.Payload{Event: "opened", Port: 443})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestSendContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := webhook.New(webhook.Config{URL: srv.URL, Timeout: time.Second})
	if err := c.Send(ctx, webhook.Payload{Event: "opened", Port: 80}); err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestSendDefaultTimeout(t *testing.T) {
	c := webhook.New(webhook.Config{URL: "http://127.0.0.1:1"})
	// Just verify it doesn't panic and returns an error quickly.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := c.Send(ctx, webhook.Payload{Event: "opened", Port: 9}); err == nil {
		t.Fatal("expected connection error")
	}
}
