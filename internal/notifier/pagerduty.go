package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

const pagerDutyEventURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyNotifier sends alerts to PagerDuty via the Events API v2.
type PagerDutyNotifier struct {
	integrationKey string
	client         *http.Client
	eventURL       string
}

type pdPayload struct {
	RoutingKey  string    `json:"routing_key"`
	EventAction string    `json:"event_action"`
	Payload     pdDetails `json:"payload"`
}

type pdDetails struct {
	Summary   string `json:"summary"`
	Source    string `json:"source"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
}

// NewPagerDutyNotifier creates a PagerDutyNotifier with the given integration key.
func NewPagerDutyNotifier(integrationKey string) (*PagerDutyNotifier, error) {
	if integrationKey == "" {
		return nil, fmt.Errorf("pagerduty: integration key is required")
	}
	return &PagerDutyNotifier{
		integrationKey: integrationKey,
		client:         &http.Client{Timeout: 10 * time.Second},
		eventURL:       pagerDutyEventURL,
	}, nil
}

// Notify sends a PagerDuty trigger event for the given secret.
func (p *PagerDutyNotifier) Notify(secret vault.Secret) error {
	days := secret.DaysUntilExpiration()
	var summary string
	if secret.IsExpired() {
		summary = fmt.Sprintf("[EXPIRED] Vault secret %q has expired", secret.Path)
	} else {
		summary = fmt.Sprintf("[EXPIRING] Vault secret %q expires in %d day(s)", secret.Path, days)
	}

	severity := "warning"
	if secret.IsExpired() {
		severity = "critical"
	}

	payload := pdPayload{
		RoutingKey:  p.integrationKey,
		EventAction: "trigger",
		Payload: pdDetails{
			Summary:   summary,
			Source:    "vaultwatch",
			Severity:  severity,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("pagerduty: failed to marshal payload: %w", err)
	}

	resp, err := p.client.Post(p.eventURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pagerduty: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
