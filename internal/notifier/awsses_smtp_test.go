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
	n, err := NewSMTPNotifier("smtp.example.com", 587, "user", "pass", "from@example.com", []string{"to@example.com"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewSMTPNotifier_MissingHost(t *testing.T) {
	_, err := NewSMTPNotifier("", 587, "user", "pass", "from@example.com", []string{"to@example.com"})
	if err == nil {
		t.Fatal("expected error for missing host")
	}
}

func TestNewSMTPNotifier_MissingFrom(t *testing.T) {
	_, err := NewSMTPNotifier("smtp.example.com", 587, "user", "pass", "", []string{"to@example.com"})
	if err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestNewSMTPNotifier_MissingTo(t *testing.T) {
	_, err := NewSMTPNotifier("smtp.example.com", 587, "user", "pass", "from@example.com", []string{})
	if err == nil {
		t.Fatal("expected error for missing to")
	}
}

func TestSMTPNotifier_ImplementsInterface(t *testing.T) {
	n, err := NewSMTPNotifier("smtp.example.com", 587, "", "", "from@example.com", []string{"to@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var _ Notifier = n
}
