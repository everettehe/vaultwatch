package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GoogleSheetsNotifier appends secret expiration alerts to a Google Sheet
// via a Google Apps Script Web App endpoint.
type GoogleSheetsNotifier struct {
	webAppURL string
	sheetName string
	client    *http.Client
}

type sheetsPayload struct {
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
	Message   string `json:"message"`
	Severity  string `json:"severity"`
	Sheet     string `json:"sheet,omitempty"`
}

// NewGoogleSheetsNotifier creates a GoogleSheetsNotifier.
// webAppURL is the deployed Apps Script endpoint; sheetName is optional.
func NewGoogleSheetsNotifier(webAppURL, sheetName string) (*GoogleSheetsNotifier, error) {
	if webAppURL == "" {
		return nil, fmt.Errorf("googlesheetsnotifier: web app URL is required")
	}
	return &GoogleSheetsNotifier{
		webAppURL: webAppURL,
		sheetName: sheetName,
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a row to the configured Google Sheet.
func (n *GoogleSheetsNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	severity := "WARNING"
	if secret.IsExpired() {
		severity = "EXPIRED"
	} else if secret.DaysUntilExpiration() <= 3 {
		severity = "CRITICAL"
	}

	payload := sheetsPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Path:      secret.Path,
		Message:   msg.Body,
		Severity:  severity,
		Sheet:     n.sheetName,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googlesheetsnotifier: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webAppURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlesheetsnotifier: post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlesheetsnotifier: unexpected status %d", resp.StatusCode)
	}
	return nil
}
