package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewDatadogNotifier_Valid(t *testing.T) {
	n, err := NewDatadogNotifier("test-api-key", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.apiURL != defaultDatadogAPIURL {
		t.Errorf("expected default API URL, got %s", n.apiURL)
	}
}

func TestNewDatadogNotifier_MissingKey(t *testing.T) {
	_, err := NewDatadogNotifier("", "")
	if err == nil {
		t.Fatal("expected error for missing api key")
	}
}

func TestDatadogNotifier_Notify_ExpiringSoon(t *testing.T) {
	var capturedHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Get("DD-API-KEY")
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n, err := NewDatadogNotifier("my-key", server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := &vault.Secret{
		Path:      "secret/my-app/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(s); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if capturedHeader != "my-key" {
		t.Errorf("expected DD-API-KEY header 'my-key', got '%s'", capturedHeader)
	}
}

func TestDatadogNotifier_Notify_Expired(t *testing.T) {
	var alertType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := decodeJSON(r.Body, &body); err == nil {
			if v, ok := body["alert_type"].(string); ok {
				alertType = v
			}
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n, _ := NewDatadogNotifier("key", server.URL)
	s := &vault.Secret{
		Path:      "secret/expired",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if alertType != "error" {
		t.Errorf("expected alert_type 'error', got '%s'", alertType)
	}
}

func TestDatadogNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewDatadogNotifier("key", server.URL)
	s := &vault.Secret{
		Path:      "secret/test",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := n.Notify(s); err == nil {
		t.Error("expected error on server error response")
	}
}
