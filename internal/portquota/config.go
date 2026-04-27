package portquota

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Rule defines a quota override for a specific port.
type Rule struct {
	Port   int
	Max    int
	Window time.Duration
}

// ParseRules parses lines of the form "port:max:window" where window is a Go
// duration string (e.g. "80:10:1m"). Blank lines and lines starting with '#'
// are ignored.
func ParseRules(lines []string) ([]Rule, error) {
	var rules []Rule
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("portquota: invalid rule %q: want port:max:window", line)
		}
		port, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil || port < 1 || port > 65535 {
			return nil, fmt.Errorf("portquota: invalid port in rule %q", line)
		}
		max, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || max < 1 {
			return nil, fmt.Errorf("portquota: invalid max in rule %q", line)
		}
		win, err := time.ParseDuration(strings.TrimSpace(parts[2]))
		if err != nil || win <= 0 {
			return nil, fmt.Errorf("portquota: invalid window in rule %q", line)
		}
		rules = append(rules, Rule{Port: port, Max: max, Window: win})
	}
	return rules, nil
}

// ApplyRules returns a map of per-port Trackers built from the supplied rules.
// Ports not listed in rules are not present in the map; callers should fall
// back to a default Tracker for unlisted ports.
func ApplyRules(rules []Rule) map[int]*Tracker {
	m := make(map[int]*Tracker, len(rules))
	for _, r := range rules {
		m[r.Port] = New(r.Max, r.Window)
	}
	return m
}
