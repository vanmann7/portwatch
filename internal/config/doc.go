// Package config provides loading, validation, and default values for
// portwatch runtime configuration.
//
// Configuration can be supplied via a JSON file. Fields not present in
// the file fall back to the defaults returned by [Default].
//
// Example JSON:
//
//	{
//	  "port_range": {"start": 1, "end": 9999},
//	  "interval":   60000000000,
//	  "state_file": "/var/lib/portwatch/state",
//	  "log_file":   "/var/log/portwatch.log"
//	}
package config
