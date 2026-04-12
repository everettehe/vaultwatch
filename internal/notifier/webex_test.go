package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewWebexNotifier_Valid(t *testing.T) {
	n, err := NewWebexNotifier("tok", "room123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewWebexNotifier_MissingToken(t *testing.T) {
	_, err := NewWebexNotifier("", "room123")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewWebexNotifier_MissingRoomID(t *testing.T) {
	_, err := NewWebexNotifier("tok", "")
	if err == nil {
		t.Fatal("expected error for missing room_id")
	}
}

func TestWebexNotifier_Notify_ExpiringSoon(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewWebexNotifier("mytoken", "roomABC")
	// Override the API URL via a custom client that redirects to the test server.
	n.client = server.Client()
	// Patch the request target by wrapping transport.
	n.client.Transport = rewriteTransport(server.URL)

	secret := &vault.Secret{
		Path:      "secret/db",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Bearer mytoken" {
		t.Errorf("expected Bearer mytoken, got %q", gotAuth)
	}
}

func TestWebexNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewWebexNotifier("tok", "room")
	n.client = server.Client()
	n.client.Transport = rewriteTransport(server.URL)

	secret := &vault.Secret{
		Path:      "secret/expired",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWebexNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewWebexNotifier("tok", "room")
	n.client = server.Client()
	n.client.Transport = rewriteTransport(server.URL)

	secret := &vault.Secret{
		Path:      "secret/db",
		ExpiresAt: time.Now().Add(2 * 24 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for server error response")
	}
}
