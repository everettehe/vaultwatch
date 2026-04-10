package notifier_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewDiscordNotifier_Valid(t *testing.T) {
	n, err := notifier.NewDiscordNotifier("https://discord.com/api/webhooks/test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewDiscordNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewDiscordNotifier("")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
	if !strings.Contains(err.Error(), "webhook URL is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDiscordNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	n, _ := notifier.NewDiscordNotifier(server.URL)
	secret := &vault.Secret{
		Path:      "secret/myapp/api-key",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, ok := received["content"].(string)
	if !ok || content == "" {
		t.Error("expected non-empty content in Discord payload")
	}
}

func TestDiscordNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	n, _ := notifier.NewDiscordNotifier(server.URL)
	secret := &vault.Secret{
		Path:      "secret/myapp/db-pass",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiscordNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notifier.NewDiscordNotifier(server.URL)
	secret := &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for server error response")
	}
}
