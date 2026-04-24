// Package portclassifier categorises TCP/UDP port numbers into well-known
// ranges: system (0–1023), user (1024–49151), and dynamic (49152–65535).
//
// Callers may supply per-port overrides via a simple text format:
//
//	# comment
//	8080: user
//	60000: dynamic
//
// Use ParseOverrides to read the override file and pass the result to New.
package portclassifier
