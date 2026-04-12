package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GrafanaNotifier sends annotations to a Grafana instance when secrets
// are expiring or have expired.
type GrafanaNotifier struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type grafanaAnnotation struct {
	Text string   `json:"text"`
	Tags []string `json:"tags"`
}

// NewGrafanaNotifier creates a new GrafanaNotifier. baseURL must be the
// root URL of the Grafana instance (e.g. https://grafana.example.com).
func NewGrafanaNotifier(baseURL, apiKey string) (*GrafanaNotifier, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("grafana: baseURL is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("grafana: apiKey is required")
	}
	return &GrafanaNotifier{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}, nil
}

// Notify posts an annotation to Grafana describing the secret's expiration state.
func (g *GrafanaNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	tags := []string{"vaultwatch", "secret-expiry"}
	if secret.IsExpired() {
		tags = append(tags, "expired")
	} else {
		tags = append(tags, "expiring-soon")
	}

	body, err := json.Marshal(grafanaAnnotation{
		Text: msg.Body,
		Tags: tags,
	})
	if err != nil {
		return fmt.Errorf("grafana: failed to marshal annotation: %w", err)
	}

	url := g.baseURL + "/api/annotations"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("grafana: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("grafana: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("grafana: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
