package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/your-org/vaultwatch/internal/vault"
)

const defaultOpsGenieURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieNotifier sends alerts to OpsGenie.
type OpsGenieNotifier struct {
	apiKey  string
	apiURL  string
	client  *http.Client
}

// NewOpsGenieNotifier creates a new OpsGenieNotifier.
func NewOpsGenieNotifier(apiKey, apiURL string) (*OpsGenieNotifier, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("opsgenie: api key is required")
	}
	if apiURL == "" {
		apiURL = defaultOpsGenieURL
	}
	return &OpsGenieNotifier{
		apiKey: apiKey,
		apiURL: apiURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends an alert to OpsGenie for the given secret.
func (n *OpsGenieNotifier) Notify(secret vault.Secret) error {
	var message string
	if secret.IsExpired() {
		message = fmt.Sprintf("Vault secret EXPIRED: %s (expired %d days ago)",
			secret.Path, -secret.DaysUntilExpiration())
	} else {
		message = fmt.Sprintf("Vault secret expiring soon: %s (expires in %d days)",
			secret.Path, secret.DaysUntilExpiration())
	}

	priority := "P3"
	if secret.IsExpired() {
		priority = "P1"
	}

	payload := map[string]interface{}{
		"message":  message,
		"priority": priority,
		"alias":    fmt.Sprintf("vaultwatch-%s", secret.Path),
		"tags":     []string{"vaultwatch", "vault", "secret-rotation"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+n.apiKey)

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
