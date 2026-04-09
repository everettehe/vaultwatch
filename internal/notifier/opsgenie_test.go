package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/vault"
)

func TestNewOpsGenieNotifier_Valid(t *testing.T) {
	n, err := NewOpsGenieNotifier("test-api-key", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.apiURL != defaultOpsGenieURL {
		t.Errorf("expected default URL %q, got %q", defaultOpsGenieURL, n.apiURL)
	}
}

func TestNewOpsGenieNotifier_MissingKey(t *testing.T) {
	_, err := NewOpsGenieNotifier("", "")
	if err == nil {
		t.Fatal("expected error for missing api key")
	}
}

func TestOpsGenieNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n, _ := NewOpsGenieNotifier("key", server.URL)
	secret := vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["priority"] != "P3" {
		t.Errorf("expected priority P3, got %v", received["priority"])
	}
}

func TestOpsGenieNotifier_Notify_Expired(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n, _ := NewOpsGenieNotifier("key", server.URL)
	secret := vault.Secret{
		Path:      "secret/db/password",
		ExpiresAt: time.Now().Add(-2 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["priority"] != "P1" {
		t.Errorf("expected priority P1 for expired secret, got %v", received["priority"])
	}
}

func TestOpsGenieNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewOpsGenieNotifier("key", server.URL)
	secret := vault.Secret{
		Path:      "secret/api/token",
		ExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for server error response")
	}
}
