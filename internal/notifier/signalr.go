package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SignalRNotifier sends alerts to an Azure SignalR-compatible HTTP endpoint.
type SignalRNotifier struct {
	endpointURL string
	accessKey   string
	hub         string
	client      *http.Client
}

type signalRPayload struct {
	Target    string        `json:"target"`
	Arguments []interface{} `json:"arguments"`
}

// NewSignalRNotifier creates a new SignalRNotifier.
func NewSignalRNotifier(endpointURL, accessKey, hub string) (*SignalRNotifier, error) {
	if endpointURL == "" {
		return nil, fmt.Errorf("signalr: endpoint URL is required")
	}
	if accessKey == "" {
		return nil, fmt.Errorf("signalr: access key is required")
	}
	if hub == "" {
		hub = "vaultwatch"
	}
	return &SignalRNotifier{
		endpointURL: endpointURL,
		accessKey:   accessKey,
		hub:         hub,
		client:      &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a secret expiration alert to the SignalR endpoint.
func (n *SignalRNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	payload := signalRPayload{
		Target:    "vaultAlert",
		Arguments: []interface{}{msg.Subject, msg.Body},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("signalr: failed to marshal payload: %w", err)
	}
	url := fmt.Sprintf("%s/api/v1/hubs/%s", n.endpointURL, n.hub)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("signalr: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+n.accessKey)
	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("signalr: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("signalr: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
