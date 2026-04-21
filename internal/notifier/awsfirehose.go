package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// firehoseClient abstracts the Firehose PutRecord call for testing.
type firehoseClient interface {
	PutRecord(ctx context.Context, params *firehose.PutRecordInput, optFns ...func(*firehose.Options)) (*firehose.PutRecordOutput, error)
}

// FirehoseNotifier sends secret expiration events to an AWS Kinesis Data Firehose delivery stream.
type FirehoseNotifier struct {
	client     firehoseClient
	deliveryStream string
}

// NewFirehoseNotifier creates a FirehoseNotifier using the default AWS credential chain.
func NewFirehoseNotifier(deliveryStream string) (*FirehoseNotifier, error) {
	if deliveryStream == "" {
		return nil, fmt.Errorf("firehose: delivery stream name is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("firehose: failed to load AWS config: %w", err)
	}
	return newFirehoseNotifierWithClient(firehose.NewFromConfig(cfg), deliveryStream)
}

func newFirehoseNotifierWithClient(client firehoseClient, deliveryStream string) (*FirehoseNotifier, error) {
	if deliveryStream == "" {
		return nil, fmt.Errorf("firehose: delivery stream name is required")
	}
	return &FirehoseNotifier{client: client, deliveryStream: deliveryStream}, nil
}

// Notify encodes the secret event as JSON and puts it onto the Firehose delivery stream.
func (n *FirehoseNotifier) Notify(ctx context.Context, secret vault.Secret) error {
	payload := map[string]interface{}{
		"path":        secret.Path,
		"expires_at":  secret.ExpiresAt.Format("2006-01-02T15:04:05Z"),
		"days_left":   secret.DaysUntilExpiration(),
		"is_expired":  secret.IsExpired(),
		"message":     FormatMessage(secret),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("firehose: failed to marshal payload: %w", err)
	}
	// Firehose records must end with a newline for line-delimited formats.
	data = append(data, '\n')
	_, err = n.client.PutRecord(ctx, &firehose.PutRecordInput{
		DeliveryStreamName: aws.String(n.deliveryStream),
		Record:             &types.Record{Data: data},
	})
	if err != nil {
		return fmt.Errorf("firehose: put record failed: %w", err)
	}
	return nil
}
