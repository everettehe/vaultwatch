package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// lambdaClient is the interface used to invoke Lambda functions.
type lambdaClient interface {
	Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
}

// LambdaNotifier invokes an AWS Lambda function with secret expiration payloads.
type LambdaNotifier struct {
	client       lambdaClient
	functionName string
}

type lambdaPayload struct {
	Path       string `json:"path"`
	DaysLeft   int    `json:"days_until_expiration"`
	Expired    bool   `json:"expired"`
	Message    string `json:"message"`
}

// NewLambdaNotifier creates a LambdaNotifier using the default AWS config.
func NewLambdaNotifier(functionName string) (*LambdaNotifier, error) {
	if functionName == "" {
		return nil, fmt.Errorf("lambda: function name is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("lambda: failed to load AWS config: %w", err)
	}
	return &LambdaNotifier{
		client:       lambda.NewFromConfig(cfg),
		functionName: functionName,
	}, nil
}

func newLambdaNotifierWithClient(client lambdaClient, functionName string) (*LambdaNotifier, error) {
	if functionName == "" {
		return nil, fmt.Errorf("lambda: function name is required")
	}
	return &LambdaNotifier{client: client, functionName: functionName}, nil
}

// Notify invokes the configured Lambda function with the secret's expiration details.
func (n *LambdaNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg, _ := FormatMessage(secret)
	payload := lambdaPayload{
		Path:     secret.Path,
		DaysLeft: secret.DaysUntilExpiration(),
		Expired:  secret.IsExpired(),
		Message:  msg,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("lambda: failed to marshal payload: %w", err)
	}
	_, err = n.client.Invoke(ctx, &lambda.InvokeInput{
		FunctionName: aws.String(n.functionName),
		Payload:      body,
	})
	if err != nil {
		return fmt.Errorf("lambda: invocation failed: %w", err)
	}
	return nil
}
