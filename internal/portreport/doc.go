// Package portreport assembles a structured Report from a snapshot of open
// ports by combining labels, severity levels, and presence-tracking data.
//
// Typical usage:
//
//	b := portreport.New(myLabeler, mySeverity, myPresence)
//	report := b.Build(openPorts)
//	fmt.Println(report.Summary())
package portreport
