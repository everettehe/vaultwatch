package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewGoogleChatNotifier_Valid(t *testing.T) {
	n, err := notifier.NewGoogleChatNotifier("https://chat.googleapis.com/v1/spaces/test/messages?key=abc")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewGoogleChatNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewGoogleChatNotifier("")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestGoogleChatNotifier_Notify_ExpiringSoon(t *testing.T) {
	var capturedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.Body
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	_ = capturedBody

	n, err := notifier.NewGoogleChatNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/my-app/api-key",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGoogleChatNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := notifier.NewGoogleChatNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/my-app/db-pass",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGoogleChatNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, err := notifier.NewGoogleChatNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/my-app/token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for server error response")
	}
}
