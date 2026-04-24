// Package portrank assigns a numeric priority rank to ports based on
// their known sensitivity, activity frequency, and optional user overrides.
// Higher rank means higher priority for alerting and processing.
package portrank

import "sync"

// defaultRanks maps well-known ports to a base priority rank (1–100).
var defaultRanks = map[int]int{
	22:   90, // SSH
	23:   95, // Telnet
	80:   60, // HTTP
	443:  65, // HTTPS
	3306: 85, // MySQL
	5432: 85, // PostgreSQL
	6379: 80, // Redis
	27017: 80, // MongoDB
	8080: 55, // HTTP alt
	8443: 60, // HTTPS alt
}

const defaultRank = 10

// Ranker assigns priority ranks to ports.
type Ranker struct {
	mu        sync.RWMutex
	overrides map[int]int
}

// New returns a new Ranker with optional user-defined overrides.
// Override values must be in the range [1, 100]; others are ignored.
func New(overrides map[int]int) *Ranker {
	r := &Ranker{
		overrides: make(map[int]int, len(overrides)),
	}
	for port, rank := range overrides {
		if rank >= 1 && rank <= 100 {
			r.overrides[port] = rank
		}
	}
	return r
}

// Rank returns the priority rank for the given port.
// Override > builtin default > fallback.
func (r *Ranker) Rank(port int) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if v, ok := r.overrides[port]; ok {
		return v
	}
	if v, ok := defaultRanks[port]; ok {
		return v
	}
	return defaultRank
}

// SetOverride adds or updates a runtime override for the given port.
// The rank must be in [1, 100] or the call is a no-op.
func (r *Ranker) SetOverride(port, rank int) {
	if rank < 1 || rank > 100 {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.overrides[port] = rank
}

// RemoveOverride deletes a user override, falling back to builtin or default.
func (r *Ranker) RemoveOverride(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.overrides, port)
}

// Top returns the n ports with the highest ranks from the provided list.
// Ties are broken by port number (lower port wins).
func (r *Ranker) Top(ports []int, n int) []int {
	type entry struct {
		port int
		rank int
	}
	entries := make([]entry, len(ports))
	for i, p := range ports {
		entries[i] = entry{port: p, rank: r.Rank(p)}
	}
	// simple insertion sort — port lists are typically small
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0; j-- {
			a, b := entries[j-1], entries[j]
			if b.rank > a.rank || (b.rank == a.rank && b.port < a.port) {
				entries[j-1], entries[j] = entries[j], entries[j-1]
			} else {
				break
			}
		}
	}
	if n > len(entries) {
		n = len(entries)
	}
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = entries[i].port
	}
	return out
}
