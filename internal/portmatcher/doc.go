// Package portmatcher implements flexible port-matching rules for use in
// filtering pipelines and alert configurations.
//
// Rules may be expressed as single ports ("443"), inclusive ranges
// ("8000-8080"), or comma-separated combinations ("22,80,443"). Multiple
// rule strings are OR-ed together: a port matches if any rule covers it.
//
// Example:
//
//	m, err := portmatcher.New([]string{"22", "80,443", "8000-8999"})
//	if err != nil { ... }
//	if m.Match(8080) { // true }
package portmatcher
