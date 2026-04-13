package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/wakeward/vaultwatch/internal/vault"
)

// sqsClient is the interface used to publish messages to SQS.
type sqsClient interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

// SQSNotifier sends secret expiration alerts to an AWS SQS queue.
type SQSNotifier struct {
	client   sqsClient
	queueURL string
}

// NewSQSNotifier creates a new SQSNotifier using the default AWS configuration.
func NewSQSNotifier(queueURL string) (*SQSNotifier, error) {
	if queueURL == "" {
		return nil, fmt.Errorf("sqs: queue URL is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sqs: failed to load AWS config: %w", err)
	}
	return newSQSNotifierWithClient(sqs.NewFromConfig(cfg), queueURL)
}

func newSQSNotifierWithClient(client sqsClient, queueURL string) (*SQSNotifier, error) {
	if queueURL == "" {
		return nil, fmt.Errorf("sqs: queue URL is required")
	}
	return &SQSNotifier{client: client, queueURL: queueURL}, nil
}

// Notify sends a secret expiration message to the configured SQS queue.
func (n *SQSNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg, err := json.Marshal(map[string]string{
		"path":    secret.Path,
		"subject": FormatMessage(secret).Subject,
		"body":    FormatMessage(secret).Body,
	})
	if err != nil {
		return fmt.Errorf("sqs: failed to marshal message: %w", err)
	}
	_, err = n.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(n.queueURL),
		MessageBody: aws.String(string(msg)),
	})
	if err != nil {
		return fmt.Errorf("sqs: failed to send message: %w", err)
	}
	return nil
}
