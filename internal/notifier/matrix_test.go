package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewMatrixNotifier_Valid(t *testing.T) {
	n, err := NewMatrixNotifier("https://matrix.example.com", "token123", "!room:example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewMatrixNotifier_MissingHomeserver(t *testing.T) {
	_, err := NewMatrixNotifier("", "token123", "!room:example.com")
	if err == nil {
		t.Fatal("expected error for missing homeserver")
	}
}

func TestNewMatrixNotifier_MissingToken(t *testing.T) {
	_, err := NewMatrixNotifier("https://matrix.example.com", "", "!room:example.com")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewMatrixNotifier_MissingRoomID(t *testing.T) {
	_, err := NewMatrixNotifier("https://matrix.example.com", "token123", "")
	if err == nil {
		t.Fatal("expected error for missing room ID")
	}
}

func TestMatrixNotifier_Notify_ExpiringSoon(t *testing.T) {
	var capturedBody matrixMessage
	var capturedAuth string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		_ = json.NewDecoder(r.Body).Decode(&capturedBody)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"event_id":"$abc123"}`))
	}))
	defer ts.Close()

	n, _ := NewMatrixNotifier(ts.URL, "mytoken", "!room:example.com")

	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if capturedAuth != "Bearer mytoken" {
		t.Errorf("expected Bearer mytoken, got %q", capturedAuth)
	}
	if capturedBody.MsgType != "m.text" {
		t.Errorf("expected msgtype m.text, got %q", capturedBody.MsgType)
	}
	if capturedBody.Body == "" {
		t.Error("expected non-empty body")
	}
}

func TestMatrixNotifier_Notify_Expired(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"event_id":"$xyz"}`))
	}))
	defer ts.Close()

	n, _ := NewMatrixNotifier(ts.URL, "tok", "!room:example.com")
	secret := &vault.Secret{
		Path:      "secret/expired",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestMatrixNotifier_Notify_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n, _ := NewMatrixNotifier(ts.URL, "tok", "!room:example.com")
	secret := &vault.Secret{
		Path:      "secret/test",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}
