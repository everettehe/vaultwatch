package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// snsFIFOPublisher abstracts the SNS Publish call for FIFO topics.
type snsFIFOPublisher interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSFIFONotifier sends notifications to an AWS SNS FIFO topic.
type SNSFIFONotifier struct {
	client   snsFIFOPublisher
	topicARN string
	groupID  string
}

// NewSNSFIFONotifier creates a new SNSFIFONotifier using the default AWS config.
func NewSNSFIFONotifier(topicARN, groupID string) (*SNSFIFONotifier, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns fifo: topic ARN is required")
	}
	if groupID == "" {
		return nil, fmt.Errorf("sns fifo: message group ID is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sns fifo: failed to load AWS config: %w", err)
	}
	return newSNSFIFONotifierWithClient(sns.NewFromConfig(cfg), topicARN, groupID)
}

func newSNSFIFONotifierWithClient(client snsFIFOPublisher, topicARN, groupID string) (*SNSFIFONotifier, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns fifo: topic ARN is required")
	}
	if groupID == "" {
		return nil, fmt.Errorf("sns fifo: message group ID is required")
	}
	return &SNSFIFONotifier{client: client, topicARN: topicARN, groupID: groupID}, nil
}

// Notify publishes a secret expiration alert to the FIFO SNS topic.
func (n *SNSFIFONotifier) Notify(ctx context.Context, secret vault.Secret) error {
	msg, err := json.Marshal(map[string]interface{}{
		"path":    secret.Path,
		"expires": secret.ExpiresAt,
		"days":    secret.DaysUntilExpiration(),
	})
	if err != nil {
		return fmt.Errorf("sns fifo: failed to marshal message: %w", err)
	}
	_, err = n.client.Publish(ctx, &sns.PublishInput{
		TopicArn:       aws.String(n.topicARN),
		Message:        aws.String(string(msg)),
		MessageGroupId: aws.String(n.groupID),
	})
	if err != nil {
		return fmt.Errorf("sns fifo: publish failed: %w", err)
	}
	return nil
}
