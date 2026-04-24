package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// iotPublisher abstracts the IoT Data Plane publish call.
type iotPublisher interface {
	Publish(ctx context.Context, params *iotdataplane.PublishInput, optFns ...func(*iotdataplane.Options)) (*iotdataplane.PublishOutput, error)
}

// AWSIoTNotifier publishes secret expiration alerts to an AWS IoT topic.
type AWSIoTNotifier struct {
	client iotPublisher
	topic  string
	region string
}

// NewAWSIoTNotifier creates an AWSIoTNotifier using the default AWS config.
func NewAWSIoTNotifier(topic, region string) (*AWSIoTNotifier, error) {
	if topic == "" {
		return nil, fmt.Errorf("awsiot: topic is required")
	}
	if region == "" {
		return nil, fmt.Errorf("awsiot: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("awsiot: failed to load AWS config: %w", err)
	}
	client := iotdataplane.NewFromConfig(cfg)
	return newAWSIoTNotifierWithClient(client, topic, region)
}

func newAWSIoTNotifierWithClient(client iotPublisher, topic, region string) (*AWSIoTNotifier, error) {
	if topic == "" {
		return nil, fmt.Errorf("awsiot: topic is required")
	}
	if region == "" {
		return nil, fmt.Errorf("awsiot: region is required")
	}
	return &AWSIoTNotifier{client: client, topic: topic, region: region}, nil
}

// Notify publishes a JSON payload to the configured IoT MQTT topic.
func (n *AWSIoTNotifier) Notify(ctx context.Context, secret vault.Secret) error {
	msg, err := json.Marshal(map[string]interface{}{
		"path":        secret.Path,
		"days_left":   secret.DaysUntilExpiration(),
		"is_expired":  secret.IsExpired(),
		"message":     FormatMessage(secret),
	})
	if err != nil {
		return fmt.Errorf("awsiot: failed to marshal payload: %w", err)
	}
	_, err = n.client.Publish(ctx, &iotdataplane.PublishInput{
		Topic:   aws.String(n.topic),
		Payload: msg,
	})
	if err != nil {
		return fmt.Errorf("awsiot: publish failed: %w", err)
	}
	return nil
}
