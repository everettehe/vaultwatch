package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewGotifyNotifier_Valid(t *testing.T) {
	n, err := NewGotifyNotifier("http://gotify.example.com", "apptoken123", 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewGotifyNotifier_MissingURL(t *testing.T) {
	_, err := NewGotifyNotifier("", "apptoken123", 5)
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestNewGotifyNotifier_MissingToken(t *testing.T) {
	_, err := NewGotifyNotifier("http://gotify.example.com", "", 5)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewGotifyNotifier_DefaultPriority(t *testing.T) {
	n, err := NewGotifyNotifier("http://gotify.example.com", "tok", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.priority != 5 {
		t.Errorf("expected default priority 5, got %d", n.priority)
	}
}

func TestGotifyNotifier_Notify_ExpiringSoon(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewGotifyNotifier(server.URL, "tok", 5)
	secret := &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(72 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGotifyNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewGotifyNotifier(server.URL, "tok", 8)
	secret := &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGotifyNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewGotifyNotifier(server.URL, "tok", 5)
	secret := &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Error("expected error for server error response")
	}
}
