package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// WebhookConfig holds configuration for a generic HTTP webhook notifier.
type WebhookConfig struct {
	URL     string
	Method  string
	Headers map[string]string
	Timeout time.Duration
}

// WebhookNotifier sends secret expiration alerts to a generic HTTP endpoint.
type WebhookNotifier struct {
	cfg    WebhookConfig
	client *http.Client
}

type webhookPayload struct {
	Path      string `json:"path"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at"`
	DaysLeft  int64  `json:"days_left"`
}

// NewWebhookNotifier creates a new WebhookNotifier from the given config.
func NewWebhookNotifier(cfg WebhookConfig) (*WebhookNotifier, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("webhook url is required")
	}
	if cfg.Method == "" {
		cfg.Method = http.MethodPost
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	return &WebhookNotifier{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}, nil
}

// Notify sends an HTTP request with secret expiration details to the configured endpoint.
func (w *WebhookNotifier) Notify(secret *vault.Secret) error {
	status := "expiring"
	if secret.IsExpired() {
		status = "expired"
	}

	payload := webhookPayload{
		Path:      secret.Path,
		Status:    status,
		ExpiresAt: secret.ExpiresAt.Format(time.RFC3339),
		DaysLeft:  secret.DaysUntilExpiration(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	req, err := http.NewRequest(w.cfg.Method, w.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range w.cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
