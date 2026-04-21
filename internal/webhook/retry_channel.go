package webhook

import (
	"context"
	"fmt"
	"time"
)

// RetryConfig controls retry behaviour for the RetryChannel.
type RetryConfig struct {
	// MaxAttempts is the total number of attempts (including the first). Default: 3.
	MaxAttempts int
	// BaseDelay is the initial back-off delay. Default: 500ms.
	BaseDelay time.Duration
}

func (rc RetryConfig) withDefaults() RetryConfig {
	if rc.MaxAttempts <= 0 {
		rc.MaxAttempts = 3
	}
	if rc.BaseDelay <= 0 {
		rc.BaseDelay = 500 * time.Millisecond
	}
	return rc
}

// RetryChannel wraps a Channel and retries failed sends with exponential back-off.
type RetryChannel struct {
	inner *Channel
	cfg   RetryConfig
}

// NewRetryChannel creates a RetryChannel wrapping inner with the given retry policy.
func NewRetryChannel(inner *Channel, cfg RetryConfig) *RetryChannel {
	return &RetryChannel{inner: inner, cfg: cfg.withDefaults()}
}

// Send attempts to deliver p, retrying on failure up to MaxAttempts times.
// Back-off doubles after each failure. Context cancellation aborts early.
func (r *RetryChannel) Send(ctx context.Context, p Payload) error {
	delay := r.cfg.BaseDelay
	var lastErr error
	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("webhook retry: context done after %d attempt(s): %w", attempt-1, err)
		}
		if err := r.inner.Send(ctx, p); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if attempt < r.cfg.MaxAttempts {
			select {
			case <-ctx.Done():
				return fmt.Errorf("webhook retry: context cancelled during back-off: %w", ctx.Err())
			case <-time.After(delay):
			}
			delay *= 2
		}
	}
	return fmt.Errorf("webhook retry: all %d attempt(s) failed: %w", r.cfg.MaxAttempts, lastErr)
}
