// Package filter implements inclusion rules for port numbers used by portwatch.
//
// A Filter is constructed from a list of range expressions such as "22" or
// "8000-9000". When no rules are provided every port is considered allowed,
// which matches the default portwatch behaviour of scanning the full range
// configured in config.Config.
package filter
