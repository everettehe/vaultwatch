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

func TestNewGoogleChatCardNotifier_Valid(t *testing.T) {
	n, err := notifier.NewGoogleChatCardNotifier("https://chat.googleapis.com/webhook" nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
("expected notifier, got nil")
	}
}

func TestNewGoogleChatCardNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewGoogleChatCardNotifier("")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func newGoogleChatCardSecret(daysUntil int) *vault.Secret {
	expiry := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: expiry,
	}
}

func TestGoogleChatCardNotifier_Notify_ExpiringSoon(t *testing.T) {
	var captured []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewGoogleChatCardNotifier(server.URL)
	secret := newGoogleChatCardSecret(5)

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(captured, &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if _, ok := payload["cards"]; !ok {
		t.Error("expected 'cards' key in payload")
	}
}

func TestGoogleChatCardNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewGoogleChatCardNotifier(server.URL)
	secret := newGoogleChatCardSecret(-1)

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGoogleChatCardNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notifier.NewGoogleChatCardNotifier(server.URL)
	secret := newGoogleChatCardSecret(3)

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server error response")
	}
}
