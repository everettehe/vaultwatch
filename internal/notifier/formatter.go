package notifier

import (
	"fmt"
	"time"

	"github.com/your-org/vaultwatch/internal/vault"
)

// MessageLevel indicates the urgency of a notification message.
type MessageLevel string

const (
	LevelWarning  MessageLevel = "WARNING"
	LevelCritical MessageLevel = "CRITICAL"
	LevelExpired  MessageLevel = "EXPIRED"
)

// Message holds the formatted notification content for a secret event.
type Message struct {
	Level   MessageLevel
	Subject string
	Body    string
	Secret  *vault.Secret
}

// FormatMessage builds a human-readable Message for the given secret.
func FormatMessage(s *vault.Secret, warningThreshold time.Duration) Message {
	days := s.DaysUntilExpiration()

	var level MessageLevel
	switch {
	case s.IsExpired():
		level = LevelExpired
	case days <= 1:
		level = LevelCritical
	default:
		level = LevelWarning
	}

	subject := formatSubject(level, s.Path)
	body := formatBody(level, s.Path, days, s.ExpiresAt)

	return Message{
		Level:   level,
		Subject: subject,
		Body:    body,
		Secret:  s,
	}
}

func formatSubject(level MessageLevel, path string) string {
	switch level {
	case LevelExpired:
		return fmt.Sprintf("[EXPIRED] Vault secret expired: %s", path)
	case LevelCritical:
		return fmt.Sprintf("[CRITICAL] Vault secret expiring within 24h: %s", path)
	default:
		return fmt.Sprintf("[WARNING] Vault secret expiring soon: %s", path)
	}
}

func formatBody(level MessageLevel, path string, days int, expiresAt time.Time) string {
	switch level {
	case LevelExpired:
		return fmt.Sprintf(
			"Vault secret at path %q has expired as of %s. Immediate rotation required.",
			path, expiresAt.UTC().Format(time.RFC1123),
		)
	case LevelCritical:
		return fmt.Sprintf(
			"Vault secret at path %q expires at %s (less than 24 hours remaining). Rotate immediately.",
			path, expiresAt.UTC().Format(time.RFC1123),
		)
	default:
		return fmt.Sprintf(
			"Vault secret at path %q expires at %s (%d days remaining). Please schedule rotation.",
			path, expiresAt.UTC().Format(time.RFC1123), days,
		)
	}
}
