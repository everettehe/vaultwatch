package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewLarkNotifier_Valid(t *testing.T) {
	n, err := notifier.NewLarkNotifier("https://open.larksuite.com/open-apis/bot/v2/hook/abc")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewLarkNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewLarkNotifier("")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestLarkNotifier_Notify_ExpiringSoon(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := notifier.NewLarkNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/my-app/api-key",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLarkNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := notifier.NewLarkNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/my-app/api-key",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLarkNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, err := notifier.NewLarkNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/my-app/api-key",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Error("expected error for server error response")
	}
}
