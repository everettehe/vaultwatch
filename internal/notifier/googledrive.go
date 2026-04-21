package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GoogleDriveNotifier appends secret expiration alerts to a Google Sheet
// via the Sheets API (treated as a drive-backed document).
type GoogleDriveNotifier struct {
	spreadsheetID string
	sheetName     string
	apiKey        string
	httpClient    *http.Client
	baseURL       string
}

type googleDriveAppendRequest struct {
	Values [][]interface{} `json:"values"`
}

// NewGoogleDriveNotifier creates a GoogleDriveNotifier.
// spreadsheetID and apiKey are required; sheetName defaults to "Alerts".
func NewGoogleDriveNotifier(spreadsheetID, sheetName, apiKey string) (*GoogleDriveNotifier, error) {
	if spreadsheetID == "" {
		return nil, fmt.Errorf("googledrive: spreadsheet ID is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("googledrive: API key is required")
	}
	if sheetName == "" {
		sheetName = "Alerts"
	}
	return &GoogleDriveNotifier{
		spreadsheetID: spreadsheetID,
		sheetName:     sheetName,
		apiKey:        apiKey,
		httpClient:    &http.Client{Timeout: 10 * time.Second},
		baseURL:       "https://sheets.googleapis.com/v4/spreadsheets",
	}, nil
}

// Notify appends a row to the configured Google Sheet.
func (n *GoogleDriveNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	row := []interface{}{
		time.Now().UTC().Format(time.RFC3339),
		secret.Path,
		fmt.Sprintf("%d", int(secret.DaysUntilExpiration())),
		msg.Body,
	}

	payload := googleDriveAppendRequest{Values: [][]interface{}{row}}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("googledrive: marshal error: %w", err)
	}

	url := fmt.Sprintf("%s/%s/values/%s:append?valueInputOption=RAW&key=%s",
		n.baseURL, n.spreadsheetID, n.sheetName, n.apiKey)

	resp, err := n.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googledrive: request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googledrive: unexpected status %d", resp.StatusCode)
	}
	return nil
}
