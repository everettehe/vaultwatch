package notifier

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func TestNewWebhookNotifier_Valid(t *testing.T) {
	n, err := NewWebhookNotifier(WebhookConfig{URL: "http://example.com/hook"})
	require.NoError(t, err)
	assert.NotNil(t, n)
}

func TestNewWebhookNotifier_MissingURL(t *testing.T) {
	_, err := NewWebhookNotifier(WebhookConfig{})
	assert.ErrorContains(t, err, "url")
}

func TestNewWebhookNotifier_DefaultMethod(t *testing.T) {
	n, err := NewWebhookNotifier(WebhookConfig{URL: "http://example.com/hook"})
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, n.cfg.Method)
}

func TestWebhookNotifier_Notify_ExpiringSoon(t *testing.T) {
	var received webhookPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&received))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := NewWebhookNotifier(WebhookConfig{URL: server.URL})
	require.NoError(t, err)

	secret := &vault.Secret{
		Path:      "secret/db/pass",
		ExpiresAt: time.Now().Add(3 * 24 * time.Hour),
	}
	require.NoError(t, n.Notify(secret))

	assert.Equal(t, "secret/db/pass", received.Path)
	assert.Equal(t, "expiring", received.Status)
}

func TestWebhookNotifier_Notify_Expired(t *testing.T) {
	var received webhookPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&received))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	n, err := NewWebhookNotifier(WebhookConfig{URL: server.URL})
	require.NoError(t, err)

	secret := &vault.Secret{
		Path:      "secret/api/token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	require.NoError(t, n.Notify(secret))
	assert.Equal(t, "expired", received.Status)
}

func TestWebhookNotifier_Notify_CustomHeaders(t *testing.T) {
	var authHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := NewWebhookNotifier(WebhookConfig{
		URL:     server.URL,
		Headers: map[string]string{"Authorization": "Bearer secret-token"},
	})
	require.NoError(t, err)

	secret := &vault.Secret{Path: "p", ExpiresAt: time.Now().Add(time.Hour)}
	require.NoError(t, n.Notify(secret))
	assert.Equal(t, "Bearer secret-token", authHeader)
}

func TestWebhookNotifier_Notify_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, err := NewWebhookNotifier(WebhookConfig{URL: server.URL})
	require.NoError(t, err)

	secret := &vault.Secret{Path: "p", ExpiresAt: time.Now().Add(time.Hour)}
	err = n.Notify(secret)
	assert.ErrorContains(t, err, "unexpected status 500")
}
