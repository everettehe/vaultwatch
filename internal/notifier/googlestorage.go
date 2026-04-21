package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GoogleStorageNotifier sends vault secret expiration events to a Google Cloud
// Storage bucket via the JSON API.
type GoogleStorageNotifier struct {
	bucketName string
	apiKey     string
	client     *http.Client
}

type storageObject struct {
	Path      string    `json:"path"`
	Status    string    `json:"status"`
	DaysLeft  int       `json:"days_until_expiration"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// NewGoogleStorageNotifier creates a GoogleStorageNotifier. bucketName and
// apiKey are required.
func NewGoogleStorageNotifier(bucketName, apiKey string) (*GoogleStorageNotifier, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("googlestorage: bucket name is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("googlestorage: api key is required")
	}
	return &GoogleStorageNotifier{
		bucketName: bucketName,
		apiKey:     apiKey,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify uploads a JSON object describing the expiring secret to the bucket.
func (n *GoogleStorageNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	obj := storageObject{
		Path:      secret.Path,
		Status:    msg.Status,
		DaysLeft:  int(secret.DaysUntilExpiration()),
		Message:   msg.Body,
		Timestamp: time.Now().UTC(),
	}

	body, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("googlestorage: marshal error: %w", err)
	}

	objectName := fmt.Sprintf("vaultwatch/%s/%d.json", secret.Path, time.Now().UnixNano())
	url := fmt.Sprintf(
		"https://storage.googleapis.com/upload/storage/v1/b/%s/o?uploadType=media&name=%s&key=%s",
		n.bucketName, objectName, n.apiKey,
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlestorage: request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("googlestorage: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlestorage: unexpected status %d", resp.StatusCode)
	}
	return nil
}
