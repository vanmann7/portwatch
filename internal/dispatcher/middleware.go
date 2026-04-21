package dispatcher

import (
	"context"
	"log"
	"time"
)

// Middleware wraps a Handler with additional behaviour.
type Middleware[T any] func(Handler[T]) Handler[T]

// Chain applies a sequence of middleware to a base handler, outermost first.
func Chain[T any](base Handler[T], mw ...Middleware[T]) Handler[T] {
	for i := len(mw) - 1; i >= 0; i-- {
		base = mw[i](base)
	}
	return base
}

// WithLogging returns a Middleware that logs each event dispatch and any
// resulting error using the standard logger.
func WithLogging[T any](label string) Middleware[T] {
	return func(next Handler[T]) Handler[T] {
		return func(ctx context.Context, event T) error {
			start := time.Now()
			err := next(ctx, event)
			if err != nil {
				log.Printf("[dispatcher] %s error after %s: %v", label, time.Since(start), err)
			} else {
				log.Printf("[dispatcher] %s ok (%s)", label, time.Since(start))
			}
			return err
		}
	}
}

// WithRecover returns a Middleware that recovers from panics inside a handler
// and converts them into errors so the dispatcher can continue running.
func WithRecover[T any]() Middleware[T] {
	return func(next Handler[T]) Handler[T] {
		return func(ctx context.Context, event T) (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[dispatcher] recovered panic: %v", r)
					err = context.DeadlineExceeded // sentinel; callers check Is
				}
			}()
			return next(ctx, event)
		}
	}
}
