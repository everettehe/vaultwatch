package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewGrafanaNotifier_Valid(t *testing.T) {
	n, err := notifier.NewGrafanaNotifier("https://grafana.example.com", "glsa_abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewGrafanaNotifier_MissingBaseURL(t *testing.T) {
	_, err := notifier.NewGrafanaNotifier("", "glsa_abc123")
	if err == nil {
		t.Fatal("expected error for missing baseURL")
	}
}

func TestNewGrafanaNotifier_MissingAPIKey(t *testing.T) {
	_, err := notifier.NewGrafanaNotifier("https://grafana.example.com", "")
	if err == nil {
		t.Fatal("expected error for missing apiKey")
	}
}

func TestGrafanaNotifier_Notify_ExpiringSoon(t *testing.T) {
	var capturedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewGrafanaNotifier(server.URL, "test-key")
	secret := &vault.Secret{
		Path:      "secret/myapp/api-key",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedAuth != "Bearer test-key" {
		t.Errorf("expected Authorization header 'Bearer test-key', got %q", capturedAuth)
	}
}

func TestGrafanaNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewGrafanaNotifier(server.URL, "test-key")
	secret := &vault.Secret{
		Path:      "secret/myapp/db-pass",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGrafanaNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notifier.NewGrafanaNotifier(server.URL, "test-key")
	secret := &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for server error response")
	}
}
