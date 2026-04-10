package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewVictorOpsNotifier_Valid(t *testing.T) {
	n, err := NewVictorOpsNotifier("https://alert.victorops.com/integrations/generic", "my-routing-key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewVictorOpsNotifier_MissingWebhook(t *testing.T) {
	_, err := NewVictorOpsNotifier("", "my-routing-key")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestNewVictorOpsNotifier_MissingRoutingKey(t *testing.T) {
	_, err := NewVictorOpsNotifier("https://alert.victorops.com/integrations/generic", "")
	if err == nil {
		t.Fatal("expected error for missing routing key")
	}
}

func TestVictorOpsNotifier_Notify_ExpiringSoon(t *testing.T) {
	var capturedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBody = make([]byte, r.ContentLength)
		r.Body.Read(capturedBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewVictorOpsNotifier(server.URL, "test-key")
	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(capturedBody) == 0 {
		t.Fatal("expected request body to be sent")
	}
}

func TestVictorOpsNotifier_Notify_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, _ := NewVictorOpsNotifier(server.URL, "test-key")
	secret := &vault.Secret{
		Path:      "secret/myapp/api",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestVictorOpsNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewVictorOpsNotifier(server.URL, "test-key")
	secret := &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error for server error response")
	}
}
