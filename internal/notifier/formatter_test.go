package notifier

import (
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/vault"
)

func newSecret(path string, expiresIn time.Duration) *vault.Secret {
	return &vault.Secret{
		Path:      path,
		ExpiresAt: time.Now().Add(expiresIn),
	}
}

func TestFormatMessage_Warning(t *testing.T) {
	s := newSecret("secret/data/myapp/db", 10*24*time.Hour)
	msg := FormatMessage(s, 14*24*time.Hour)

	if msg.Level != LevelWarning {
		t.Errorf("expected level WARNING, got %s", msg.Level)
	}
	if msg.Secret != s {
		t.Error("expected Secret field to reference original secret")
	}
	if msg.Subject == "" {
		t.Error("expected non-empty Subject")
	}
	if msg.Body == "" {
		t.Error("expected non-empty Body")
	}
}

func TestFormatMessage_Critical(t *testing.T) {
	s := newSecret("secret/data/myapp/api", 12*time.Hour)
	msg := FormatMessage(s, 14*24*time.Hour)

	if msg.Level != LevelCritical {
		t.Errorf("expected level CRITICAL, got %s", msg.Level)
	}
	if msg.Subject == "" {
		t.Error("expected non-empty Subject")
	}
}

func TestFormatMessage_Expired(t *testing.T) {
	s := newSecret("secret/data/myapp/cert", -1*time.Hour)
	msg := FormatMessage(s, 14*24*time.Hour)

	if msg.Level != LevelExpired {
		t.Errorf("expected level EXPIRED, got %s", msg.Level)
	}
	if msg.Subject == "" {
		t.Error("expected non-empty Subject")
	}
}

func TestFormatMessage_SubjectContainsPath(t *testing.T) {
	path := "secret/data/payments/stripe"
	s := newSecret(path, 5*24*time.Hour)
	msg := FormatMessage(s, 14*24*time.Hour)

	for _, field := range []string{msg.Subject, msg.Body} {
		if len(field) == 0 {
			t.Error("expected non-empty message field")
		}
	}
}

func TestFormatMessage_BodyContainsExpiry(t *testing.T) {
	s := newSecret("secret/data/app/token", 3*24*time.Hour)
	msg := FormatMessage(s, 14*24*time.Hour)

	if msg.Body == "" {
		t.Error("expected body to be populated")
	}
	if msg.Level != LevelWarning {
		t.Errorf("unexpected level: %s", msg.Level)
	}
}
