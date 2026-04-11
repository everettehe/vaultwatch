package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// CampfireNotifier sends alerts to a Basecamp Campfire room via webhook.
type CampfireNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewCampfireNotifier creates a new CampfireNotifier.
// webhookURL must be a valid Basecamp Campfire webhook URL.
func NewCampfireNotifier(webhookURL string) (*CampfireNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("campfire: webhook URL is required")
	}
	return &CampfireNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Notify sends a secret expiration alert to the configured Campfire room.
func (n *CampfireNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	payload := map[string]string{
		"content": fmt.Sprintf("%s: %s", msg.Subject, msg.Body),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("campfire: failed to marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("campfire: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("campfire: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
