package notifier

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// TwilioNotifier sends SMS alerts via the Twilio REST API.
type TwilioNotifier struct {
	accountSID string
	authToken  string
	from       string
	to         string
	client     *http.Client
}

// NewTwilioNotifier creates a new TwilioNotifier.
func NewTwilioNotifier(accountSID, authToken, from, to string) (*TwilioNotifier, error) {
	if accountSID == "" {
		return nil, fmt.Errorf("twilio: account_sid is required")
	}
	if authToken == "" {
		return nil, fmt.Errorf("twilio: auth_token is required")
	}
	if from == "" {
		return nil, fmt.Errorf("twilio: from number is required")
	}
	if to == "" {
		return nil, fmt.Errorf("twilio: to number is required")
	}
	return &TwilioNotifier{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		to:         to,
		client:     &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify sends an SMS message for the given secret.
func (n *TwilioNotifier) Notify(s vault.Secret) error {
	msg, _ := FormatMessage(s)
	body := url.Values{}
	body.Set("From", n.from)
	body.Set("To", n.to)
	body.Set("Body", msg.Subject+": "+msg.Body)

	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", n.accountSID)
	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(body.Encode()))
	if err != nil {
		return fmt.Errorf("twilio: build request: %w", err)
	}
	req.SetBasicAuth(n.accountSID, n.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("twilio: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		var e struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(raw, &e)
		return fmt.Errorf("twilio: unexpected status %d: %s", resp.StatusCode, e.Message)
	}
	return nil
}
