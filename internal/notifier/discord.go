package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// DiscordNotifier sends secret expiration alerts to a Discord channel
// via an incoming webhook.
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewDiscordNotifier creates a new DiscordNotifier. It returns an error if
// webhookURL is empty.
func NewDiscordNotifier(webhookURL string) (*DiscordNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("discord: webhook URL is required")
	}
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

type discordPayload struct {
	Content string `json:"content"`
}

// Notify sends a formatted alert message to the configured Discord webhook.
func (d *DiscordNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	payload := discordPayload{
		Content: fmt.Sprintf("**%s**\n%s", msg.Subject, msg.Body),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: failed to marshal payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("discord: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
