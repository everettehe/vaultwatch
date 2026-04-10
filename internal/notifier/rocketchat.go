package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// RocketChatNotifier sends alerts to a Rocket.Chat channel via incoming webhook.
type RocketChatNotifier struct {
	webhookURL string
	channel    string
	client     *http.Client
}

type rocketChatPayload struct {
	Text    string `json:"text"`
	Channel string `json:"channel,omitempty"`
}

// NewRocketChatNotifier creates a new RocketChatNotifier.
// webhookURL is required; channel is optional (uses webhook default if empty).
func NewRocketChatNotifier(webhookURL, channel string) (*RocketChatNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("rocketchat: webhook URL is required")
	}
	return &RocketChatNotifier{
		webhookURL: webhookURL,
		channel:    channel,
		client:     &http.Client{},
	}, nil
}

// Notify sends a Rocket.Chat message for the given secret.
func (r *RocketChatNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	payload := rocketChatPayload{
		Text:    fmt.Sprintf("%s\n%s", msg.Subject, msg.Body),
		Channel: r.channel,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rocketchat: failed to marshal payload: %w", err)
	}

	resp, err := r.client.Post(r.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("rocketchat: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rocketchat: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
