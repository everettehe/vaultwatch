package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SignaldNotifier sends notifications via the signald HTTP API (Signal messenger).
type SignaldNotifier struct {
	baseURL    string
	sender     string
	recipients []string
	client     *http.Client
}

type signaldMessage struct {
	Username     string   `json:"username"`
	Recipients   []string `json:"recipients"`
	MessageBody  string   `json:"messageBody"`
}

// NewSignaldNotifier constructs a SignaldNotifier.
func NewSignaldNotifier(baseURL, sender string, recipients []string) (*SignaldNotifier, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("signald: base URL is required")
	}
	if sender == "" {
		return nil, fmt.Errorf("signald: sender is required")
	}
	if len(recipients) == 0 {
		return nil, fmt.Errorf("signald: at least one recipient is required")
	}
	return &SignaldNotifier{
		baseURL:    baseURL,
		sender:     sender,
		recipients: recipients,
		client:     &http.Client{},
	}, nil
}

// Notify sends a Signal message for the given secret.
func (n *SignaldNotifier) Notify(s *vault.Secret) error {
	msg, _ := FormatMessage(s)
	payload := signaldMessage{
		Username:    n.sender,
		Recipients:  n.recipients,
		MessageBody: msg,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signald: failed to marshal payload: %w", err)
	}
	resp, err := n.client.Post(n.baseURL+"/v2/send", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signald: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("signald: unexpected status %d", resp.StatusCode)
	}
	return nil
}
