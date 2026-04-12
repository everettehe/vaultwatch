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

func TestNewSignaldNotifier_Valid(t *testing.T) {
	n, err := notifier.NewSignaldNotifier("http://localhost:8080", "+10000000000", []string{"+19999999999"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewSignaldNotifier_MissingBaseURL(t *testing.T) {
	_, err := notifier.NewSignaldNotifier("", "+10000000000", []string{"+19999999999"})
	if err == nil {
		t.Fatal("expected error for missing base URL")
	}
}

func TestNewSignaldNotifier_MissingSender(t *testing.T) {
	_, err := notifier.NewSignaldNotifier("http://localhost:8080", "", []string{"+19999999999"})
	if err == nil {
		t.Fatal("expected error for missing sender")
	}
}

func TestNewSignaldNotifier_MissingRecipients(t *testing.T) {
	_, err := notifier.NewSignaldNotifier("http://localhost:8080", "+10000000000", nil)
	if err == nil {
		t.Fatal("expected error for missing recipients")
	}
}

func TestSignaldNotifier_Notify_ExpiringSoon(t *testing.T) {
	var captured map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/send" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewDecoder(r.Body).Decode(&captured)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, err := notifier.NewSignaldNotifier(ts.URL, "+10000000000", []string{"+19999999999"})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/my-app/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
	if captured["username"] != "+10000000000" {
		t.Errorf("unexpected sender: %v", captured["username"])
	}
}

func TestSignaldNotifier_Notify_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, err := notifier.NewSignaldNotifier(ts.URL, "+10000000000", []string{"+19999999999"})
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/my-app/db",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on server error response")
	}
}
