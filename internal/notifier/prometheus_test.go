package notifier

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewPrometheusNotifier_Valid(t *testing.T) {
	n, err := NewPrometheusNotifier("http://localhost:9091", "vaultwatch")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewPrometheusNotifier_DefaultJob(t *testing.T) {
	n, err := NewPrometheusNotifier("http://localhost:9091", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n.job != "vaultwatch" {
		t.Errorf("expected default job 'vaultwatch', got %q", n.job)
	}
}

func TestNewPrometheusNotifier_MissingURL(t *testing.T) {
	_, err := NewPrometheusNotifier("", "vaultwatch")
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestPrometheusNotifier_Notify_ExpiringSoon(t *testing.T) {
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		buf := make([]byte, 1024)
		n, _ := r.Body.Read(buf)
		receivedBody = string(buf[:n])
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier, err := NewPrometheusNotifier(server.URL, "vaultwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := notifier.Notify(secret); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if receivedBody == "" {
		t.Error("expected non-empty body sent to pushgateway")
	}
}

func TestPrometheusNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	notifier, err := NewPrometheusNotifier(server.URL, "vaultwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}

	if err := notifier.Notify(secret); err == nil {
		t.Error("expected error on server 500 response")
	}
}

func TestSanitizeLabelValue(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"secret/myapp/db", "secret_myapp_db"},
		{"/leading/slash", "leading_slash"},
		{"no_slashes", "no_slashes"},
	}
	for _, tt := range tests {
		got := sanitizeLabelValue(tt.input)
		if got != tt.expected {
			t.Errorf("sanitizeLabelValue(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
