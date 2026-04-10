package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// ZendutyNotifier sends alerts to Zenduty via its alert API.
type ZendutyNotifier struct {
	apiKey      string
	serviceID   string
	integrKey   string
	httpClient  *http.Client
}

type zendutyPayload struct {
	Message   string `json:"message"`
	AlertType string `json:"alert_type"`
	Summary   string `json:"summary"`
	EntityID  string `json:"entity_id"`
}

// NewZendutyNotifier creates a new ZendutyNotifier.
// apiKey, serviceID, and integrationKey are required.
func NewZendutyNotifier(apiKey, serviceID, integrationKey string) (*ZendutyNotifier, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("zenduty: api key is required")
	}
	if serviceID == "" {
		return nil, fmt.Errorf("zenduty: service ID is required")
	}
	if integrationKey == "" {
		return nil, fmt.Errorf("zenduty: integration key is required")
	}
	return &ZendutyNotifier{
		apiKey:     apiKey,
		serviceID:  serviceID,
		integrKey:  integrationKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends an alert to Zenduty for the given secret.
func (z *ZendutyNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	alertType := "warning"
	if secret.IsExpired() {
		alertType = "critical"
	}

	payload := zendutyPayload{
		Message:   msg.Subject,
		AlertType: alertType,
		Summary:   msg.Body,
		EntityID:  secret.Path,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("zenduty: failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("https://www.zenduty.com/api/v1/integrations/%s/alerts/", z.integrKey)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("zenduty: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+z.apiKey)

	resp, err := z.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("zenduty: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("zenduty: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
