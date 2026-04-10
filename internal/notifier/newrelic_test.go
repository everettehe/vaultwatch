package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewNewRelicNotifier_Valid(t *testing.T) {
	n, err := NewNewRelicNotifier("123456", "test-api-key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewNewRelicNotifier_MissingAccountID(t *testing.T) {
	_, err := NewNewRelicNotifier("", "test-api-key")
	if err == nil {
		t.Fatal("expected error for missing account ID")
	}
}

func TestNewNewRelicNotifier_MissingAPIKey(t *testing.T) {
	_, err := NewNewRelicNotifier("123456", "")
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestNewRelicNotifier_Notify_ExpiringSoon(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Insert-Key") == "" {
			t.Error("expected X-Insert-Key header")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewNewRelicNotifier("123456", "test-api-key")
	n.baseURL = server.URL

	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewRelicNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewNewRelicNotifier("123456", "test-api-key")
	n.baseURL = server.URL

	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewRelicNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewNewRelicNotifier("123456", "test-api-key")
	n.baseURL = server.URL

	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server error response")
	}
}
