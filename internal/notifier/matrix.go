package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// MatrixNotifier sends alerts to a Matrix room via the Matrix HTTP API.
type MatrixNotifier struct {
	homeserver string
	token      string
	roomID     string
	client     *http.Client
}

type matrixMessage struct {
	MsgType       string `json:"msgtype"`
	Body          string `json:"body"`
	FormattedBody string `json:"formatted_body,omitempty"`
	Format        string `json:"format,omitempty"`
}

// NewMatrixNotifier creates a MatrixNotifier. homeserver should be the base URL
// (e.g. https://matrix.example.com), token is the access token, and roomID is
// the fully-qualified Matrix room ID (e.g. !abc123:example.com).
func NewMatrixNotifier(homeserver, token, roomID string) (*MatrixNotifier, error) {
	if strings.TrimSpace(homeserver) == "" {
		return nil, fmt.Errorf("matrix: homeserver URL is required")
	}
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("matrix: access token is required")
	}
	if strings.TrimSpace(roomID) == "" {
		return nil, fmt.Errorf("matrix: room ID is required")
	}
	return &MatrixNotifier{
		homeserver: strings.TrimRight(homeserver, "/"),
		token:      token,
		roomID:     roomID,
		client:     &http.Client{},
	}, nil
}

// Notify sends a secret expiration alert to the configured Matrix room.
func (m *MatrixNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	payload := matrixMessage{
		MsgType: "m.text",
		Body:    msg.Subject + "\n" + msg.Body,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("matrix: failed to marshal message: %w", err)
	}

	// Matrix event send endpoint
	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/send/m.room.message",
		m.homeserver, m.roomID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("matrix: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("matrix: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("matrix: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
