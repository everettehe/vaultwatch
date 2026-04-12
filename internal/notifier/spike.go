package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SpikeNotifier sends alerts to Spike.sh via their webhook API.
type SpikeNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewSpikeNotifier creates a new SpikeNotifier.
// webhookURL is required and must be a valid Spike.sh webhook URL.
func NewSpikeNotifier(webhookURL string) (*SpikeNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("spike: webhook URL is required")
	}
	return &SpikeNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

type spikePayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Notify sends a secret expiration alert to Spike.sh.
func (n *SpikeNotifier) Notify(secret vault.Secret) error {
	msg := FormatMessage(secret)

	status := "warning"
	if secret.IsExpired() {
		status = "critical"
	}

	payload := spikePayload{
		Title:   msg.Subject,
		Message: msg.Body,
		Status:  status,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("spike: failed to marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("spike: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("spike: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
