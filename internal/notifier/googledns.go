package notifier

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GoogleDNSNotifier posts secret expiration alerts via a Google Cloud DNS
// webhook endpoint (e.g. a Cloud Function fronting DNS metadata).
type GoogleDNSNotifier struct {
	webhookURL string
	project    string
	client     *http.Client
}

// NewGoogleDNSNotifier creates a GoogleDNSNotifier.
// webhookURL and project are required.
func NewGoogleDNSNotifier(webhookURL, project string) (*GoogleDNSNotifier, error) {
	if strings.TrimSpace(webhookURL) == "" {
		return nil, fmt.Errorf("googledns: webhook URL is required")
	}
	if strings.TrimSpace(project) == "" {
		return nil, fmt.Errorf("googledns: project is required")
	}
	return &GoogleDNSNotifier{
		webhookURL: webhookURL,
		project:    project,
		client:     &http.Client{},
	}, nil
}

// Notify sends a formatted alert message to the configured webhook.
func (n *GoogleDNSNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	body := fmt.Sprintf(`{"project":%q,"message":%q}`, n.project, msg.Body)
	resp, err := n.client.Post(n.webhookURL, "application/json", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("googledns: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("googledns: unexpected status %d", resp.StatusCode)
	}
	return nil
}
