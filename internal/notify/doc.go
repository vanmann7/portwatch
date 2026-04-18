// Package notify provides notification channel abstractions for portwatch.
//
// It defines the Channel interface and a Dispatcher that fans out
// notifications to multiple channels simultaneously. The built-in
// StdoutChannel writes timestamped messages to any io.Writer.
package notify
