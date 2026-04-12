package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewSpikeNotifier_Valid(t *testing.T) {
	n, err := notifier.NewSpikeNotifier("https://hooks.spike.sh/abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewSpikeNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewSpikeNotifier("")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestSpikeNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewSpikeNotifier(server.URL)
	secret := vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !received {
		t.Fatal("expected server to receive request")
	}
}

func TestSpikeNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewSpikeNotifier(server.URL)
	secret := vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSpikeNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notifier.NewSpikeNotifier(server.URL)
	secret := vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for server error response")
	}
}
