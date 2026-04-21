package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func newHTTPPostSecret() *vault.Secret {
	return &vault.Secret{
		Path:      "secret/data/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
}

func TestNewHTTPPostNotifier_Valid(t *testing.T) {
	n, err := NewHTTPPostNotifier("http://example.com/hook", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewHTTPPostNotifier_MissingURL(t *testing.T) {
	_, err := NewHTTPPostNotifier("", nil)
	if err == nil {
		t.Fatal("expected error for missing url")
	}
}

func TestHTTPPostNotifier_Notify_ExpiringSoon(t *testing.T) {
	var got map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewHTTPPostNotifier(ts.URL, map[string]string{"X-Token": "abc"})
	if err := n.Notify(newHTTPPostSecret()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["path"] != "secret/data/db" {
		t.Errorf("unexpected path: %s", got["path"])
	}
	if got["status"] == "" {
		t.Error("expected non-empty status")
	}
}

func TestHTTPPostNotifier_Notify_Expired(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n, _ := NewHTTPPostNotifier(ts.URL, nil)
	s := &vault.Secret{Path: "secret/data/old", ExpiresAt: time.Now().Add(-1 * time.Hour)}
	if err := n.Notify(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHTTPPostNotifier_Notify_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, _ := NewHTTPPostNotifier(ts.URL, nil)
	if err := n.Notify(newHTTPPostSecret()); err == nil {
		t.Fatal("expected error on 500")
	}
}
