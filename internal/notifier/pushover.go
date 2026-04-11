package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/yourusername/vaultwatch/internal/vault"
)

const pushoverAPIURL = "https://api.push.json"

// PushoverNotifier sends notifications via the Pushover API.
type PushoverNotifier struct {
	userKey  string
	apiToken string
	client   *http.Client
}

// NewPushoverNotifier creates a new PushoverNotifier.
// userKey is the Pushover user or group key.
// apiToken is the application API token.
func NewPushoverNotifier(userKey, apiToken string) (*PushoverNotifier, error) {
	if userKey == "" {
		return nil, fmt.Errorf("pushover: user key is required")
	}
	if apiToken == "" {
		return nil, fmt.Errorf("pushover: api token is required")
	}
	return &PushoverNotifier{
		userKey:  userKey,
		apiToken: apiToken,
		client:   &http.Client{},
	}, nil
}

// Notify sends a Pushover notification for the given secret.
func (p *PushoverNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	form := url.Values{}
	form.Set("token", p.apiToken)
	form.Set("user", p.userKey)
	form.Set("title", msg.Subject)
	form.Set("message", msg.Body)

	if secret.IsExpired() {
		form.Set("priority", "1")
	} else {
		form.Set("priority", "0")
	}

	resp, err := p.client.Post(pushoverAPIURL, "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("pushover: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pushover: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
