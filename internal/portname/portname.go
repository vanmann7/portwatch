// Package portname maps port numbers to human-readable service names.
// It provides a lightweight lookup layer with support for custom overrides
// and falls back to a curated built-in table for well-known ports.
package portname

const defaultName = "unknown"

// Resolver maps port numbers to service names.
type Resolver struct {
	overrides map[int]string
}

// builtinNames contains IANA-assigned service names for common ports.
var builtinNames = map[int]string{
	20:   "ftp-data",
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	67:   "dhcp",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	194:  "irc",
	389:  "ldap",
	443:  "https",
	465:  "smtps",
	514:  "syslog",
	587:  "submission",
	636:  "ldaps",
	993:  "imaps",
	995:  "pop3s",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// New returns a Resolver with optional custom overrides.
// Overrides take precedence over the built-in table.
func New(overrides map[int]string) *Resolver {
	if overrides == nil {
		overrides = make(map[int]string)
	}
	return &Resolver{overrides: overrides}
}

// Lookup returns the service name for port. If no match is found in
// overrides or the built-in table, it returns "unknown".
func (r *Resolver) Lookup(port int) string {
	if name, ok := r.overrides[port]; ok {
		return name
	}
	if name, ok := builtinNames[port]; ok {
		return name
	}
	return defaultName
}

// IsKnown reports whether the port has a named service entry.
func (r *Resolver) IsKnown(port int) bool {
	return r.Lookup(port) != defaultName
}
