package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewGoogleDNSNotifier_Valid(t *testing.T) {
	n, err := notifier.NewGoogleDNSNotifier("https://example.com/hook", "my-project")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewGoogleDNSNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewGoogleDNSNotifier("", "my-project")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestNewGoogleDNSNotifier_MissingProject(t *testing.T) {
	_, err := notifier.NewGoogleDNSNotifier("https://example.com/hook", "")
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestGoogleDNSNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n, _ := notifier.NewGoogleDNSNotifier(srv.URL+"/hook", "proj")
	secret := &vault.Secret{Path: "secret/db", ExpiresAt: time.Now().Add(48 * time.Hour)}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received != "/hook" {
		t.Errorf("expected /hook, got %s", received)
	}
}

func TestGoogleDNSNotifier_Notify_Expired(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n, _ := notifier.NewGoogleDNSNotifier(srv.URL, "proj")
	secret := &vault.Secret{Path: "secret/db", ExpiresAt: time.Now().Add(-1 * time.Hour)}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGoogleDNSNotifier_Notify_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n, _ := notifier.NewGoogleDNSNotifier(srv.URL, "proj")
	secret := &vault.Secret{Path: "secret/db", ExpiresAt: time.Now().Add(24 * time.Hour)}
	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for 500 response")
	}
}
