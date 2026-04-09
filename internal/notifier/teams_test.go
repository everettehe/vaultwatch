package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewTeamsNotifier_Valid(t *testing.T) {
	n, err := notifier.NewTeamsNotifier("https://example.webhook.office.com/webhookb2/abc")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewTeamsNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewTeamsNotifier("")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestTeamsNotifier_Notify_ExpiringSoon(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := notifier.NewTeamsNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !called {
		t.Error("expected server to be called")
	}
}

func TestTeamsNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := notifier.NewTeamsNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := vault.Secret{
		Path:      "secret/api/key",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestTeamsNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, err := notifier.NewTeamsNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := vault.Secret{
		Path:      "secret/api/key",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Error("expected error for server error response")
	}
}
