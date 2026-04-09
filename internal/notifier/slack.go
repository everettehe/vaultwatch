package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SlackNotifier sends secret expiration alerts to a Slack webhook.
type SlackNotifier struct {
	webhookURL string
	channel    string
	client     *http.Client
}

type slackMessage struct {
	Channel string       `json:"channel,omitempty"`
	Blocks  []slackBlock `json:"blocks"`
}

type slackBlock struct {
	Type string    `json:"type"`
	Text *slackText `json:"text,omitempty"`
}

type slackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// NewSlackNotifier creates a new SlackNotifier with the given webhook URL and optional channel.
func NewSlackNotifier(webhookURL, channel string) (*SlackNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("slack webhook URL must not be empty")
	}
	return &SlackNotifier{
		webhookURL: webhookURL,
		channel:    channel,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a Slack message for the given secret expiration event.
func (s *SlackNotifier) Notify(secret vault.Secret) error {
	var text string
	if secret.IsExpired() {
		text = fmt.Sprintf(":rotating_light: *Secret Expired*\nPath: `%s`\nExpired at: %s",
			secret.Path, secret.ExpiresAt.Format(time.RFC3339))
	} else {
		days := secret.DaysUntilExpiration()
		text = fmt.Sprintf(":warning: *Secret Expiring Soon*\nPath: `%s`\nExpires in: *%d day(s)* (%s)",
			secret.Path, days, secret.ExpiresAt.Format(time.RFC3339))
	}

	msg := slackMessage{
		Channel: s.channel,
		Blocks: []slackBlock{
			{
				Type: "section",
				Text: &slackText{
					Type: "mrkdwn",
					Text: text,
				},
			},
		},
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("slack: failed to marshal message: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack: unexpected response status: %s", resp.Status)
	}
	return nil
}
