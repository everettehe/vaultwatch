package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// IncidentNotifier sends alerts to Incident.io via their Alerts API.
type IncidentNotifier struct {
	apiKey string
	alertsURL string
	client *http.Client
}

type incidentPayload struct {
	Title       string            `json:"title"`
	Message     string            `json:"message"`
	Status      string            `json:"status"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// NewIncidentNotifier creates a new Incident.io notifier.
// apiKey must be a valid Incident.io API key.
func NewIncidentNotifier(apiKey string) (*IncidentNotifier, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("incident.io: api key is required")
	}
	return &IncidentNotifier{
		apiKey:    apiKey,
		alertsURL: "https://api.incident.io/v2/alert_events/http",
		client:    &http.Client{},
	}, nil
}

// Notify sends an alert to Incident.io for the given secret.
func (n *IncidentNotifier) Notify(s *vault.Secret) error {
	msg := FormatMessage(s)

	status := "firing"
	if s.IsExpired() {
		status = "firing"
	}

	payload := incidentPayload{
		Title:   msg.Subject,
		Message: msg.Body,
		Status:  status,
		Metadata: map[string]string{
			"path":        s.Path,
			"days_until":  fmt.Sprintf("%d", s.DaysUntilExpiration()),
			"source":      "vaultwatch",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("incident.io: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.alertsURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("incident.io: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.apiKey)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("incident.io: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("incident.io: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
