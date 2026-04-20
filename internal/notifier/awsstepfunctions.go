package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

// sfnClient is the interface for AWS Step Functions operations used by StepFunctionsNotifier.
type sfnClient interface {
	StartExecution(ctx context.Context, params *sfn.StartExecutionInput, optFns ...func(*sfn.Options)) (*sfn.StartExecutionOutput, error)
}

// StepFunctionsNotifier triggers an AWS Step Functions state machine execution on secret expiry.
type StepFunctionsNotifier struct {
	client       sfnClient
	stateMachineARN string
}

// NewStepFunctionsNotifier creates a StepFunctionsNotifier using the default AWS credential chain.
func NewStepFunctionsNotifier(stateMachineARN string) (*StepFunctionsNotifier, error) {
	if stateMachineARN == "" {
		return nil, fmt.Errorf("step functions: state machine ARN is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("step functions: failed to load AWS config: %w", err)
	}
	return newStepFunctionsNotifierWithClient(sfn.NewFromConfig(cfg), stateMachineARN)
}

func newStepFunctionsNotifierWithClient(client sfnClient, stateMachineARN string) (*StepFunctionsNotifier, error) {
	if stateMachineARN == "" {
		return nil, fmt.Errorf("step functions: state machine ARN is required")
	}
	return &StepFunctionsNotifier{client: client, stateMachineARN: stateMachineARN}, nil
}

// Notify starts a Step Functions execution with the secret expiry details as input.
func (n *StepFunctionsNotifier) Notify(ctx context.Context, secret ExpiringSecret) error {
	payload := map[string]interface{}{
		"path":        secret.Path,
		"days_left":   secret.DaysUntilExpiration,
		"expires_at":  secret.ExpiresAt.Format(time.RFC3339),
		"is_expired":  secret.IsExpired,
		"message":     FormatMessage(secret),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("step functions: failed to marshal input: %w", err)
	}
	name := fmt.Sprintf("vaultwatch-%d", time.Now().UnixNano())
	_, err = n.client.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(n.stateMachineARN),
		Input:           aws.String(string(data)),
		Name:            aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("step functions: failed to start execution: %w", err)
	}
	return nil
}
