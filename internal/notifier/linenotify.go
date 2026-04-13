package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

const lineNotifyAPIURL = "https://notify-api.line.me/api/notify"

// LineNotifyNotifier sends alerts via LINE Notify.
type LineNotifyNotifier struct {
	token  string
	client *http.Client
}

// NewLineNotifyNotifier creates a new LineNotifyNotifier.
// token is the LINE Notify personal access token.
func NewLineNotifyNotifier(token string) (*LineNotifyNotifier, error) {
	if token == "" {
		return nil, fmt.Errorf("linenotify: token is required")
	}
	return &LineNotifyNotifier{
		token: token,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends a LINE Notify message for the given secret.
func (n *LineNotifyNotifier) Notify(secret vault.Secret) error {
	msg := FormatMessage(secret)

	form := url.Values{}
	form.Set("message", fmt.Sprintf("\n%s\n%s", msg.Subject, msg.Body))

	req, err := http.NewRequest(http.MethodPost, lineNotifyAPIURL, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("linenotify: failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+n.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("linenotify: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("linenotify: unexpected status %d", resp.StatusCode)
	}
	return nil
}
