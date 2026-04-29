package notifier

import (
	"fmt"
	"net/smtp"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SMTPNotifier sends notifications via SMTP (e.g. AWS SES SMTP endpoint).
type SMTPNotifier struct {
	host     string
	port     int
	username string
	password string
	from     string
	to       []string
}

// NewSMTPNotifier creates a new SMTPNotifier.
func NewSMTPNotifier(host string, port int, username, password, from string, to []string) (*SMTPNotifier, error) {
	if host == "" {
		return nil, fmt.Errorf("smtp: host is required")
	}
	if from == "" {
		return nil, fmt.Errorf("smtp: from address is required")
	}
	if len(to) == 0 {
		return nil, fmt.Errorf("smtp: at least one recipient is required")
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

// Notify sends an email via SMTP for the given secret.
func (n *SMTPNotifier) Notify(s *vault.Secret) error {
	msg := FormatMessage(s)
	body := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		n.from, n.to[0], msg.Subject, msg.Body)
	addr := fmt.Sprintf("%s:%d", n.host, n.port)
	var auth smtp.Auth
	if n.username != "" {
		auth = smtp.PlainAuth("", n.username, n.password, n.host)
	}
	return smtp.SendMail(addr, auth, n.from, n.to, []byte(body))
}
