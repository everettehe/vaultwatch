package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GoogleChatNotifier sends alerts to a Google Chat space via an incoming webhook.
type GoogleChatNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewGoogleChatNotifier creates a new GoogleChatNotifier.
// webhookURL must be a valid Google Chat incoming webhook URL.
func NewGoogleChatNotifier(webhookURL string) (*GoogleChatNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("googlechat: webhook URL is required")
	}
	return &GoogleChatNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Notify sends a Google Chat message for the given secret.
func (n *GoogleChatNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	payload := map[string]string{
		"text": fmt.Sprintf("*%s*\n%s", msg.Subject, msg.Body),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlechat: failed to marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechat: failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
