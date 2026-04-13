package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewCampfireNotifier_Valid(t *testing.T) {
	n, err := notifier.NewCampfireNotifier("https://3.basecamp.com/12345/integrations/abc/buckets/xyz/chats/1/lines")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewCampfireNotifier_MissingWebhook(t *testing.T) {
	_, err := notifier.NewCampfireNotifier("")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

// newTestCampfireServer creates a test HTTP server and a CampfireNotifier pointed at it.
// The provided handler is used to respond to incoming requests.
func newTestCampfireServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *notifier.CampfireNotifier) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	n, err := notifier.NewCampfireNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error creating notifier: %v", err)
	}
	return server, n
}

func TestCampfireNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received bool
	_, n := newTestCampfireServer(t, func(w http.ResponseWriter, r *http.Request) {
		received = true
		w.WriteHeader(http.StatusOK)
	})

	secret := &vault.Secret{
		Path:      "secret/data/myapp/api-key",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !received {
		t.Fatal("expected server to receive a request")
	}
}

func TestCampfireNotifier_Notify_Expired(t *testing.T) {
	_, n := newTestCampfireServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	secret := &vault.Secret{
		Path:      "secret/data/myapp/db-pass",
		ExpiresAt: time.Now().Add(-24 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCampfireNotifier_Notify_ServerError(t *testing.T) {
	_, n := newTestCampfireServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	secret := &vault.Secret{
		Path:      "secret/data/myapp/token",
		ExpiresAt: time.Now().Add(12 * time.Hour),
	}

	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server error response")
	}
}
