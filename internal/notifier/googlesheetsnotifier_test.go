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

func newSheetsSecret(daysUntil int) *vault.Secret {
	exp := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return &vault.Secret{Path: "secret/db/password", ExpiresAt: exp}
}

func TestNewGoogleSheetsNotifier_Valid(t *testing.T) {
	n, err := notifier.NewGoogleSheetsNotifier("https://script.google.com/macros/s/abc/exec", "Alerts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewGoogleSheetsNotifier_MissingURL(t *testing.T) {
	_, err := notifier.NewGoogleSheetsNotifier("", "")
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestGoogleSheetsNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewGoogleSheetsNotifier(server.URL, "Vault")
	if err := n.Notify(newSheetsSecret(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["path"] != "secret/db/password" {
		t.Errorf("expected path secret/db/password, got %s", received["path"])
	}
	if received["sheet"] != "Vault" {
		t.Errorf("expected sheet Vault, got %s", received["sheet"])
	}
}

func TestGoogleSheetsNotifier_Notify_Expired(t *testing.T) {
	var received map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := notifier.NewGoogleSheetsNotifier(server.URL, "")
	if err := n.Notify(newSheetsSecret(-1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["severity"] != "EXPIRED" {
		t.Errorf("expected severity EXPIRED, got %s", received["severity"])
	}
}

func TestGoogleSheetsNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := notifier.NewGoogleSheetsNotifier(server.URL, "")
	if err := n.Notify(newSheetsSecret(5)); err == nil {
		t.Fatal("expected error on server error")
	}
}
