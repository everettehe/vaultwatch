package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewSignalRNotifier_Valid(t *testing.T) {
	n, err := notifier.NewSignalRNotifier("https://example.service.signalr.net", "mykey", "alerts")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewSignalRNotifier_MissingURL(t *testing.T) {
	_, err := notifier.NewSignalRNotifier("", "mykey", "alerts")
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestNewSignalRNotifier_MissingAccessKey(t *testing.T) {
	_, err := notifier.NewSignalRNotifier("https://example.service.signalr.net", "", "alerts")
	if err == nil {
		t.Fatal("expected error for missing access key")
	}
}

func TestNewSignalRNotifier_DefaultHub(t *testing.T) {
	n, err := notifier.NewSignalRNotifier("https://example.service.signalr.net", "mykey", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestSignalRNotifier_Notify_ExpiringSoon(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewSignalRNotifier(server.URL, "mykey", "alerts")
	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSignalRNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewSignalRNotifier(server.URL, "mykey", "alerts")
	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSignalRNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notifier.NewSignalRNotifier(server.URL, "mykey", "alerts")
	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for server error response")
	}
}
