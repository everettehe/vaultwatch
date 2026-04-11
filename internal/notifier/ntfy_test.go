package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewNtfyNotifier_Valid(t *testing.T) {
	n, err := NewNtfyNotifier("https://ntfy.sh", "vaultwatch-alerts")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.topic != "vaultwatch-alerts" {
		t.Errorf("expected topic 'vaultwatch-alerts', got %q", n.topic)
	}
}

func TestNewNtfyNotifier_DefaultServer(t *testing.T) {
	n, err := NewNtfyNotifier("", "vaultwatch-alerts")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.serverURL != "https://ntfy.sh" {
		t.Errorf("expected default server URL, got %q", n.serverURL)
	}
}

func TestNewNtfyNotifier_MissingTopic(t *testing.T) {
	_, err := NewNtfyNotifier("https://ntfy.sh", "")
	if err == nil {
		t.Fatal("expected error for missing topic")
	}
}

func TestNtfyNotifier_Notify_ExpiringSoon(t *testing.T) {
	var receivedTitle, receivedPriority string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedTitle = r.Header.Get("Title")
		receivedPriority = r.Header.Get("Priority")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewNtfyNotifier(server.URL, "alerts")
	secret := &vault.Secret{
		Path:      "secret/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedTitle == "" {
		t.Error("expected Title header to be set")
	}
	if receivedPriority != "high" {
		t.Errorf("expected priority 'high', got %q", receivedPriority)
	}
}

func TestNtfyNotifier_Notify_Expired(t *testing.T) {
	var receivedPriority string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPriority = r.Header.Get("Priority")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewNtfyNotifier(server.URL, "alerts")
	secret := &vault.Secret{
		Path:      "secret/db",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if receivedPriority != "urgent" {
		t.Errorf("expected priority 'urgent', got %q", receivedPriority)
	}
}

func TestNtfyNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewNtfyNotifier(server.URL, "alerts")
	secret := &vault.Secret{
		Path:      "secret/db",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server error response")
	}
}
