package audit

import (
	"fmt"
	"os"
)

// FileLogger wraps Logger with an underlying file handle.
type FileLogger struct {
	*Logger
	f *os.File
}

// NewFileLogger opens or creates path for append-only audit logging.
func NewFileLogger(path string) (*FileLogger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return nil, fmt.Errorf("audit: open file: %w", err)
	}
	return &FileLogger{Logger: New(f), f: f}, nil
}

// Close closes the underlying file.
func (fl *FileLogger) Close() error {
	return fl.f.Close()
}
