package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// VictorOpsNotifier sends alerts to VictorOps (Splunk On-Call) via REST endpoint.
type VictorOpsNotifier struct {
	webhookURL string
	routingKey string
	client     *http.Client
}

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
	Timestamp         int64  `json:"timestamp"`
}

// NewVictorOpsNotifier creates a new VictorOpsNotifier.
// webhookURL should be the VictorOps REST endpoint URL including the routing key.
func NewVictorOpsNotifier(webhookURL, routingKey string) (*VictorOpsNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("victorops: webhook URL is required")
	}
	if routingKey == "" {
		return nil, fmt.Errorf("victorops: routing key is required")
	}
	return &VictorOpsNotifier{
		webhookURL: webhookURL,
		routingKey: routingKey,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a VictorOps alert for the given secret.
func (v *VictorOpsNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	messageType := "WARNING"
	if secret.IsExpired() {
		messageType = "CRITICAL"
	}

	payload := victorOpsPayload{
		MessageType:       messageType,
		EntityID:          fmt.Sprintf("vaultwatch/%s", secret.Path),
		EntityDisplayName: msg.Subject,
		StateMessage:      msg.Body,
		Timestamp:         time.Now().Unix(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s", v.webhookURL, v.routingKey)
	resp, err := v.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("victorops: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("victorops: unexpected status code %d: %s", resp.StatusCode, bytes.TrimSpace(respBody))
	}
	return nil
}
