package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// StatusPageNotifier sends incident updates to Atlassian Statuspage.
type StatusPageNotifier struct {
	apiKey    string
	pageID    string
	baseURL   string
	client    *http.Client
}

type statusPageIncident struct {
	Incident statusPageIncidentBody `json:"incident"`
}

type statusPageIncidentBody struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Body   string `json:"body"`
}

// NewStatusPageNotifier creates a StatusPageNotifier.
func NewStatusPageNotifier(apiKey, pageID string) (*StatusPageNotifier, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("statuspage: api key is required")
	}
	if pageID == "" {
		return nil, fmt.Errorf("statuspage: page ID is required")
	}
	return &StatusPageNotifier{
		apiKey:  apiKey,
		pageID:  pageID,
		baseURL: "https://api.statuspage.io/v1",
		client:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify creates a Statuspage incident for the given secret.
func (n *StatusPageNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	status := "investigating"
	if secret.IsExpired() {
		status = "identified"
	}

	payload := statusPageIncident{
		Incident: statusPageIncidentBody{
			Name:   msg.Subject,
			Status: status,
			Body:   msg.Body,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("statuspage: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/pages/%s/incidents", n.baseURL, n.pageID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("statuspage: create request: %w", err)
	}
	req.Header.Set("Authorization", "OAuth "+n.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("statuspage: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("statuspage: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
