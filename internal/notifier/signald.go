package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SignaldNotifier sends notifications via a local signald HTTP bridge.
type SignaldNotifier struct {
	baseURL    string
	sender     string
	recipients []string
	client     *http.Client
}

type signaldPayload struct {
	Username   string   `json:"username"`
	Recipients []string `json:"recipients"`
	Message    string   `json:"message"`
}

// NewSignaldNotifier creates a new SignaldNotifier.
// baseURL is the signald HTTP bridge URL (e.g. http://localhost:8080),
// sender is the registered Signal number, and recipients is a list of
// phone numbers or group IDs to notify.
func NewSignaldNotifier(baseURL, sender string, recipients []string) (*SignaldNotifier, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("signald: base URL is required")
	}
	if sender == "" {
		return nil, fmt.Errorf("signald: sender number is required")
	}
	if len(recipients) == 0 {
		return nil, fmt.Errorf("signald: at least one recipient is required")
	}
	return &SignaldNotifier{
		baseURL:    baseURL,
		sender:     sender,
		recipients: recipients,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a Signal message for the given secret via signald.
func (n *SignaldNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	payload := signaldPayload{
		Username:   n.sender,
		Recipients: n.recipients,
		Message:    fmt.Sprintf("%s\n%s", msg.Subject, msg.Body),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signald: failed to marshal payload: %w", err)
	}
	resp, err := n.client.Post(n.baseURL+"/v1/send", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signald: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("signald: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
