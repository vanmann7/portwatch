package portguard

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseRules reads guard rules from r.
// Each non-blank, non-comment line has the form:
//
//	<verdict> <port>[-<port>]
//
// where verdict is allow, deny, or warn.
func ParseRules(r io.Reader) ([]Rule, error) {
	var rules []Rule
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("portguard: malformed line %q", line)
		}
		v, err := parseVerdict(parts[0])
		if err != nil {
			return nil, err
		}
		low, high, err := parsePortRange(parts[1])
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{Low: low, High: high, Verdict: v})
	}
	return rules, scanner.Err()
}

func parseVerdict(s string) (Verdict, error) {
	switch strings.ToLower(s) {
	case "allow":
		return Allow, nil
	case "deny":
		return Deny, nil
	case "warn":
		return Warn, nil
	default:
		return 0, fmt.Errorf("portguard: unknown verdict %q", s)
	}
}

func parsePortRange(s string) (int, int, error) {
	if idx := strings.IndexByte(s, '-'); idx >= 0 {
		lo, err := strconv.Atoi(s[:idx])
		if err != nil {
			return 0, 0, fmt.Errorf("portguard: invalid port %q", s[:idx])
		}
		hi, err := strconv.Atoi(s[idx+1:])
		if err != nil {
			return 0, 0, fmt.Errorf("portguard: invalid port %q", s[idx+1:])
		}
		return lo, hi, nil
	}
	p, err := strconv.Atoi(s)
	if err != nil {
		return 0, 0, fmt.Errorf("portguard: invalid port %q", s)
	}
	return p, p, nil
}
