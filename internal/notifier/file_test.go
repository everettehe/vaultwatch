package notifier_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewFileNotifier_EmptyPath(t *testing.T) {
	_, err := notifier.NewFileNotifier("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestNewFileNotifier_InvalidPath(t *testing.T) {
	_, err := notifier.NewFileNotifier("/nonexistent-dir/vaultwatch.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestNewFileNotifier_Valid(t *testing.T) {
	tmp, err := os.CreateTemp("", "vaultwatch-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	n, err := notifier.NewFileNotifier(tmp.Name())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil FileNotifier")
	}
}

func TestFileNotifier_Notify_WritesLine(t *testing.T) {
	tmp, err := os.CreateTemp("", "vaultwatch-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	n, err := notifier.NewFileNotifier(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	contents := string(data)
	if !strings.Contains(contents, "secret/db/password") {
		t.Errorf("expected log to contain secret path, got: %s", contents)
	}
}

func TestFileNotifier_Notify_MultipleWrites(t *testing.T) {
	tmp, err := os.CreateTemp("", "vaultwatch-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	n, _ := notifier.NewFileNotifier(tmp.Name())

	for i := 0; i < 3; i++ {
		secret := vault.Secret{
			Path:      "secret/key",
			ExpiresAt: time.Now().Add(time.Duration(i+1) * 24 * time.Hour),
		}
		if err := n.Notify(secret); err != nil {
			t.Fatalf("Notify[%d] returned error: %v", i, err)
		}
	}

	data, _ := os.ReadFile(tmp.Name())
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 log lines, got %d", len(lines))
	}
}
