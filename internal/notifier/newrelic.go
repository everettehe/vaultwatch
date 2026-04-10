package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

const newRelicEventsURL = "https://insights-collector.newrelic.com/v1/accounts/%s/events"

// NewRelicNotifier sends custom events to New Relic Insights.
type NewRelicNotifier struct {
	accountID string
	apiKey    string
	client    *http.Client
	baseURL   string
}

type newRelicEvent struct {
	EventType  string  `json:"eventType"`
	SecretPath string  `json:"secretPath"`
	DaysLeft   float64 `json:"daysUntilExpiration"`
	Severity   string  `json:"severity"`
	Message    string  `json:"message"`
	Timestamp  int64   `json:"timestamp"`
}

// NewNewRelicNotifier creates a New Relic notifier.
func NewNewRelicNotifier(accountID, apiKey string) (*NewRelicNotifier, error) {
	if accountID == "" {
		return nil, fmt.Errorf("newrelic: account ID is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("newrelic: API key is required")
	}
	return &NewRelicNotifier{
		accountID: accountID,
		apiKey:    apiKey,
		client:    &http.Client{Timeout: 10 * time.Second},
		baseURL:   fmt.Sprintf(newRelicEventsURL, accountID),
	}, nil
}

// Notify sends a custom event to New Relic Insights.
func (n *NewRelicNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	event := newRelicEvent{
		EventType:  "VaultSecretExpiration",
		SecretPath: secret.Path,
		DaysLeft:   secret.DaysUntilExpiration(),
		Severity:   msg.Severity,
		Message:    msg.Body,
		Timestamp:  time.Now().Unix(),
	}

	payload, err := json.Marshal([]newRelicEvent{event})
	if err != nil {
		return fmt.Errorf("newrelic: failed to marshal event: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.baseURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("newrelic: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Insert-Key", n.apiKey)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("newrelic: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("newrelic: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
