package notifier_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewBearyChatNotifier_Valid(t *testing.T) {
	_, err := notifier.NewBearyChatNotifier("https://hook.bearychat.com/abc")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewBearyChatNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewBearyChatNotifier("")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestBearyChatNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := notifier.NewBearyChatNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	if received["text"] == nil {
		t.Error("expected text field in payload")
	}
}

func TestBearyChatNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewBearyChatNotifier(server.URL)
	secret := &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
}

func TestBearyChatNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notifier.NewBearyChatNotifier(server.URL)
	secret := &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server error response")
	}
}
