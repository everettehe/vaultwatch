// Package runner ties together scanning, monitoring, and notification
// into a single periodic execution loop.
package runner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/your-org/vaultwatch/internal/config"
	"github.com/your-org/vaultwatch/internal/monitor"
	"github.com/your-org/vaultwatch/internal/scanner"
)

// Runner orchestrates a periodic scan-and-alert cycle.
type Runner struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	monitor *monitor.Monitor
	interval time.Duration
}

// New creates a Runner from the provided dependencies.
func New(cfg *config.Config, s *scanner.Scanner, m *monitor.Monitor) (*Runner, error) {
	if cfg == nil {
		return nil, fmt.Errorf("runner: config must not be nil")
	}
	if s == nil {
		return nil, fmt.Errorf("runner: scanner must not be nil")
	}
	if m == nil {
		return nil, fmt.Errorf("runner: monitor must not be nil")
	}

	interval := cfg.ScanInterval
	if interval <= 0 {
		interval = 5 * time.Minute
	}

	return &Runner{
		cfg:      cfg,
		scanner:  s,
		monitor:  m,
		interval: interval,
	}, nil
}

// RunOnce performs a single scan and processes each discovered secret.
func (r *Runner) RunOnce(ctx context.Context) error {
	secrets, err := r.scanner.Scan(ctx)
	if err != nil {
		return fmt.Errorf("runner: scan failed: %w", err)
	}

	log.Printf("runner: discovered %d secrets", len(secrets))

	for _, s := range secrets {
		if err := r.monitor.Check(ctx, s); err != nil {
			log.Printf("runner: check failed for %s: %v", s.Path, err)
		}
	}
	return nil
}

// Start runs the scan loop until ctx is cancelled.
func (r *Runner) Start(ctx context.Context) error {
	log.Printf("runner: starting with interval %s", r.interval)

	if err := r.RunOnce(ctx); err != nil {
		log.Printf("runner: initial run error: %v", err)
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("runner: shutting down")
			return ctx.Err()
		case <-ticker.C:
			if err := r.RunOnce(ctx); err != nil {
				log.Printf("runner: run error: %v", err)
			}
		}
	}
}
