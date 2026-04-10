package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewSplunkNotifier_Valid(t *testing.T) {
	n, err := NewSplunkNotifier("http://splunk:8088/services/collector", "my-token", "vaultwatch", "_json", "main")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewSplunkNotifier_MissingURL(t *testing.T) {
	_, err := NewSplunkNotifier("", "my-token", "", "", "")
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestNewSplunkNotifier_MissingToken(t *testing.T) {
	_, err := NewSplunkNotifier("http://splunk:8088/services/collector", "", "", "", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestSplunkNotifier_Notify_ExpiringSoon(t *testing.T) {
	var captured splunkEvent
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Splunk test-token" {
			t.Errorf("unexpected Authorization header: %s", r.Header.Get("Authorization"))
		}
		if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewSplunkNotifier(ts.URL, "test-token", "vaultwatch", "_json", "")
	secret := &vault.Secret{
		Path:      "secret/db",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Event["secret_path"] != "secret/db" {
		t.Errorf("expected secret_path=secret/db, got %v", captured.Event["secret_path"])
	}
}

func TestSplunkNotifier_Notify_Expired(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewSplunkNotifier(ts.URL, "test-token", "", "", "")
	secret := &vault.Secret{
		Path:      "secret/expired",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSplunkNotifier_Notify_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, _ := NewSplunkNotifier(ts.URL, "test-token", "", "", "")
	secret := &vault.Secret{
		Path:      "secret/db",
		ExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server error response")
	}
}
