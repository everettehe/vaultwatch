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

func newGoogleTaskSecret(daysUntil int) vault.Secret {
	expiry := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return vault.Secret{Path: "secret/myapp/api-key", ExpiresAt: expiry}
}

func TestNewGoogleTaskNotifier_Valid(t *testing.T) {
	n, err := notifier.NewGoogleTaskNotifier("https://script.google.com/macros/s/abc/exec", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewGoogleTaskNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewGoogleTaskNotifier("", "")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestNewGoogleTaskNotifier_DefaultTasklist(t *testing.T) {
	n, err := notifier.NewGoogleTaskNotifier("https://example.com/webhook", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestGoogleTaskNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := notifier.NewGoogleTaskNotifier(ts.URL, "mylist")
	if err := n.Notify(newGoogleTaskSecret(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["title"] == "" {
		t.Error("expected non-empty title")
	}
}

func TestGoogleTaskNotifier_Notify_Expired(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := notifier.NewGoogleTaskNotifier(ts.URL, "")
	if err := n.Notify(newGoogleTaskSecret(-1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGoogleTaskNotifier_Notify_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, _ := notifier.NewGoogleTaskNotifier(ts.URL, "")
	if err := n.Notify(newGoogleTaskSecret(3)); err == nil {
		t.Fatal("expected error on server error")
	}
}
