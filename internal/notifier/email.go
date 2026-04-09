package notifier

import (
	"fmt"
	"net/smtp"
	"time"

	"vaultwatch/internal/vault"
)

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
	To       []string
}

// EmailNotifier sends notifications via email
type EmailNotifier struct {
	config EmailConfig
}

// NewEmailNotifier creates a new email notifier
func NewEmailNotifier(config EmailConfig) (*EmailNotifier, error) {
	if config.SMTPHost == "" {
		return nil, fmt.Errorf("SMTP host is required")
	}
	if config.From == "" {
		return nil, fmt.Errorf("from address is required")
	}
	if len(config.To) == 0 {
		return nil, fmt.Errorf("at least one recipient is required")
	}

	return &EmailNotifier{
		config: config,
	}, nil
}

// Notify sends an email notification
func (e *EmailNotifier) Notify(secret vault.Secret) error {
	subject := e.buildSubject(secret)
	body := e.buildBody(secret)

	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n",
		e.config.From,
		e.config.To[0],
		subject,
		body,
	)

	addr := fmt.Sprintf("%s:%d", e.config.SMTPHost, e.config.SMTPPort)
	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.SMTPHost)

	err := smtp.SendMail(addr, auth, e.config.From, e.config.To, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (e *EmailNotifier) buildSubject(secret vault.Secret) string {
	if secret.IsExpired() {
		return fmt.Sprintf("[VAULTWATCH] Secret Expired: %s", secret.Path)
	}
	return fmt.Sprintf("[VAULTWATCH] Secret Expiring Soon: %s", secret.Path)
}

func (e *EmailNotifier) buildBody(secret vault.Secret) string {
	days := secret.DaysUntilExpiration()

	if secret.IsExpired() {
		return fmt.Sprintf(
			"The following Vault secret has EXPIRED:\n\n"+
				"Path: %s\n"+
				"Expired: %d days ago\n"+
				"Expiration Date: %s\n\n"+
				"Please rotate this secret immediately.\n",
			secret.Path,
			-days,
			secret.ExpirationTime.Format(time.RFC3339),
		)
	}

	return fmt.Sprintf(
		"The following Vault secret is expiring soon:\n\n"+
			"Path: %s\n"+
			"Days Until Expiration: %d\n"+
			"Expiration Date: %s\n\n"+
			"Please plan to rotate this secret before it expires.\n",
		secret.Path,
		days,
		secret.ExpirationTime.Format(time.RFC3339),
	)
}
