package webhook_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/webhook"
)

func TestRetrySucceedsOnFirstAttempt(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ch := webhook.New(webhook.Config{URL: srv.URL})
	rc := webhook.NewRetryChannel(ch, webhook.RetryConfig{MaxAttempts: 3, BaseDelay: 10 * time.Millisecond})

	if err := rc.Send(context.Background(), webhook.Payload{Event: "opened", Port: 80}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryRetriesOnFailure(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ch := webhook.New(webhook.Config{URL: srv.URL})
	rc := webhook.NewRetryChannel(ch, webhook.RetryConfig{MaxAttempts: 3, BaseDelay: 5 * time.Millisecond})

	if err := rc.Send(context.Background(), webhook.Payload{Event: "closed", Port: 443}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetryExhaustsAttempts(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	ch := webhook.New(webhook.Config{URL: srv.URL})
	rc := webhook.NewRetryChannel(ch, webhook.RetryConfig{MaxAttempts: 2, BaseDelay: 5 * time.Millisecond})

	if err := rc.Send(context.Background(), webhook.Payload{Event: "opened", Port: 8443}); err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestRetryAbortsOnContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	ch := webhook.New(webhook.Config{URL: srv.URL})
	rc := webhook.NewRetryChannel(ch, webhook.RetryConfig{MaxAttempts: 5, BaseDelay: 200 * time.Millisecond})

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	if err := rc.Send(ctx, webhook.Payload{Event: "opened", Port: 9090}); err == nil {
		t.Fatal("expected error on context cancel")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("retry took too long after cancel: %v", elapsed)
	}
}
