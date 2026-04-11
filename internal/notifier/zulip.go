package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// ZulipNotifier sends alerts to a Zulip stream via the Zulip REST API.
type ZulipNotifier struct {
	baseURL  string
	email    string
	apiKey   string
	stream   string
	topic    string
	client   *http.Client
}

// NewZulipNotifier creates a new ZulipNotifier.
// baseURL is the Zulip server URL (e.g. https://yourorg.zulipchat.com),
// email and apiKey are bot credentials, stream and topic define the destination.
func NewZulipNotifier(baseURL, email, apiKey, stream, topic string) (*ZulipNotifier, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("zulip: base URL is required")
	}
	if email == "" {
		return nil, fmt.Errorf("zulip: bot email is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("zulip: API key is required")
	}
	if stream == "" {
		return nil, fmt.Errorf("zulip: stream is required")
	}
	if topic == "" {
		topic = "vaultwatch alerts"
	}
	return &ZulipNotifier{
		baseURL: baseURL,
		email:   email,
		apiKey:  apiKey,
		stream:  stream,
		topic:   topic,
		client:  &http.Client{},
	}, nil
}

// Notify sends a Zulip message for the given secret.
func (z *ZulipNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)

	form := url.Values{}
	form.Set("type", "stream")
	form.Set("to", z.stream)
	form.Set("topic", z.topic)
	form.Set("content", fmt.Sprintf("**%s**\n%s", msg.Subject, msg.Body))

	endpoint := z.baseURL + "/api/v1/messages"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("zulip: failed to build request: %w", err)
	}
	req.SetBasicAuth(z.email, z.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("zulip: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		return fmt.Errorf("zulip: unexpected status %d", resp.StatusCode)
	}
	return nil
}
