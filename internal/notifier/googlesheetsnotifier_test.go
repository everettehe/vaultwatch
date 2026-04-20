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

func newSheetsSecret(days int) vault.Secret {
	return vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(time.Duration(days) * 24 * time.Hour),
	}
}

func TestNewGoogleSheetsNotifier_Valid(t *testing.T) {
	n, err := notifier.NewGoogleSheetsNotifier("https://script.google.com/macros/s/abc/exec")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewGoogleSheetsNotifier_MissingURL(t *testing.T) {
	_, err := notifier.NewGoogleSheetsNotifier("")
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestGoogleSheetsNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewGoogleSheetsNotifier(server.URL)
	secret := newSheetsSecret(5)
	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if received["path"] != "secret/myapp/db" {
		t.Errorf("unexpected path: %v", received["path"])
	}
}

func TestGoogleSheetsNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewGoogleSheetsNotifier(server.URL)
	secret := newSheetsSecret(-1)
	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGoogleSheetsNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notifier.NewGoogleSheetsNotifier(server.URL)
	secret := newSheetsSecret(3)
	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for server error response")
	}
}
