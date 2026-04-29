// Package portreport assembles a human-readable or structured summary
// of the current port-scan state, combining open ports, their labels,
// severities, and uptime information into a single Report value.
package portreport

import (
	"fmt"
	"sort"
	"time"
)

// Entry describes a single open port in a report.
type Entry struct {
	Port     int
	Label    string
	Severity string
	FirstSeen time.Time
	Uptime   time.Duration
}

// Report is a point-in-time snapshot of all tracked open ports.
type Report struct {
	GeneratedAt time.Time
	Entries     []Entry
}

// Builder collects port metadata from pluggable sources and builds a Report.
type Builder struct {
	labeler  Labeler
	severity Severitier
	presence Presencer
	clock    func() time.Time
}

// Labeler returns a human-readable label for a port number.
type Labeler interface {
	Label(port int) string
}

// Severitier returns a severity string (info/warning/critical) for a port.
type Severitier interface {
	Level(port int) string
}

// Presencer returns the first-seen time and current uptime for a port.
type Presencer interface {
	FirstSeen(port int) (time.Time, bool)
	Uptime(port int) (time.Duration, bool)
}

// New creates a Builder with the supplied dependencies.
func New(l Labeler, s Severitier, p Presencer) *Builder {
	return &Builder{labeler: l, severity: s, presence: p, clock: time.Now}
}

// Build constructs a Report for the given set of open ports.
func (b *Builder) Build(ports []int) Report {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	entries := make([]Entry, 0, len(sorted))
	for _, p := range sorted {
		e := Entry{
			Port:     p,
			Label:    b.labeler.Label(p),
			Severity: b.severity.Level(p),
		}
		if fs, ok := b.presence.FirstSeen(p); ok {
			e.FirstSeen = fs
		}
		if up, ok := b.presence.Uptime(p); ok {
			e.Uptime = up
		}
		entries = append(entries, e)
	}
	return Report{GeneratedAt: b.clock(), Entries: entries}
}

// Summary returns a compact text summary of the report.
func (r Report) Summary() string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("portreport: no open ports at %s", r.GeneratedAt.Format(time.RFC3339))
	}
	return fmt.Sprintf("portreport: %d open port(s) at %s",
		len(r.Entries), r.GeneratedAt.Format(time.RFC3339))
}
