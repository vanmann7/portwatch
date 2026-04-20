// Package aggregator provides a time-windowed event aggregator for port
// change notifications.
//
// Events (opened / closed port actions) are collected over a configurable
// window duration. When the window expires, a single [Summary] containing
// all unique opened and closed ports is emitted on the output channel.
//
// This reduces downstream alert noise when many ports change simultaneously,
// for example during a service restart or a bulk deployment.
package aggregator
