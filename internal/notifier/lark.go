package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// LarkNotifier sends alerts to a Lark (Feishu) channel via incoming webhook.
type LarkNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewLarkNotifier creates a new LarkNotifier.
func NewLarkNotifier(webhookURL string) (*LarkNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("lark: webhook URL is required")
	}
	return &LarkNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Notify sends a Lark message for the given secret.
func (n *LarkNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	payload := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": fmt.Sprintf("%s\n%s", msg.Subject, msg.Body),
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("lark: failed to marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("lark: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lark: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
