package notifier_test

import (
	"log/syslog"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewSyslogNotifier_DefaultTag(t *testing.T) {
	n, err := notifier.NewSyslogNotifier("", syslog.LOG_LOCAL0)
	if err != nil {
		// syslog may not be available in all CI environments; skip gracefully.
		t.Skipf("syslog not available: %v", err)
	}
	defer n.Close()

	if n == nil {
		t.Fatal("expected non-nil SyslogNotifier")
	}
}

func TestNewSyslogNotifier_CustomTag(t *testing.T) {
	n, err := notifier.NewSyslogNotifier("myapp", syslog.LOG_DAEMON)
	if err != nil {
		t.Skipf("syslog not available: %v", err)
	}
	defer n.Close()

	if n == nil {
		t.Fatal("expected non-nil SyslogNotifier")
	}
}

func TestSyslogNotifier_Notify_Expiring(t *testing.T) {
	n, err := notifier.NewSyslogNotifier("vaultwatch-test", syslog.LOG_LOCAL0)
	if err != nil {
		t.Skipf("syslog not available: %v", err)
	}
	defer n.Close()

	secret := vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestSyslogNotifier_Notify_Expired(t *testing.T) {
	n, err := notifier.NewSyslogNotifier("vaultwatch-test", syslog.LOG_LOCAL0)
	if err != nil {
		t.Skipf("syslog not available: %v", err)
	}
	defer n.Close()

	secret := vault.Secret{
		Path:      "secret/api/key",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error for expired secret, got: %v", err)
	}
}
