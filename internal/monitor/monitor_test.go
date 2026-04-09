package monitor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// mockFetcher implements SecretFetcher for tests.
type mockFetcher struct {
	secret *vault.Secret
	err    error
}

func (m *mockFetcher) GetSecretMetadata(_ context.Context, _ string) (*vault.Secret, error) {
	return m.secret, m.err
}

// mockNotifier records Notify calls.
type mockNotifier struct {
	called int
	err    error
}

func (m *mockNotifier) Notify(_ context.Context, _ *vault.Secret) error {
	m.called++
	return m.err
}

func expiringSecret(daysUntil int) *vault.Secret {
	expiry := time.Now().Add(time.Duration(daysUntil) * 24 * time.Hour)
	return &vault.Secret{Path: "secret/test", ExpiresAt: &expiry}
}

func TestNew_MissingClient(t *testing.T) {
	_, err := New(nil, &mockNotifier{}, Config{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestNew_MissingNotifier(t *testing.T) {
	_, err := New(&mockFetcher{}, nil, Config{Paths: []string{"secret/a"}})
	if err == nil {
		t.Fatal("expected error for nil notifier")
	}
}

func TestNew_MissingPaths(t *testing.T) {
	_, err := New(&mockFetcher{}, &mockNotifier{}, Config{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestNew_Defaults(t *testing.T) {
	m, err := New(&mockFetcher{}, &mockNotifier{}, Config{Paths: []string{"secret/a"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.interval != time.Hour {
		t.Errorf("expected default interval 1h, got %s", m.interval)
	}
	if m.warnDays != 7 {
		t.Errorf("expected default warn_days 7, got %d", m.warnDays)
	}
}

func TestCheckAll_NotifiesExpiringSoon(t *testing.T) {
	n := &mockNotifier{}
	f := &mockFetcher{secret: expiringSecret(3)}
	m, _ := New(f, n, Config{Paths: []string{"secret/a"}, WarnDays: 7})

	m.checkAll(context.Background())

	if n.called != 1 {
		t.Errorf("expected 1 notification, got %d", n.called)
	}
}

func TestCheckAll_SkipsHealthySecret(t *testing.T) {
	n := &mockNotifier{}
	f := &mockFetcher{secret: expiringSecret(30)}
	m, _ := New(f, n, Config{Paths: []string{"secret/a"}, WarnDays: 7})

	m.checkAll(context.Background())

	if n.called != 0 {
		t.Errorf("expected 0 notifications, got %d", n.called)
	}
}

func TestCheckAll_FetchError(t *testing.T) {
	n := &mockNotifier{}
	f := &mockFetcher{err: errors.New("vault unavailable")}
	m, _ := New(f, n, Config{Paths: []string{"secret/a"}, WarnDays: 7})

	// Should not panic; error is logged and skipped.
	m.checkAll(context.Background())

	if n.called != 0 {
		t.Errorf("expected 0 notifications on fetch error, got %d", n.called)
	}
}
