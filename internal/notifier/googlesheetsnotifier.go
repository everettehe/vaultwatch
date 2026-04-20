package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GoogleSheetsNotifier appends secret expiration events to a Google Sheets
// spreadsheet via a Google Apps Script Web App URL.
type GoogleSheetsNotifier struct {
	webAppURL string
	client    *http.Client
}

// NewGoogleSheetsNotifier creates a new GoogleSheetsNotifier.
func NewGoogleSheetsNotifier(webAppURL string) (*GoogleSheetsNotifier, error) {
	if webAppURL == "" {
		return nil, fmt.Errorf("googlesheetsnotifier: web app URL is required")
	}
	return &GoogleSheetsNotifier{
		webAppURL: webAppURL,
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends secret expiration data to the configured Google Sheets web app.
func (n *GoogleSheetsNotifier) Notify(secret vault.Secret) error {
	payload := map[string]interface{}{
		"path":       secret.Path,
		"days":       secret.DaysUntilExpiration(),
		"expires_at": secret.ExpiresAt.Format(time.RFC3339),
		"expired":    secret.IsExpired(),
		"message":    FormatMessage(secret),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlesheetsnotifier: failed to marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webAppURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlesheetsnotifier: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlesheetsnotifier: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
