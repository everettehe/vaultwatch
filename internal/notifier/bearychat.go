package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// BearyChat sends notifications to a BearyChat incoming webhook.
type BearyChat struct {
	webhookURL string
	client     *http.Client
}

// NewBearyChatNotifier creates a new BearyChat notifier.
func NewBearyChatNotifier(webhookURL string) (*BearyChat, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("bearychat: webhook URL is required")
	}
	return &BearyChat{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Notify sends a secret expiration alert to BearyChat.
func (b *BearyChat) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	payload := map[string]interface{}{
		"text":     msg.Subject,
		"markdown": true,
		"attachments": []map[string]string{
			{"text": msg.Body},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("bearychat: failed to marshal payload: %w", err)
	}

	resp, err := b.client.Post(b.webhookURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("bearychat: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bearychat: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
