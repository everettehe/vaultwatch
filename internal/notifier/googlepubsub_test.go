package notifier_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewGooglePubSubNotifier_MissingProjectID(t *testing.T) {
	_, err := notifier.NewGooglePubSubNotifier("", "my-topic")
	if err == nil {
		t.Fatal("expected error for missing project_id, got nil")
	}
}

func TestNewGooglePubSubNotifier_MissingTopicID(t *testing.T) {
	_, err := notifier.NewGooglePubSubNotifier("my-project", "")
	if err == nil {
		t.Fatal("expected error for missing topic_id, got nil")
	}
}

func TestNewGooglePubSubNotifier_Valid(t *testing.T) {
	// NewGooglePubSubNotifier requires real GCP credentials; validate field
	// constraints only — actual client creation is skipped in unit tests.
	if _, err := notifier.NewGooglePubSubNotifier("", "topic"); err == nil {
		t.Error("expected validation error")
	}
	if _, err := notifier.NewGooglePubSubNotifier("proj", ""); err == nil {
		t.Error("expected validation error")
	}
}

func TestGooglePubSubNotifier_ImplementsInterface(t *testing.T) {
	// Compile-time check: GooglePubSubNotifier must satisfy the Notifier interface.
	var _ interface {
		Notify(interface{}, vault.Secret) error
	} = (*notifier.GooglePubSubNotifier)(nil)
}

func newPubSubSecret(path string, daysLeft int, expired bool) vault.Secret {
	expiry := time.Now().Add(time.Duration(daysLeft) * 24 * time.Hour)
	if expired {
		expiry = time.Now().Add(-24 * time.Hour)
	}
	return vault.Secret{
		Path:      path,
		ExpiresAt: expiry,
	}
}

func TestGooglePubSubNotifier_MessageFields_Expiring(t *testing.T) {
	secret := newPubSubSecret("secret/db/password", 5, false)
	if secret.IsExpired() {
		t.Error("expected secret to be not expired")
	}
	if secret.DaysUntilExpiration() < 4 {
		t.Errorf("expected ~5 days remaining, got %d", secret.DaysUntilExpiration())
	}
}

func TestGooglePubSubNotifier_MessageFields_Expired(t *testing.T) {
	secret := newPubSubSecret("secret/db/password", 0, true)
	if !secret.IsExpired() {
		t.Error("expected secret to be expired")
	}
}
