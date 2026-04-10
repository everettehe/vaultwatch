package notifier_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewJiraNotifier_Valid(t *testing.T) {
	n, err := notifier.NewJiraNotifier("https://example.atlassian.net", "token", "OPS", "Bug")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestNewJiraNotifier_MissingURL(t *testing.T) {
	_, err := notifier.NewJiraNotifier("", "token", "OPS", "Task")
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestNewJiraNotifier_MissingToken(t *testing.T) {
	_, err := notifier.NewJiraNotifier("https://example.atlassian.net", "", "OPS", "Task")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewJiraNotifier_MissingProject(t *testing.T) {
	_, err := notifier.NewJiraNotifier("https://example.atlassian.net", "token", "", "Task")
	if err == nil {
		t.Fatal("expected error for missing project")
	}
}

func TestNewJiraNotifier_DefaultIssueType(t *testing.T) {
	n, err := notifier.NewJiraNotifier("https://example.atlassian.net", "token", "OPS", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected notifier, got nil")
	}
}

func TestJiraNotifier_Notify_ExpiringSoon(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	n, err := notifier.NewJiraNotifier(server.URL, "token", "OPS", "Task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}
	if err := n.Notify(secret); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestJiraNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, err := notifier.NewJiraNotifier(server.URL, "token", "OPS", "Task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secret := &vault.Secret{
		Path:      "secret/myapp/db",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}
	if err := n.Notify(secret); err == nil {
		t.Error("expected error on server error response")
	}
}
