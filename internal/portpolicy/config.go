package portpolicy

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseRules reads lines of the form:
//
//	<action> <name> <port|min-max>
//
// Lines beginning with '#' and blank lines are ignored.
func ParseRules(r io.Reader) ([]Rule, error) {
	var rules []Rule
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("portpolicy: malformed line: %q", line)
		}
		var action Action
		switch strings.ToLower(parts[0]) {
		case "allow":
			action = Allow
		case "deny":
			action = Deny
		default:
			return nil, fmt.Errorf("portpolicy: unknown action %q", parts[0])
		}
		name := parts[1]
		min, max, err := parsePortRange(parts[2])
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{Name: name, Action: action, Min: min, Max: max})
	}
	return rules, sc.Err()
}

func parsePortRange(s string) (int, int, error) {
	if idx := strings.IndexByte(s, '-'); idx >= 0 {
		lo, err := strconv.Atoi(s[:idx])
		if err != nil {
			return 0, 0, fmt.Errorf("portpolicy: invalid port %q", s[:idx])
		}
		hi, err := strconv.Atoi(s[idx+1:])
		if err != nil {
			return 0, 0, fmt.Errorf("portpolicy: invalid port %q", s[idx+1:])
		}
		return lo, hi, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, 0, fmt.Errorf("portpolicy: invalid port %q", s)
	}
	return v, v, nil
}
