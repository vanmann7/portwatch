// Package severity provides a Classifier that assigns severity levels
// (Info, Warning, Critical) to port numbers based on well-known service
// associations and optional caller-supplied overrides.
//
// Built-in rules:
//   - Critical: privileged and commonly exploited ports (SSH, Telnet, RDP, …)
//   - Warning:  well-known service ports and the IANA registered range (1024–49151)
//   - Info:     everything else (dynamic / private range)
//
// Overrides supplied to New always take precedence over built-in rules,
// allowing operators to tune alerting thresholds for their environment.
package severity
