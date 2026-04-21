// Package formatter converts port change events into human-readable text
// or machine-readable JSON strings suitable for writing to stdout, log
// files, or downstream notification channels.
//
// Two formats are supported:
//
//	FormatText  – a single-line human-readable string.
//	FormatJSON  – a compact JSON object with RFC3339 timestamp.
//
// Usage:
//
//	f := formatter.New(formatter.FormatJSON)
//	line, err := f.Format(formatter.Event{
//	    Port: 443, Proto: "tcp", Action: "opened", Service: "https",
//	})
package formatter
