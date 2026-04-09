package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewSlackNotifier_Valid(t *testing.T) {
	n, err := NewSlackNotifier("https://hooks.slack.com/test", "#alerts")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewSlackNotifier_MissingWebhook(t *testing.T) {
	_, err := NewSlackNotifier("", "#alerts")
	if err == nil {
		t.Fatal("expected error for empty webhook URL")
	}
}

func TestSlackNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received slackMessage

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewSlackNotifier(server.URL, "#vault-alerts")
	secret := vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Channel != "#vault-alerts" {
		t.Errorf("expected channel #vault-alerts, got %s", received.Channel)
	}
	if len(received.Blocks) == 0 {
		t.Fatal("expected at least one block in message")
	}
	if !strings.Contains(received.Blocks[0].Text.Text, "secret/db/password") {
		t.Errorf("expected secret path in message text, got: %s", received.Blocks[0].Text.Text)
	}
}

func TestSlackNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewSlackNotifier(server.URL, "")
	secret := vault.Secret{
		Path:      "secret/api/key",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSlackNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewSlackNotifier(server.URL, "#alerts")
	secret := vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server 500 response")
	}
}
