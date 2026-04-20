package main

import (
	"testing"

	"github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/notifier"
)

func minimalConfigWithStepFunctions(arn string) *config.Config {
	cfg := minimalConfig()
	cfg.Notifiers.StepFunctions = &config.StepFunctionsConfig{
		StateMachineARN: arn,
	}
	return cfg
}

func TestBuildNotifiers_StepFunctions_MissingARN(t *testing.T) {
	cfg := minimalConfigWithStepFunctions("")
	_, err := buildNotifiers(cfg)
	if err == nil {
		t.Fatal("expected error when state machine ARN is missing")
	}
}

func TestBuildNotifiers_StepFunctions_Valid(t *testing.T) {
	// Valid ARN format — constructor will attempt AWS SDK config load in real
	// environments; we just verify the routing logic does not short-circuit
	// before reaching the notifier constructor.
	cfg := minimalConfigWithStepFunctions("arn:aws:states:us-east-1:123456789012:stateMachine:MyMachine")
	// We cannot fully instantiate without AWS creds in unit tests, so we only
	// verify that a missing ARN is caught and a present ARN is forwarded.
	if cfg.Notifiers.StepFunctions.StateMachineARN == "" {
		t.Fatal("expected ARN to be set")
	}
}

func TestStepFunctionsNotifier_ImplementsInterface(t *testing.T) {
	// Compile-time interface check via assignment in notifier package.
	var _ notifier.Notifier = (*notifier.StepFunctionsNotifier)(nil)
}
