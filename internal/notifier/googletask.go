package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GoogleTaskNotifier sends vault secret expiration alerts by creating
// Google Tasks via a webhook-compatible endpoint (e.g., Apps Script).
type GoogleTaskNotifier struct {
	webhookURL string
	tasklist   string
	client     *http.Client
}

type googleTaskPayload struct {
	Title   string `json:"title"`
	Notes   string `json:"notes"`
	Due     string `json:"due"`
	Status  string `json:"status"`
}

// NewGoogleTaskNotifier creates a GoogleTaskNotifier.
// webhookURL is required; tasklist is optional (defaults to "@default").
func NewGoogleTaskNotifier(webhookURL, tasklist string) (*GoogleTaskNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("googletask: webhook URL is required")
	}
	if tasklist == "" {
		tasklist = "@default"
	}
	return &GoogleTaskNotifier{
		webhookURL: webhookURL,
		tasklist:   tasklist,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify creates a Google Task for the expiring/expired secret.
func (n *GoogleTaskNotifier) Notify(secret vault.Secret) error {
	msg := FormatMessage(secret)
	due := time.Now().Add(time.Duration(secret.DaysUntilExpiration()) * 24 * time.Hour).
		UTC().Format(time.RFC3339)

	payload := googleTaskPayload{
		Title:  msg.Subject,
		Notes:  msg.Body,
		Due:    due,
		Status: "needsAction",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googletask: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googletask: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("googletask: unexpected status %d", resp.StatusCode)
	}
	return nil
}
