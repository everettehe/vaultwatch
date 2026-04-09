package notifier

import (
	"fmt"
	"log"
	"time"

	"vaultwatch/internal/vault"
)

// Notifier defines the interface for sending notifications
type Notifier interface {
	Notify(secret vault.Secret) error
}

// LogNotifier sends notifications to stdout/stderr
type LogNotifier struct {
	logger *log.Logger
}

// NewLogNotifier creates a new log-based notifier
func NewLogNotifier(logger *log.Logger) *LogNotifier {
	if logger == nil {
		logger = log.Default()
	}
	return &LogNotifier{
		logger: logger,
	}
}

// Notify sends a notification via logging
func (l *LogNotifier) Notify(secret vault.Secret) error {
	days := secret.DaysUntilExpiration()
	var message string

	if secret.IsExpired() {
		message = fmt.Sprintf("[EXPIRED] Secret '%s' expired %d days ago",
			secret.Path, -days)
	} else {
		message = fmt.Sprintf("[WARNING] Secret '%s' expires in %d days (on %s)",
			secret.Path, days, secret.ExpirationTime.Format(time.RFC3339))
	}

	l.logger.Println(message)
	return nil
}

// MultiNotifier sends notifications to multiple notifiers
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier creates a notifier that sends to multiple destinations
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
	}
}

// Notify sends notifications to all configured notifiers
func (m *MultiNotifier) Notify(secret vault.Secret) error {
	var errs []error

	for _, notifier := range m.notifiers {
		if err := notifier.Notify(secret); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send %d notifications: %v", len(errs), errs)
	}

	return nil
}

// ShouldNotify determines if a secret should trigger a notification
func ShouldNotify(secret vault.Secret, warningThreshold time.Duration) bool {
	return secret.IsExpired() || secret.IsExpiringSoon(warningThreshold)
}
