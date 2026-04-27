package notifier

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newSMTPSecret() *vault.Secret {
	return &vault.Secret{
		Path:      "secret/smtp-test",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
}

func TestNewSMTPNotifier_Valid(t *testing.T) {
	n, err := NewSMTPNotifier("smtp.example.com", "user", "pass", "from@example.com", "to@example.com", 587)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewSMTPNotifier_MissingHost(t *testing.T) {
	_, err := NewSMTPNotifier("", "user", "pass", "from@example.com", "to@example.com", 587)
	if err == nil {
		t.Fatal("expected error for missing host")
	}
}

func TestNewSMTPNotifier_MissingFrom(t *testing.T) {
	_, err := NewSMTPNotifier("smtp.example.com", "user", "pass", "", "to@example.com", 587)
	if err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestNewSMTPNotifier_MissingTo(t *testing.T) {
	_, err := NewSMTPNotifier("smtp.example.com", "user", "pass", "from@example.com", "", 587)
	if err == nil {
		t.Fatal("expected error for missing to")
	}
}

func TestNewSMTPNotifier_DefaultPort(t *testing.T) {
	n, err := NewSMTPNotifier("smtp.example.com", "", "", "from@example.com", "to@example.com", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.port != 587 {
		t.Errorf("expected default port 587, got %d", n.port)
	}
}

func TestSMTPNotifier_ImplementsInterface(t *testing.T) {
	n, _ := NewSMTPNotifier("smtp.example.com", "", "", "from@example.com", "to@example.com", 587)
	var _ Notifier = n
}
