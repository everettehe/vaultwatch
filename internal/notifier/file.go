package notifier

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// FileNotifier appends alert messages to a log file on disk.
type FileNotifier struct {
	mu   sync.Mutex
	path string
}

// NewFileNotifier creates a FileNotifier that writes to the given file path.
// The file is created if it does not exist.
func NewFileNotifier(path string) (*FileNotifier, error) {
	if path == "" {
		return nil, fmt.Errorf("file notifier: path must not be empty")
	}

	// Verify we can open/create the file before returning.
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("file notifier: cannot open file %q: %w", path, err)
	}
	f.Close()

	return &FileNotifier{path: path}, nil
}

// Notify appends a timestamped alert line to the configured file.
func (f *FileNotifier) Notify(secret vault.Secret) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	file, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("file notifier: cannot open file: %w", err)
	}
	defer file.Close()

	msg := FormatMessage(secret)
	line := fmt.Sprintf("[%s] %s: %s\n",
		time.Now().UTC().Format(time.RFC3339),
		msg.Subject,
		msg.Body,
	)

	if _, err := file.WriteString(line); err != nil {
		return fmt.Errorf("file notifier: failed to write: %w", err)
	}
	return nil
}
