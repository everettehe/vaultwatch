package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// GotifyNotifier sends notifications to a self-hosted Gotify server.
type GotifyNotifier struct {
	serverURL string
	token     string
	priority  int
	client    *http.Client
}

// NewGotifyNotifier creates a new GotifyNotifier.
// serverURL is the base URL of the Gotify instance, token is the app token.
func NewGotifyNotifier(serverURL, token string, priority int) (*GotifyNotifier, error) {
	if strings.TrimSpace(serverURL) == "" {
		return nil, fmt.Errorf("gotify: server URL is required")
	}
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("gotify: app token is required")
	}
	if priority <= 0 {
		priority = 5
	}
	return &GotifyNotifier{
		serverURL: strings.TrimRight(serverURL, "/"),
		token:     token,
		priority:  priority,
		client:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a message to Gotify for the given secret.
func (g *GotifyNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	params := url.Values{}
	params.Set("token", g.token)

	body := url.Values{}
	body.Set("title", msg.Subject)
	body.Set("message", msg.Body)
	body.Set("priority", fmt.Sprintf("%d", g.priority))

	endpoint := fmt.Sprintf("%s/message?%s", g.serverURL, params.Encode())
	resp, err := g.client.PostForm(endpoint, body)
	if err != nil {
		return fmt.Errorf("gotify: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gotify: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
