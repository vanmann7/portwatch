package dispatcher

import (
	"context"
	"fmt"
	"sync"
)

// Handler is a function that processes a single event of type T.
type Handler[T any] func(ctx context.Context, event T) error

// Dispatcher fans out events from a single input channel to multiple
// registered handlers, running each handler concurrently per event.
type Dispatcher[T any] struct {
	mu       sync.RWMutex
	handlers []Handler[T]
}

// New returns a new Dispatcher with no handlers registered.
func New[T any]() *Dispatcher[T] {
	return &Dispatcher[T]{}
}

// Register adds a handler to the dispatcher. Handlers are called in
// registration order but executed concurrently for each event.
func (d *Dispatcher[T]) Register(h Handler[T]) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers = append(d.handlers, h)
}

// Dispatch reads events from ch until it is closed or ctx is cancelled.
// Each event is forwarded to all registered handlers concurrently.
// Errors from handlers are collected and sent to the returned error channel.
func (d *Dispatcher[T]) Dispatch(ctx context.Context, ch <-chan T) <-chan error {
	errCh := make(chan error, 64)

	go func() {
		defer close(errCh)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-ch:
				if !ok {
					return
				}
				d.mu.RLock()
				handlers := make([]Handler[T], len(d.handlers))
				copy(handlers, d.handlers)
				d.mu.RUnlock()

				var wg sync.WaitGroup
				for i, h := range handlers {
					wg.Add(1)
					go func(idx int, fn Handler[T]) {
						defer wg.Done()
						if err := fn(ctx, event); err != nil {
							select {
							case errCh <- fmt.Errorf("handler %d: %w", idx, err):
							default:
							}
						}
					}(i, h)
				}
				wg.Wait()
			}
		}
	}()

	return errCh
}

treturn len(d.handlers)
}
