package notifier

import (
	"fmt"
	"log/syslog"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// SyslogNotifier sends alerts to the system syslog daemon.
type SyslogNotifier struct {
	writer   *syslog.Writer
	tag      string
	facility syslog.Priority
}

// NewSyslogNotifier creates a SyslogNotifier using the given tag and syslog facility.
// facility should be one of syslog.LOG_LOCAL0–LOG_LOCAL7 or syslog.LOG_DAEMON, etc.
func NewSyslogNotifier(tag string, facility syslog.Priority) (*SyslogNotifier, error) {
	if tag == "" {
		tag = "vaultwatch"
	}

	w, err := syslog.New(facility|syslog.LOG_WARNING, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog notifier: failed to connect to syslog: %w", err)
	}

	return &SyslogNotifier{
		writer:   w,
		tag:      tag,
		facility: facility,
	}, nil
}

// Notify writes a syslog message for the given secret.
// Expired secrets are logged at LOG_ERR; expiring secrets at LOG_WARNING.
func (s *SyslogNotifier) Notify(secret vault.Secret) error {
	msg := FormatMessage(secret)

	var err error
	if secret.IsExpired() {
		err = s.writer.Err(msg.Body)
	} else {
		err = s.writer.Warning(msg.Body)
	}

	if err != nil {
		return fmt.Errorf("syslog notifier: failed to write message: %w", err)
	}
	return nil
}

// Close releases the underlying syslog connection.
func (s *SyslogNotifier) Close() error {
	return s.writer.Close()
}
