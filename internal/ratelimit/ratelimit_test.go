package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllowConsumesTokens(t *testing.T) {
	l := ratelimit.New(3, time.Second)

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}

	if l.Allow() {
		t.Fatal("expected Allow()=false after exhausting tokens")
	}
}

func TestResetRestoresTokens(t *testing.T) {
	l := ratelimit.New(2, time.Second)
	l.Allow()
	l.Allow()

	if l.Allow() {
		t.Fatal("tokens should be exhausted before reset")
	}

	l.Reset()

	if !l.Allow() {
		t.Fatal("expected Allow()=true after reset")
	}
}

func TestRefillAfterInterval(t *testing.T) {
	l := ratelimit.New(2, 50*time.Millisecond)
	l.Allow()
	l.Allow()

	if l.Allow() {
		t.Fatal("should be rate-limited before interval")
	}

	time.Sleep(60 * time.Millisecond)

	if !l.Allow() {
		t.Fatal("expected token refill after interval elapsed")
	}
}

func TestNewLimiterFullTokens(t *testing.T) {
	l := ratelimit.New(5, time.Second)
	count := 0
	for l.Allow() {
		count++
	}
	if count != 5 {
		t.Fatalf("expected 5 allowed, got %d", count)
	}
}
