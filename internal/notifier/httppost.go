package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// httpPostNotifier sends a JSON POST request to a generic HTTP endpoint.
// It is similar to WebhookNotifier but always uses POST and includes
// structured JSON rather than a plain text body.
type httpPostNotifier struct {
	url    string
	headers map[string]string
	client *http.Client
}

// NewHTTPPostNotifier creates a notifier that POSTs structured JSON to url.
// Optional headers (e.g. Authorization) may be supplied via the headers map.
func NewHTTPPostNotifier(url string, headers map[string]string) (*httpPostNotifier, error) {
	if url == "" {
		return nil, fmt.Errorf("httppost: url is required")
	}
	if headers == nil {
		headers = map[string]string{}
	}
	return &httpPostNotifier{
		url:     url,
		headers: headers,
		client:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (n *httpPostNotifier) Notify(s *vault.Secret) error {
	msg := FormatMessage(s)
	payload := map[string]string{
		"path":    s.Path,
		"status":  msg.Status,
		"subject": msg.Subject,
		"body":    msg.Body,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("httppost: marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, n.url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("httppost: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range n.headers {
		req.Header.Set(k, v)
	}
	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("httppost: do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("httppost: unexpected status %d", resp.StatusCode)
	}
	return nil
}
