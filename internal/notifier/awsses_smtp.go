package notifier

import (
	"fmt"
	"net/smtp"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SMTPNotifier sends notifications via SMTP (e.g. AWS SES SMTP interface).
type SMTPNotifier struct {
	host string
	port string
	username string
	password string
	from string
	to string
}

// NewSMTPNotifier creates a new SMTPNotifier.
func NewSMTPNotifier(host, port, username, password, from, to string) (*SMTPNotifier, error) {
	if host == "" {
		return nil, fmt.Errorf("smtp: host is required")
	}
	if from == "" {
		return nil, fmt.Errorf("smtp: from address is required")
	}
	if to == "" {
		return nil, fmt.Errorf("smtp: to address is required")
	}
	if port == "" {
		port = "587"
	}
	return &SMTPNotifier{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		to:       to,
	}, nil
}

// Notify sends an email notification via SMTP.
func (n *SMTPNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	body := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		n.from, n.to, msg.Subject, msg.Body)

	addr := fmt.Sprintf("%s:%s", n.host, n.port)
	var auth smtp.Auth
	if n.username != "" && n.password != "" {
		auth = smtp.PlainAuth("", n.username, n.password, n.host)
	}
	return smtp.SendMail(addr, auth, n.from, []string{n.to}, []byte(body))
}
