// Package portgroup provides grouping of ports into named logical categories
// such as "web", "database", or "admin" for richer alert context.
package portgroup

import "fmt"

// Group represents a named collection of port numbers.
type Group struct {
	Name  string
	Ports map[int]struct{}
}

// Grouper maps individual ports to their group names.
type Grouper struct {
	groups []*Group
}

// Rule defines a named group and the ports or port ranges that belong to it.
type Rule struct {
	Name  string
	Ports []string // e.g. "80", "443", "8000-8099"
}

// New constructs a Grouper from a slice of Rules.
// Returns an error if any rule contains an invalid port or range.
func New(rules []Rule) (*Grouper, error) {
	g := &Grouper{}
	for _, r := range rules {
		grp := &Group{Name: r.Name, Ports: make(map[int]struct{})}
		for _, p := range r.Ports {
			ports, err := parsePorts(p)
			if err != nil {
				return nil, fmt.Errorf("portgroup: rule %q: %w", r.Name, err)
			}
			for _, port := range ports {
				grp.Ports[port] = struct{}{}
			}
		}
		g.groups = append(g.groups, grp)
	}
	return g, nil
}

// Lookup returns the name of the first group that contains port.
// If no group matches, it returns an empty string.
func (g *Grouper) Lookup(port int) string {
	for _, grp := range g.groups {
		if _, ok := grp.Ports[port]; ok {
			return grp.Name
		}
	}
	return ""
}

// Groups returns a snapshot of all registered group names.
func (g *Grouper) Groups() []string {
	names := make([]string, 0, len(g.groups))
	for _, grp := range g.groups {
		names = append(names, grp.Name)
	}
	return names
}

// parsePorts parses a single port string ("80") or range ("8000-8099").
func parsePorts(s string) ([]int, error) {
	var lo, hi int
	if n, err := fmt.Sscanf(s, "%d-%d", &lo, &hi); n == 2 && err == nil {
		if lo < 1 || hi > 65535 || lo > hi {
			return nil, fmt.Errorf("invalid range %q", s)
		}
		ports := make([]int, 0, hi-lo+1)
		for p := lo; p <= hi; p++ {
			ports = append(ports, p)
		}
		return ports, nil
	}
	var port int
	if _, err := fmt.Sscanf(s, "%d", &port); err != nil {
		return nil, fmt.Errorf("invalid port %q", s)
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("port %d out of range", port)
	}
	return []int{port}, nil
}
