package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// snsPublisher abstracts the SNS publish call for testing.
type snsPublisher interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSNotifier sends Vault secret expiration alerts to an AWS SNS topic.
type SNSNotifier struct {
	client   snsPublisher
	topicARN string
}

// NewSNSNotifier creates an SNSNotifier that publishes to the given topic ARN.
// It loads AWS credentials from the default credential chain.
func NewSNSNotifier(topicARN string) (*SNSNotifier, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns: topic ARN must not be empty")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("sns: failed to load AWS config: %w", err)
	}

	return &SNSNotifier{
		client:   sns.NewFromConfig(cfg),
		topicARN: topicARN,
	}, nil
}

// newSNSNotifierWithClient creates an SNSNotifier with a custom publisher (for testing).
func newSNSNotifierWithClient(topicARN string, client snsPublisher) (*SNSNotifier, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("sns: topic ARN must not be empty")
	}
	return &SNSNotifier{client: client, topicARN: topicARN}, nil
}

// Notify publishes a secret expiration alert to the configured SNS topic.
func (n *SNSNotifier) Notify(secret *vault.Secret) error {
	msg := FormatMessage(secret)
	_, err := n.client.Publish(context.Background(), &sns.PublishInput{
		TopicArn: aws.String(n.topicARN),
		Subject:  aws.String(msg.Subject),
		Message:  aws.String(msg.Body),
	})
	if err != nil {
		return fmt.Errorf("sns: publish failed: %w", err)
	}
	return nil
}
