// Package resolve maps port numbers to well-known service names.
package resolve

import (
	"fmt"
	"net"
)

// Resolver looks up service names for port numbers.
type Resolver struct {
	cache map[uint16]string
	overrides map[uint16]string
}

// New returns a Resolver with an optional map of user-defined overrides.
func New(overrides map[uint16]string) *Resolver {
	if overrides == nil {
		overrides = make(map[uint16]string)
	}
	return &Resolver{
		cache:     make(map[uint16]string),
		overrides: overrides,
	}
}

// Name returns a human-readable label for the given port.
// It checks user overrides first, then the local cache, then the OS service
// database via net.LookupPort. Falls back to "unknown" on failure.
func (r *Resolver) Name(port uint16) string {
	if name, ok := r.overrides[port]; ok {
		return name
	}
	if name, ok := r.cache[port]; ok {
		return name
	}
	name := r.lookup(port)
	r.cache[port] = name
	return name
}

func (r *Resolver) lookup(port uint16) string {
	// net.LookupPort is the inverse direction; use the service file via
	// a raw lookup through the IANA list embedded in the stdlib.
	if svc, err := net.LookupAddr(fmt.Sprintf("%d", port)); err == nil && len(svc) > 0 {
		return svc[0]
	}
	// Fallback: check a small built-in table for common ports.
	if name, ok := commonPorts[port]; ok {
		return name
	}
	return "unknown"
}

// Reset clears the internal lookup cache.
func (r *Resolver) Reset() {
	r.cache = make(map[uint16]string)
}

// commonPorts is a small static table for frequently seen ports.
var commonPorts = map[uint16]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	27017: "mongodb",
}
