package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewZulipNotifier_Valid(t *testing.T) {
	n, err := NewZulipNotifier("https://org.zulipchat.com", "bot@org.com", "apikey", "alerts", "vault")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewZulipNotifier_DefaultTopic(t *testing.T) {
	n, err := NewZulipNotifier("https://org.zulipchat.com", "bot@org.com", "apikey", "alerts", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.topic != "vaultwatch alerts" {
		t.Errorf("expected default topic, got %q", n.topic)
	}
}

func TestNewZulipNotifier_MissingBaseURL(t *testing.T) {
	_, err := NewZulipNotifier("", "bot@org.com", "apikey", "alerts", "vault")
	if err == nil {
		t.Fatal("expected error for missing base URL")
	}
}

func TestNewZulipNotifier_MissingEmail(t *testing.T) {
	_, err := NewZulipNotifier("https://org.zulipchat.com", "", "apikey", "alerts", "vault")
	if err == nil {
		t.Fatal("expected error for missing email")
	}
}

func TestNewZulipNotifier_MissingAPIKey(t *testing.T) {
	_, err := NewZulipNotifier("https://org.zulipchat.com", "bot@org.com", "", "alerts", "vault")
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
}

func TestNewZulipNotifier_MissingStream(t *testing.T) {
	_, err := NewZulipNotifier("https://org.zulipchat.com", "bot@org.com", "apikey", "", "vault")
	if err == nil {
		t.Fatal("expected error for missing stream")
	}
}

func TestZulipNotifier_Notify_ExpiringSoon(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}
		if r.FormValue("type") != "stream" {
			t.Errorf("expected type=stream, got %q", r.FormValue("type"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"success"}`))
	}))
	defer server.Close()

	n, _ := NewZulipNotifier(server.URL, "bot@org.com", "apikey", "alerts", "vault")
	secret := &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestZulipNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewZulipNotifier(server.URL, "bot@org.com", "apikey", "alerts", "vault")
	secret := &vault.Secret{
		Path:      "secret/myapp/token",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server error")
	}
}
