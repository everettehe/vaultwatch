package notifier_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewCircleCINotifier_Valid(t *testing.T) {
	n, err := notifier.NewCircleCINotifier("mytoken", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewCircleCINotifier_MissingToken(t *testing.T) {
	_, err := notifier.NewCircleCINotifier("", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestCircleCINotifier_Notify_ExpiringSoon(t *testing.T) {
	var received map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Circle-Token") == "" {
			t.Error("expected Circle-Token header")
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, err := notifier.NewCircleCINotifier("tok", ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
	if received["subject"] == "" {
		t.Error("expected non-empty subject in payload")
	}
}

func TestCircleCINotifier_Notify_Expired(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, err := notifier.NewCircleCINotifier("tok", ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
}

func TestCircleCINotifier_Notify_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n, err := notifier.NewCircleCINotifier("bad-tok", ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
