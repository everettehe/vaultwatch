package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

const clickSendAPIURL = "https://rest.clicksend.com/v3/sms/send"

// ClickSendNotifier sends SMS notifications via ClickSend.
type ClickSendNotifier struct {
	username string
	apiKey   string
	to       string
	client   *http.Client
}

// NewClickSendNotifier creates a new ClickSendNotifier.
func NewClickSendNotifier(username, apiKey, to string) (*ClickSendNotifier, error) {
	if username == "" {
		return nil, fmt.Errorf("clicksend: username is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("clicksend: api key is required")
	}
	if to == "" {
		return nil, fmt.Errorf("clicksend: recipient phone number is required")
	}
	return &ClickSendNotifier{
		username: username,
		apiKey:   apiKey,
		to:       to,
		client:   &http.Client{},
	}, nil
}

// Notify sends an SMS notification for the given secret.
func (n *ClickSendNotifier) Notify(secret *vault.Secret) error {
	msg, err := FormatMessage(secret)
	if err != nil {
		return fmt.Errorf("clicksend: format message: %w", err)
	}

	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"source": "vaultwatch", "to": n.to, "body": msg.Body},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("clicksend: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, clickSendAPIURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("clicksend: create request: %w", err)
	}
	req.SetBasicAuth(n.username, n.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("clicksend: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("clicksend: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
