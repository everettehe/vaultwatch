package notifier

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
	"time"

	"vaultwatch/internal/vault"
)

func TestLogNotifier_Notify_ExpiringSoon(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	notifier := NewLogNotifier(logger)

	secret := vault.Secret{
		Path:           "secret/data/test",
		ExpirationTime: time.Now().Add(5 * 24 * time.Hour),
	}

	err := notifier.Notify(secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[WARNING]") {
		t.Errorf("expected [WARNING] in output, got: %s", output)
	}
	if !strings.Contains(output, "secret/data/test") {
		t.Errorf("expected secret path in output, got: %s", output)
	}
	if !strings.Contains(output, "expires in 5 days") {
		t.Errorf("expected 'expires in 5 days' in output, got: %s", output)
	}
}

func TestLogNotifier_Notify_Expired(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	notifier := NewLogNotifier(logger)

	secret := vault.Secret{
		Path:           "secret/data/expired",
		ExpirationTime: time.Now().Add(-3 * 24 * time.Hour),
	}

	err := notifier.Notify(secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[EXPIRED]") {
		t.Errorf("expected [EXPIRED] in output, got: %s", output)
	}
	if !strings.Contains(output, "expired 3 days ago") {
		t.Errorf("expected 'expired 3 days ago' in output, got: %s", output)
	}
}

func TestShouldNotify(t *testing.T) {
	tests := []struct {
		name      string
		secret    vault.Secret
		threshold time.Duration
		want      bool
	}{
		{
			name:      "expired secret",
			secret:    vault.Secret{ExpirationTime: time.Now().Add(-24 * time.Hour)},
			threshold: 7 * 24 * time.Hour,
			want:      true,
		},
		{
			name:      "expiring soon",
			secret:    vault.Secret{ExpirationTime: time.Now().Add(3 * 24 * time.Hour)},
			threshold: 7 * 24 * time.Hour,
			want:      true,
		},
		{
			name:      "not expiring soon",
			secret:    vault.Secret{ExpirationTime: time.Now().Add(30 * 24 * time.Hour)},
			threshold: 7 * 24 * time.Hour,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldNotify(tt.secret, tt.threshold)
			if got != tt.want {
				t.Errorf("ShouldNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockNotifier struct {
	called bool
	err    error
}

func (m *mockNotifier) Notify(secret vault.Secret) error {
	m.called = true
	return m.err
}

func TestMultiNotifier_Notify(t *testing.T) {
	mock1 := &mockNotifier{}
	mock2 := &mockNotifier{err: errors.New("test error")}

	multi := NewMultiNotifier(mock1, mock2)
	secret := vault.Secret{Path: "test"}

	err := multi.Notify(secret)

	if !mock1.called {
		t.Error("expected first notifier to be called")
	}
	if !mock2.called {
		t.Error("expected second notifier to be called")
	}
	if err == nil {
		t.Error("expected error from multi notifier")
	}
}
