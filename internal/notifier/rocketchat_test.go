package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewRocketChatNotifier_Valid(t *testing.T) {
	n, err := notifier.NewRocketChatNotifier("https://chat.example.com/hooks/abc", "#alerts")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewRocketChatNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewRocketChatNotifier("", "#alerts")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestRocketChatNotifier_Notify_ExpiringSoon(t *testing.T) {
	var capturedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBody = make([]byte, r.ContentLength)
		r.Body.Read(capturedBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := notifier.NewRocketChatNotifier(server.URL, "#vault-alerts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/my-app/api-key",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
	if len(capturedBody) == 0 {
		t.Error("expected request body to be non-empty")
	}
}

func TestRocketChatNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewRocketChatNotifier(server.URL, "")
	secret := &vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
}

func TestRocketChatNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notifier.NewRocketChatNotifier(server.URL, "#alerts")
	secret := &vault.Secret{
		Path:      "secret/app/token",
		ExpiresAt: time.Now().Add(2 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server 500 response")
	}
}
