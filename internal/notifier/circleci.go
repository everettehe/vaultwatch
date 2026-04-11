package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// CircleCINotifier sends Vault secret expiration alerts as CircleCI pipeline
// environment variable update requests or as a simple webhook-style notification
// via the CircleCI API.
type CircleCINotifier struct {
	token   string
	baseURL string
	client  *http.Client
}

type circleCIPayload struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// NewCircleCINotifier creates a new CircleCINotifier.
// token is the CircleCI personal API token.
// baseURL is the CircleCI API base URL (e.g. https://circleci.com/api/v2).
func NewCircleCINotifier(token, baseURL string) (*CircleCINotifier, error) {
	if token == "" {
		return nil, fmt.Errorf("circleci: token is required")
	}
	if baseURL == "" {
		baseURL = "https://circleci.com/api/v2"
	}
	return &CircleCINotifier{
		token:   token,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}

// Notify sends a secret expiration alert to the CircleCI notification endpoint.
func (n *CircleCINotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	payload := circleCIPayload{
		Subject: msg.Subject,
		Body:    msg.Body,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("circleci: failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/me/notifications", n.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("circleci: failed to create request: %w", err)
	}
	req.Header.Set("Circle-Token", n.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("circleci: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("circleci: unexpected status %d", resp.StatusCode)
	}
	return nil
}
