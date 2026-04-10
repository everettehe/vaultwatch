package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourusername/vaultwatch/internal/vault"
)

const telegramAPIBase = "https://api.telegram.org/bot"

// TelegramNotifier sends alerts via Telegram Bot API.
type TelegramNotifier struct {
	botToken string
	chatID   string
	client   *http.Client
}

type telegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// NewTelegramNotifier creates a TelegramNotifier.
// botToken and chatID are required.
func NewTelegramNotifier(botToken, chatID string) (*TelegramNotifier, error) {
	if botToken == "" {
		return nil, fmt.Errorf("telegram: bot token is required")
	}
	if chatID == "" {
		return nil, fmt.Errorf("telegram: chat ID is required")
	}
	return &TelegramNotifier{
		botToken: botToken,
		chatID:   chatID,
		client:   &http.Client{},
	}, nil
}

// Notify sends a Telegram message for the given secret event.
func (t *TelegramNotifier) Notify(s *vault.Secret) error {
	msg := FormatMessage(s)
	payload := telegramMessage{
		ChatID:    t.chatID,
		Text:      fmt.Sprintf("*%s*\n%s", msg.Subject, msg.Body),
		ParseMode: "Markdown",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram: marshal payload: %w", err)
	}
	url := fmt.Sprintf("%s%s/sendMessage", telegramAPIBase, t.botToken)
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram: send message: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}
	return nil
}
