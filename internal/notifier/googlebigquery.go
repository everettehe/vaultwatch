package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// BigQueryNotifier sends secret expiration events to a Google BigQuery table
// via the BigQuery streaming insert REST API.
type BigQueryNotifier struct {
	projectID string
	datasetID string
	tableID   string
	apiKey    string
	client    *http.Client
}

// NewBigQueryNotifier constructs a BigQueryNotifier from the provided config.
func NewBigQueryNotifier(projectID, datasetID, tableID, apiKey string) (*BigQueryNotifier, error) {
	if projectID == "" {
		return nil, fmt.Errorf("bigquery: project_id is required")
	}
	if datasetID == "" {
		return nil, fmt.Errorf("bigquery: dataset_id is required")
	}
	if tableID == "" {
		return nil, fmt.Errorf("bigquery: table_id is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("bigquery: api_key is required")
	}
	return &BigQueryNotifier{
		projectID: projectID,
		datasetID: datasetID,
		tableID:   tableID,
		apiKey:    apiKey,
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a streaming insert to BigQuery with the secret expiration details.
func (n *BigQueryNotifier) Notify(ctx context.Context, secret vault.Secret) error {
	row := map[string]interface{}{
		"path":        secret.Path,
		"days_left":   secret.DaysUntilExpiration(),
		"expires_at":  secret.ExpiresAt.Format(time.RFC3339),
		"is_expired":  secret.IsExpired(),
		"notified_at": time.Now().UTC().Format(time.RFC3339),
	}

	body := map[string]interface{}{
		"rows": []map[string]interface{}{
			{"insertId": fmt.Sprintf("%s-%d", secret.Path, time.Now().UnixNano()), "json": row},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("bigquery: failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf(
		"https://bigquery.googleapis.com/bigquery/v2/projects/%s/datasets/%s/tables/%s/insertAll?key=%s",
		n.projectID, n.datasetID, n.tableID, n.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("bigquery: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("bigquery: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("bigquery: unexpected status %d", resp.StatusCode)
	}
	return nil
}
