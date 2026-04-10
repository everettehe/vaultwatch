package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewMattermostNotifier_Valid(t *testing.T) {
	n, err := notifier.NewMattermostNotifier("https://mattermost.example.com/hooks/abc", "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewMattermostNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewMattermostNotifier("", "", "")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestMattermostNotifier_Notify_ExpiringSoon(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewMattermostNotifier(server.URL, "#alerts", "vaultwatch")
	secret := &vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected webhook to be called")
	}
}

func TestMattermostNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewMattermostNotifier(server.URL, "", "")
	secret := &vault.Secret{
		Path:      "secret/api/key",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMattermostNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notifier.NewMattermostNotifier(server.URL, "", "")
	secret := &vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for server error response")
	}
}
