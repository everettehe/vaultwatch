package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// SecretFetcher defines the interface for fetching secrets from Vault.
type SecretFetcher interface {
	GetSecretMetadata(ctx context.Context, path string) (*vault.Secret, error)
}

// Monitor periodically checks Vault secrets and triggers notifications.
type Monitor struct {
	client   SecretFetcher
	notifier notifier.Notifier
	paths    []string
	interval time.Duration
	warnDays int
}

// Config holds configuration for the Monitor.
type Config struct {
	Paths    []string
	Interval time.Duration
	WarnDays int
}

// New creates a new Monitor instance.
func New(client SecretFetcher, n notifier.Notifier, cfg Config) (*Monitor, error) {
	if client == nil {
		return nil, fmt.Errorf("vault client is required")
	}
	if n == nil {
		return nil, fmt.Errorf("notifier is required")
	}
	if len(cfg.Paths) == 0 {
		return nil, fmt.Errorf("at least one secret path is required")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 1 * time.Hour
	}
	if cfg.WarnDays <= 0 {
		cfg.WarnDays = 7
	}
	return &Monitor{
		client:   client,
		notifier: n,
		paths:    cfg.Paths,
		interval: cfg.Interval,
		warnDays: cfg.WarnDays,
	}, nil
}

// Run starts the monitoring loop, blocking until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) error {
	log.Printf("monitor: starting, interval=%s warn_days=%d paths=%v", m.interval, m.warnDays, m.paths)
	m.checkAll(ctx)

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkAll(ctx)
		case <-ctx.Done():
			log.Println("monitor: shutting down")
			return ctx.Err()
		}
	}
}

// checkAll iterates over all configured paths and notifies where appropriate.
func (m *Monitor) checkAll(ctx context.Context) {
	for _, path := range m.paths {
		secret, err := m.client.GetSecretMetadata(ctx, path)
		if err != nil {
			log.Printf("monitor: failed to fetch secret %q: %v", path, err)
			continue
		}
		if notifier.ShouldNotify(secret, m.warnDays) {
			if err := m.notifier.Notify(ctx, secret); err != nil {
				log.Printf("monitor: failed to notify for secret %q: %v", path, err)
			}
		}
	}
}
