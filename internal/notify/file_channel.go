package notify

import (
	"fmt"
	"os"
	"time"
)

// FileChannel appends notifications to a log file.
type FileChannel struct {
	path string
}

// NewFileChannel creates a FileChannel that appends to the file at path.
func NewFileChannel(path string) *FileChannel {
	return &FileChannel{path: path}
}

// Send appends a formatted log line to the file.
func (f *FileChannel) Send(subject, body string) error {
	file, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("notify: open log file: %w", err)
	}
	defer file.Close()
	line := fmt.Sprintf("[%s] %s: %s\n", time.Now().Format(time.RFC3339), subject, body)
	_, err = file.WriteString(line)
	return err
}
