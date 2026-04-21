package watchdog

import (
	"context"
	"sync"
	"time"
)

// Status represents the current health state of the watchdog.
type Status int

const (
	StatusHealthy Status = iota
	StatusStale
	StatusDead
)

// Watchdog monitors whether a component is regularly checking in.
// If the component misses its deadline, the watchdog transitions to
// stale and eventually dead states.
type Watchdog struct {
	mu          sync.Mutex
	lastPing    time.Time
	staleAfter  time.Duration
	deadAfter   time.Duration
	onStale     func()
	onDead      func()
	status      Status
}

// New creates a Watchdog. staleAfter and deadAfter define the thresholds
// after which the component is considered stale or dead without a ping.
func New(staleAfter, deadAfter time.Duration, onStale, onDead func()) *Watchdog {
	w := &Watchdog{
		lastPing:   time.Now(),
		staleAfter: staleAfter,
		deadAfter:  deadAfter,
		onStale:    onStale,
		onDead:     onDead,
		status:     StatusHealthy,
	}
	return w
}

// Ping resets the watchdog timer, marking the component as alive.
func (w *Watchdog) Ping() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastPing = time.Now()
	w.status = StatusHealthy
}

// Status returns the current watchdog status.
func (w *Watchdog) Status() Status {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.status
}

// Run starts the watchdog loop. It checks the component at the given
// interval until the context is cancelled.
func (w *Watchdog) Run(ctx context.Context, checkInterval time.Duration) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.check()
		}
	}
}

) check() {
	wtdefer w.mu.Unlock()
	elapsed := time.Since(w.lastPing)
	switch {
	case elapsed w.dead w.status != StatusDead:
		w.status = StatusDead
		if w.onDead != nil {
			go w.onDead()
		}
	case elapsed >= w.staleAfter && w.status == StatusHealthy:
		w.status = StatusStale
		if w.onStale != nil {
			go w.onStale()
		}
	}
}
