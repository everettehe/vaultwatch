package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewPagerDutyNotifier_Valid(t *testing.T) {
	n, err := NewPagerDutyNotifier("test-key-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewPagerDutyNotifier_MissingKey(t *testing.T) {
	_, err := NewPagerDutyNotifier("")
	if err == nil {
		t.Fatal("expected error for missing integration key")
	}
}

func TestPagerDutyNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received pdPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n, _ := NewPagerDutyNotifier("key-abc")
	n.eventURL = server.URL

	secret := vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.EventAction != "trigger" {
		t.Errorf("expected event_action=trigger, got %q", received.EventAction)
	}
	if received.Payload.Severity != "warning" {
		t.Errorf("expected severity=warning, got %q", received.Payload.Severity)
	}
}

func TestPagerDutyNotifier_Notify_Expired(t *testing.T) {
	var received pdPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n, _ := NewPagerDutyNotifier("key-abc")
	n.eventURL = server.URL

	secret := vault.Secret{
		Path:      "secret/api/key",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Payload.Severity != "critical" {
		t.Errorf("expected severity=critical, got %q", received.Payload.Severity)
	}
}

func TestPagerDutyNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewPagerDutyNotifier("key-abc")
	n.eventURL = server.URL

	secret := vault.Secret{
		Path:      "secret/token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server 500 response")
	}
}
