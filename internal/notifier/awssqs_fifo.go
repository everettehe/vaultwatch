package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	vaultsecret "github.com/yourusername/vaultwatch/internal/vault"
)

// sqsFIFOClient is the interface for sending messages to SQS FIFO queues.
type sqsFIFOClient interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

// SQSFIFONotifier sends vault secret expiration alerts to an AWS SQS FIFO queue.
type SQSFIFONotifier struct {
	client       sqsFIFOClient
	queueURL     string
	messageGroup string
}

// NewSQSFIFONotifier creates a new SQSFIFONotifier using the default AWS configuration.
func NewSQSFIFONotifier(queueURL, messageGroup string) (*SQSFIFONotifier, error) {
	if queueURL == "" {
		return nil, fmt.Errorf("SQS FIFO queue URL is required")
	}
	if messageGroup == "" {
		return nil, fmt.Errorf("SQS FIFO message group ID is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	return newSQSFIFONotifierWithClient(sqs.NewFromConfig(cfg), queueURL, messageGroup)
}

func newSQSFIFONotifierWithClient(client sqsFIFOClient, queueURL, messageGroup string) (*SQSFIFONotifier, error) {
	if queueURL == "" {
		return nil, fmt.Errorf("SQS FIFO queue URL is required")
	}
	if messageGroup == "" {
		return nil, fmt.Errorf("SQS FIFO message group ID is required")
	}
	return &SQSFIFONotifier{
		client:       client,
		queueURL:     queueURL,
		messageGroup: messageGroup,
	}, nil
}

// Notify sends a secret expiration alert to the configured SQS FIFO queue.
func (n *SQSFIFONotifier) Notify(ctx context.Context, secret *vaultsecret.Secret) error {
	msg, err := json.Marshal(map[string]interface{}{
		"path":        secret.Path,
		"days_left":   secret.DaysUntilExpiration(),
		"expires_at":  secret.ExpiresAt,
		"is_expired":  secret.IsExpired(),
		"message":     FormatMessage(secret),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal SQS FIFO message: %w", err)
	}
	_, err = n.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:               aws.String(n.queueURL),
		MessageBody:            aws.String(string(msg)),
		MessageGroupId:         aws.String(n.messageGroup),
		MessageDeduplicationId: aws.String(fmt.Sprintf("%s-%d", secret.Path, secret.ExpiresAt.Unix())),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS FIFO queue: %w", err)
	}
	return nil
}
