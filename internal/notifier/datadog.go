package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

const defaultDatadogAPIURL = "https://api.datadoghq.com/api/v1/events"

// DatadogNotifier sends Vault secret expiration events to Datadog.
type DatadogNotifier struct {
	apiKey string
	apiURL string
	client *http.Client
}

type datadogEvent struct {
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	AlertType string   `json:"alert_type"`
	Tags      []string `json:"tags"`
}

// NewDatadogNotifier creates a DatadogNotifier. apiKey is required.
func NewDatadogNotifier(apiKey, apiURL string) (*DatadogNotifier, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("datadog: api key is required")
	}
	if apiURL == "" {
		apiURL = defaultDatadogAPIURL
	}
	return &DatadogNotifier{
		apiKey: apiKey,
		apiURL: apiURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a Datadog event for the given secret.
func (d *DatadogNotifier) Notify(s *vault.Secret) error {
	msg := FormatMessage(s)
	alertType := "warning"
	if s.IsExpired() {
		alertType = "error"
	}

	event := datadogEvent{
		Title:     msg.Subject,
		Text:      msg.Body,
		AlertType: alertType,
		Tags:      []string{"source:vaultwatch", fmt.Sprintf("secret_path:%s", s.Path)},
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("datadog: failed to marshal event: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, d.apiURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("datadog: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("datadog: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
