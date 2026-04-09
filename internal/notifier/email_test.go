package notifier

import (
	"net/smtp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmailNotifier_Valid(t *testing.T) {
	n, err := NewEmailNotifier(EmailConfig{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		From:     "alerts@example.com",
		To:       []string{"ops@example.com"},
		Username: "user",
		Password: "pass",
	})
	require.NoError(t, err)
	assert.NotNil(t, n)
}

func TestNewEmailNotifier_MissingHost(t *testing.T) {
	_, err := NewEmailNotifier(EmailConfig{
		SMTPPort: 587,
		From:     "alerts@example.com",
		To:       []string{"ops@example.com"},
	})
	assert.ErrorContains(t, err, "smtp host")
}

func TestNewEmailNotifier_MissingFrom(t *testing.T) {
	_, err := NewEmailNotifier(EmailConfig{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		To:       []string{"ops@example.com"},
	})
	assert.ErrorContains(t, err, "from address")
}

func TestNewEmailNotifier_MissingTo(t *testing.T) {
	_, err := NewEmailNotifier(EmailConfig{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		From:     "alerts@example.com",
	})
	assert.ErrorContains(t, err, "recipient")
}

func TestEmailNotifier_Notify_Expired(t *testing.T) {
	var capturedMsg []byte
	var capturedAuth smtp.Auth
	var capturedAddr string

	n, err := NewEmailNotifier(EmailConfig{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		From:     "alerts@example.com",
		To:       []string{"ops@example.com"},
	})
	require.NoError(t, err)

	n.sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		capturedAddr = addr
		capturedAuth = a
		capturedMsg = msg
		return nil
	}

	secret := testSecret("secret/db/password", time.Now().Add(-24*time.Hour))
	err = n.Notify(secret)
	require.NoError(t, err)

	assert.Equal(t, "smtp.example.com:587", capturedAddr)
	assert.Nil(t, capturedAuth)
	assert.Contains(t, string(capturedMsg), "EXPIRED")
	assert.Contains(t, string(capturedMsg), "secret/db/password")
}

func TestEmailNotifier_Notify_ExpiringSoon(t *testing.T) {
	n, err := NewEmailNotifier(EmailConfig{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		From:     "alerts@example.com",
		To:       []string{"ops@example.com"},
	})
	require.NoError(t, err)

	var capturedMsg []byte
	n.sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		capturedMsg = msg
		return nil
	}

	secret := testSecret("secret/api/key", time.Now().Add(5*24*time.Hour))
	err = n.Notify(secret)
	require.NoError(t, err)

	assert.Contains(t, string(capturedMsg), "expiring")
	assert.Contains(t, string(capturedMsg), "secret/api/key")
}

func testSecret(path string, expiry time.Time) interface{ Path() string; ExpiresAt() time.Time } {
	return &mockSecret{path: path, expiresAt: expiry}
}

type mockSecret struct {
	path      string
	expiresAt time.Time
}

func (m *mockSecret) Path() string        { return m.path }
func (m *mockSecret) ExpiresAt() time.Time { return m.expiresAt }
