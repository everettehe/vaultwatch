// Package main is the entry point for the vaultwatch CLI.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/monitor"
	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/runner"
	"github.com/yourusername/vaultwatch/internal/scanner"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfgPath := "configs/vaultwatch.yaml"
	for i, a := range args {
		if (a == "--config" || a == "-c") && i+1 < len(args) {
			cfgPath = args[i+1]
		}
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	vaultClient, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	lister := vault.NewLister(vaultClient)
	scn := scanner.New(lister)

	notifiers, err := buildNotifiers(cfg)
	if err != nil {
		return fmt.Errorf("building notifiers: %w", err)
	}

	mon, err := monitor.New(vaultClient, notifier.NewMultiNotifier(notifiers...), cfg)
	if err != nil {
		return fmt.Errorf("creating monitor: %w", err)
	}

	r, err := runner.New(cfg, scn, mon)
	if err != nil {
		return fmt.Errorf("creating runner: %w", err)
	}

	return r.Run()
}

func buildNotifiers(cfg *config.Config) ([]notifier.Notifier, error) {
	var notifiers []notifier.Notifier

	notifiers = append(notifiers, notifier.NewLogNotifier())

	if cfg.Notifiers.Slack.WebhookURL != "" {
		sn, err := notifier.NewSlackNotifier(cfg.Notifiers.Slack.WebhookURL)
		if err != nil {
			return nil, fmt.Errorf("slack notifier: %w", err)
		}
		notifiers = append(notifiers, sn)
	}

	if cfg.Notifiers.SNS.TopicARN != "" {
		snsN, err := notifier.NewSNSNotifier(cfg.Notifiers.SNS.TopicARN)
		if err != nil {
			return nil, fmt.Errorf("sns notifier: %w", err)
		}
		notifiers = append(notifiers, snsN)
	}

	if len(notifiers) == 0 {
		return nil, errors.New("no notifiers configured")
	}
	return notifiers, nil
}
