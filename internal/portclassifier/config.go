package portclassifier

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseOverrides reads lines of the form "<port>:<category>" from r and
// returns a map suitable for passing to New. Blank lines and lines starting
// with '#' are ignored.
func ParseOverrides(r io.Reader) (map[uint16]string, error) {
	out := make(map[uint16]string)
	scanner := bufio.NewScanner(r)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("portclassifier: line %d: expected <port>:<category>, got %q", lineNo, line)
		}
		rawPort := strings.TrimSpace(parts[0])
		rawCat := strings.TrimSpace(parts[1])
		n, err := strconv.ParseUint(rawPort, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("portclassifier: line %d: invalid port %q: %w", lineNo, rawPort, err)
		}
		out[uint16(n)] = rawCat
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("portclassifier: reading overrides: %w", err)
	}
	return out, nil
}
