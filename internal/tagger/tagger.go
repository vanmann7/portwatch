// Package tagger assigns human-readable tags to port events based on
// well-known service mappings and user-defined rules.
package tagger

import "fmt"

// well-known maps port numbers to service tag names.
var wellKnown = map[int]string{
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
	8443: "https-alt",
	27017: "mongodb",
}

// Tagger resolves a tag string for a given port number.
type Tagger struct {
	overrides map[int]string
}

// New returns a Tagger with an optional set of user-defined overrides.
// Overrides take precedence over well-known mappings.
func New(overrides map[int]string) *Tagger {
	if overrides == nil {
		overrides = make(map[int]string)
	}
	return &Tagger{overrides: overrides}
}

// Tag returns the tag associated with port. If no mapping exists the tag
// is the string representation of the port number prefixed with "port-".
func (t *Tagger) Tag(port int) string {
	if tag, ok := t.overrides[port]; ok {
		return tag
	}
	if tag, ok := wellKnown[port]; ok {
		return tag
	}
	return fmt.Sprintf("port-%d", port)
}

// TagAll returns a map of port → tag for every port in the provided slice.
func (t *Tagger) TagAll(ports []int) map[int]string {
	out := make(map[int]string, len(ports))
	for _, p := range ports {
		out[p] = t.Tag(p)
	}
	return out
}
