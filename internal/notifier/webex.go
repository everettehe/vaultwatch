package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// WebexNotifier sends alerts to a Cisco Webex space via the Webex REST API.
type WebexNotifier struct {
	token  string
	roomID string
	client *http.Client
}

type webexMessage struct {
	RoomID   string `json:"roomId"`
	Markdown string `json:"markdown"`
}

// NewWebexNotifier creates a WebexNotifier. token is a Webex bot token and
// roomID is the destination space ID.
func NewWebexNotifier(token, roomID string) (*WebexNotifier, error) {
	if token == "" {
		return nil, fmt.Errorf("webex: token is required")
	}
	if roomID == "" {
		return nil, fmt.Errorf("webex: room_id is required")
	}
	return &WebexNotifier{
		token:  token,
		roomID: roomID,
		client: &http.Client{},
	}, nil
}

// Notify sends a Webex message for the given secret.
func (w *WebexNotifier) Notify(secret *vault.Secret) error {
	msg, err := FormatMessage(secret)
	if err != nil {
		return fmt.Errorf("webex: format message: %w", err)
	}

	body, err := json.Marshal(webexMessage{
		RoomID:   w.roomID,
		Markdown: fmt.Sprintf("**%s**\n\n%s", msg.Subject, msg.Body),
	})
	if err != nil {
		return fmt.Errorf("webex: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://webexapis.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webex: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.token)

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webex: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webex: unexpected status %d", resp.StatusCode)
	}
	return nil
}
