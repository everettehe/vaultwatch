package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GoogleStorageNotifier uploads alert payloads to a GCS bucket via the
// JSON API.
type GoogleStorageNotifier struct {
	bucket string
	apiKey string
	client *http.Client
}

// NewGoogleStorageNotifier creates a GoogleStorageNotifier.
func NewGoogleStorageNotifier(bucket, apiKey string) (*GoogleStorageNotifier, error) {
	if bucket == "" {
		return nil, fmt.Errorf("google_storage: bucket is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("google_storage: api_key is required")
	}
	return &GoogleStorageNotifier{
		bucket: bucket,
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify uploads a JSON record for the given secret to GCS.
func (n *GoogleStorageNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	payload, err := json.Marshal(map[string]string{
		"path":    secret.Path,
		"message": msg.Body,
		"status":  msg.Subject,
	})
	if err != nil {
		return fmt.Errorf("google_storage: marshal payload: %w", err)
	}

	objectName := fmt.Sprintf("vaultwatch/%s/%d.json",
		secret.Path, time.Now().UTC().UnixNano())

	url := fmt.Sprintf(
		"https://storage.googleapis.com/upload/storage/v1/b/%s/o?uploadType=media&name=%s&key=%s",
		n.bucket, objectName, n.apiKey,
	)

	resp, err := n.client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("google_storage: upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("google_storage: unexpected status %d", resp.StatusCode)
	}
	return nil
}
