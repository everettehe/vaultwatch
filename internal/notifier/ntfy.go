package notifier

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// NtfyNotifier sends notifications via ntfy.sh or a self-hosted ntfy server.
type NtfyNotifier struct {
	serverURL string
	topic     string
	client    *http.Client
}

// NewNtfyNotifier creates a new NtfyNotifier.
// serverURL defaults to "https://ntfy.sh" if empty.
func NewNtfyNotifier(serverURL, topic string) (*NtfyNotifier, error) {
	if topic == "" {
		return nil, fmt.Errorf("ntfy: topic is required")
	}
	if serverURL == "" {
		serverURL = "https://ntfy.sh"
	}
	serverURL = strings.TrimRight(serverURL, "/")
	return &NtfyNotifier{
		serverURL: serverURL,
		topic:     topic,
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends an alert for the given secret via ntfy.
func (n *NtfyNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	url := fmt.Sprintf("%s/%s", n.serverURL, n.topic)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(msg.Body))
	if err != nil {
		return fmt.Errorf("ntfy: failed to build request: %w", err)
	}
	req.Header.Set("Title", msg.Subject)
	req.Header.Set("Content-Type", "text/plain")
	if secret.IsExpired() {
		req.Header.Set("Priority", "urgent")
		req.Header.Set("Tags", "rotating_light")
	} else {
		req.Header.Set("Priority", "high")
		req.Header.Set("Tags", "warning")
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("ntfy: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
