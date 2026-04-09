package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/monitor"
	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/runner"
	"github.com/yourusername/vaultwatch/internal/scanner"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfgPath := os.Getenv("VAULTWATCH_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/vaultwatch.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	vaultClient, err := vault.NewClient(cfg.Vault)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	lister := vault.NewLister(vaultClient)
	scn := scanner.New(lister, cfg.Scanner.Paths)

	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		return fmt.Errorf("building notifiers: %w", err)
	}

	mon := monitor.New(vaultClient, notifiers, cfg.Monitor)

	r, err := runner.New(cfg, scn, mon)
	if err != nil {
		return fmt.Errorf("creating runner: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	return r.Run(ctx)
}

func buildNotifiers(cfg *config.Config) (notifier.Notifier, error) {
	var notifiers []notifier.Notifier

	notifiers = append(notifiers, notifier.NewLogNotifier())

	if cfg.Notifiers.Slack.WebhookURL != "" {
		sn, err := notifier.NewSlackNotifier(cfg.Notifiers.Slack)
		if err != nil {
			return nil, fmt.Errorf("slack notifier: %w", err)
		}
		notifiers = append(notifiers, sn)
	}

	if cfg.Notifiers.Email.Host != "" {
		en, err := notifier.NewEmailNotifier(cfg.Notifiers.Email)
		if err != nil {
			return nil, fmt.Errorf("email notifier: %w", err)
		}
		notifiers = append(notifiers, en)
	}

	return notifier.NewMultiNotifier(notifiers...), nil
}
