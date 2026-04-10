package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// MattermostNotifier sends alerts to a Mattermost channel via incoming webhook.
type MattermostNotifier struct {
	webhookURL string
	channel    string
	username   string
	client     *http.Client
}

type mattermostPayload struct {
	Text     string `json:"text"`
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username,omitempty"`
}

// NewMattermostNotifier creates a MattermostNotifier.
// webhookURL is required; channel and username are optional overrides.
func NewMattermostNotifier(webhookURL, channel, username string) (*MattermostNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("mattermost: webhook URL is required")
	}
	return &MattermostNotifier{
		webhookURL: webhookURL,
		channel:    channel,
		username:   username,
		client:     &http.Client{},
	}, nil
}

// Notify sends a secret expiration alert to Mattermost.
func (m *MattermostNotifier) Notify(secret *vault.Secret) error {
	msg, err := FormatMessage(secret)
	if err != nil {
		return fmt.Errorf("mattermost: format message: %w", err)
	}

	payload := mattermostPayload{
		Text:     msg.Body,
		Channel:  m.channel,
		Username: m.username,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: marshal payload: %w", err)
	}

	resp, err := m.client.Post(m.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mattermost: post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
