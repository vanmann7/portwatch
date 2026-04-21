package labeler

import (
	"fmt"
	"strings"
)

// RuleConfig holds the raw string representation of a labeling rule
// as read from a config file, e.g. "8000-8100:my-service".
type RuleConfig struct {
	Spec  string
	Label string
}

// ParseRules converts a slice of "spec:label" strings into Rule values.
// Spec may be a single port ("80") or a range ("8000-8100").
func ParseRules(entries []string) ([]Rule, error) {
	out := make([]Rule, 0, len(entries))
	for _, entry := range entries {
		r, err := parseEntry(entry)
		if err != nil {
			return nil, fmt.Errorf("labeler: %w", err)
		}
		out = append(out, r)
	}
	return out, nil
}

func parseEntry(s string) (Rule, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 || parts[1] == "" {
		return Rule{}, fmt.Errorf("invalid rule %q: expected <port|range>:<label>", s)
	}
	spec, label := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

	var low, high int
	if idx := strings.Index(spec, "-"); idx >= 0 {
		if _, err := fmt.Sscanf(spec[:idx], "%d", &low); err != nil {
			return Rule{}, fmt.Errorf("invalid range low in %q", s)
		}
		if _, err := fmt.Sscanf(spec[idx+1:], "%d", &high); err != nil {
			return Rule{}, fmt.Errorf("invalid range high in %q", s)
		}
	} else {
		if _, err := fmt.Sscanf(spec, "%d", &low); err != nil {
			return Rule{}, fmt.Errorf("invalid port in %q", s)
		}
		high = low
	}
	if low < 1 || high > 65535 || low > high {
		return Rule{}, fmt.Errorf("port range %d-%d out of bounds", low, high)
	}
	return Rule{Low: low, High: high, Label: label}, nil
}
