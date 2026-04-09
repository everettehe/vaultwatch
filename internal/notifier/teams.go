package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// TeamsNotifier sends alerts to a Microsoft Teams channel via an incoming webhook.
type TeamsNotifier struct {
	webhookURL string
	client     *http.Client
}

type teamsPayload struct {
	Type       string         `json:"@type"`
	Context    string         `json:"@context"`
	ThemeColor string         `json:"themeColor"`
	Summary    string         `json:"summary"`
	Sections   []teamsSection `json:"sections"`
}

type teamsSection struct {
	ActivityTitle string      `json:"activityTitle"`
	Facts         []teamsFact `json:"facts"`
	Markdown      bool        `json:"markdown"`
}

type teamsFact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewTeamsNotifier creates a TeamsNotifier with the given webhook URL.
func NewTeamsNotifier(webhookURL string) (*TeamsNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("teams webhook URL is required")
	}
	return &TeamsNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a Teams message for the given secret.
func (t *TeamsNotifier) Notify(secret vault.Secret) error {
	color := "FFA500"
	status := fmt.Sprintf("Expiring in %d day(s)", secret.DaysUntilExpiration())
	if secret.IsExpired() {
		color = "FF0000"
		status = "Expired"
	}

	payload := teamsPayload{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: color,
		Summary:    fmt.Sprintf("VaultWatch Alert: %s", secret.Path),
		Sections: []teamsSection{
			{
				ActivityTitle: "VaultWatch Secret Alert",
				Markdown:      true,
				Facts: []teamsFact{
					{Name: "Path", Value: secret.Path},
					{Name: "Status", Value: status},
					{Name: "Expiration", Value: secret.ExpiresAt.Format(time.RFC3339)},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
