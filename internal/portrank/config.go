package portrank

import (
	"fmt"
	"strconv"
	"strings"
)

// RuleEntry holds a raw config rule before parsing.
// Format: "<port>:<rank>"  e.g. "2222:88"
type RuleEntry struct {
	Port int
	Rank int
}

// ParseRules parses a slice of "port:rank" strings into RuleEntry values.
// Lines that are blank or start with '#' are skipped.
// Returns an error describing the first malformed entry encountered.
func ParseRules(lines []string) ([]RuleEntry, error) {
	var out []RuleEntry
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("portrank: invalid rule %q: expected port:rank", line)
		}
		port, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil || port < 1 || port > 65535 {
			return nil, fmt.Errorf("portrank: invalid port in rule %q", line)
		}
		rank, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || rank < 1 || rank > 100 {
			return nil, fmt.Errorf("portrank: invalid rank in rule %q: must be 1–100", line)
		}
		out = append(out, RuleEntry{Port: port, Rank: rank})
	}
	return out, nil
}

// NewFromRules constructs a Ranker from parsed RuleEntry overrides.
func NewFromRules(rules []RuleEntry) *Ranker {
	overrides := make(map[int]int, len(rules))
	for _, e := range rules {
		overrides[e.Port] = e.Rank
	}
	return New(overrides)
}
