// Package scanner provides TCP port scanning functionality for portwatch.
//
// Usage:
//
//	s := scanner.New("127.0.0.1")
//	ports, err := s.Scan(1, 1024)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, p := range ports {
//		fmt.Printf("open: %s/%d\n", p.Protocol, p.Number)
//	}
//
// The scanner performs sequential TCP dial attempts with a configurable
// timeout. It is intentionally simple and suitable for local-host monitoring
// rather than large-scale network surveys.
package scanner
