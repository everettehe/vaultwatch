package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SplunkNotifier sends alerts to a Splunk HTTP Event Collector (HEC) endpoint.
type SplunkNotifier struct {
	url        string
	token      string
	source     string
	sourceType string
	index      string
	client     *http.Client
}

type splunkEvent struct {
	Time       float64        `json:"time"`
	Source     string         `json:"source,omitempty"`
	SourceType string         `json:"sourcetype,omitempty"`
	Index      string         `json:"index,omitempty"`
	Event      map[string]any `json:"event"`
}

// NewSplunkNotifier creates a SplunkNotifier that posts to the given HEC URL.
// token is the Splunk HEC token. source, sourceType, and index are optional.
func NewSplunkNotifier(url, token, source, sourceType, index string) (*SplunkNotifier, error) {
	if url == "" {
		return nil, fmt.Errorf("splunk: HEC url is required")
	}
	if token == "" {
		return nil, fmt.Errorf("splunk: HEC token is required")
	}
	return &SplunkNotifier{
		url:        url,
		token:      token,
		source:     source,
		sourceType: sourceType,
		index:      index,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a Vault secret expiration event to Splunk HEC.
func (s *SplunkNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	event := splunkEvent{
		Time:       float64(time.Now().Unix()),
		Source:     s.source,
		SourceType: s.sourceType,
		Index:      s.index,
		Event: map[string]any{
			"message":          msg.Body,
			"secret_path":      secret.Path,
			"days_until_expiry": secret.DaysUntilExpiration(),
			"severity":         msg.Severity,
		},
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("splunk: failed to marshal event: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("splunk: failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Splunk "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
