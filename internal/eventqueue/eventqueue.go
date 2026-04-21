// Package eventqueue provides a bounded, thread-safe FIFO queue for
// port-change events. When the queue is full the oldest entry is evicted
// so that the daemon never blocks the scanner goroutine.
package eventqueue

import "sync"

// Event represents a single port-change notification to be processed.
type Event struct {
	Port   int
	Kind   string // "opened" | "closed"
	Source string // originating scanner label
}

// Queue is a bounded FIFO queue of Events.
type Queue struct {
	mu      sync.Mutex
	items   []Event
	maxSize int
}

// New returns a Queue that holds at most maxSize events.
// If maxSize is less than 1 it is set to 64.
func New(maxSize int) *Queue {
	if maxSize < 1 {
		maxSize = 64
	}
	return &Queue{maxSize: maxSize}
}

// Push appends e to the queue. If the queue is already full the oldest
// event is silently dropped to make room.
func (q *Queue) Push(e Event) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.maxSize {
		q.items = q.items[1:]
	}
	q.items = append(q.items, e)
}

// Pop removes and returns the oldest event and true. If the queue is
// empty it returns the zero Event and false.
func (q *Queue) Pop() (Event, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return Event{}, false
	}
	e := q.items[0]
	q.items = q.items[1:]
	return e, true
}

// Len returns the current number of events in the queue.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Drain removes and returns all events currently in the queue.
func (q *Queue) Drain() []Event {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]Event, len(q.items))
	copy(out, q.items)
	q.items = q.items[:0]
	return out
}
