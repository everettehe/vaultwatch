package runner_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/config"
	"github.com/your-org/vaultwatch/internal/monitor"
	"github.com/your-org/vaultwatch/internal/runner"
	"github.com/your-org/vaultwatch/internal/scanner"
)

func minimalConfig() *config.Config {
	return &config.Config{
		ScanInterval: 10 * time.Millisecond,
	}
}

// stubLister satisfies vault.Lister for scanner construction.
type stubLister struct {
	paths []string
	err   error
}

func (s *stubLister) List(_ context.Context, _ string) ([]string, error) {
	return s.paths, s.err
}

func TestNew_NilConfig(t *testing.T) {
	_, err := runner.New(nil, &scanner.Scanner{}, &monitor.Monitor{})
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestNew_NilScanner(t *testing.T) {
	_, err := runner.New(minimalConfig(), nil, &monitor.Monitor{})
	if err == nil {
		t.Fatal("expected error for nil scanner")
	}
}

func TestNew_NilMonitor(t *testing.T) {
	_, err := runner.New(minimalConfig(), &scanner.Scanner{}, nil)
	if err == nil {
		t.Fatal("expected error for nil monitor")
	}
}

func TestNew_Valid(t *testing.T) {
	cfg := minimalConfig()
	s := &scanner.Scanner{}
	m := &monitor.Monitor{}

	r, err := runner.New(cfg, s, m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestStart_CancelledContext(t *testing.T) {
	cfg := minimalConfig()
	s := &scanner.Scanner{}
	m := &monitor.Monitor{}

	r, err := runner.New(cfg, s, m)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err = r.Start(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRunOnce_ScanError(t *testing.T) {
	cfg := minimalConfig()
	// scanner with a lister that errors
	lister := &stubLister{err: errors.New("vault unavailable")}
	s, _ := scanner.New(lister, []string{"secret/"})
	m := &monitor.Monitor{}

	r, err := runner.New(cfg, s, m)
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	err = r.RunOnce(context.Background())
	if err == nil {
		t.Fatal("expected error from RunOnce when scanner fails")
	}
}
