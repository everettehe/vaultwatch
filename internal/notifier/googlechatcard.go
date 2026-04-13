package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GoogleChatCardNotifier sends rich card messages to a Google Chat webhook.
type GoogleChatCardNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewGoogleChatCardNotifier creates a new GoogleChatCardNotifier.
func NewGoogleChatCardNotifier(webhookURL string) (*GoogleChatCardNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("googlechatcard: webhook URL is required")
	}
	return &GoogleChatCardNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Notify sends a structured card message to Google Chat.
func (n *GoogleChatCardNotifier) Notify(secret *vault.Secret) error {
	msg := n.buildCard(secret)

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("googlechatcard: failed to marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("googlechatcard: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("googlechatcard: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (n *GoogleChatCardNotifier) buildCard(secret *vault.Secret) map[string]any {
	subj, body := FormatMessage(secret)
	color := "#FFA500"
	if secret.IsExpired() {
		color = "#FF0000"
	}
	return map[string]any{
		"cards": []map[string]any{
			{
				"header": map[string]any{
					"title":    subj,
					"subtitle": secret.Path,
				},
				"sections": []map[string]any{
					{
						"widgets": []map[string]any{
							{
								"textParagraph": map[string]any{
									"text": fmt.Sprintf("<font color=\"%s\">%s</font>", color, body),
								},
							},
						},
					},
				},
			},
		},
	}
}
