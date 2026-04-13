package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"

	vaultsecret "github.com/yourusername/vaultwatch/internal/vault"
)

type kinesisClient interface {
	PutRecord(ctx context.Context, params *kinesis.PutRecordInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error)
}

// KinesisNotifier sends secret expiration events to an AWS Kinesis Data Stream.
type KinesisNotifier struct {
	client     kinesisClient
	streamName string
	partitionKey string
}

type kinesisPayload struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Path    string `json:"path"`
	Status  string `json:"status"`
}

// NewKinesisNotifier creates a KinesisNotifier using the default AWS config.
func NewKinesisNotifier(streamName, partitionKey string) (*KinesisNotifier, error) {
	if streamName == "" {
		return nil, fmt.Errorf("kinesis: stream name is required")
	}
	if partitionKey == "" {
		partitionKey = "vaultwatch"
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("kinesis: failed to load AWS config: %w", err)
	}
	return newKinesisNotifierWithClient(kinesis.NewFromConfig(cfg), streamName, partitionKey), nil
}

func newKinesisNotifierWithClient(client kinesisClient, streamName, partitionKey string) *KinesisNotifier {
	return &KinesisNotifier{
		client:       client,
		streamName:   streamName,
		partitionKey: partitionKey,
	}
}

// Notify sends a secret expiration event to the configured Kinesis stream.
func (n *KinesisNotifier) Notify(ctx context.Context, secret *vaultsecret.Secret) error {
	msg := FormatMessage(secret)
	payload := kinesisPayload{
		Subject: msg.Subject,
		Body:    msg.Body,
		Path:    secret.Path,
		Status:  msg.Status,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("kinesis: failed to marshal payload: %w", err)
	}
	_, err = n.client.PutRecord(ctx, &kinesis.PutRecordInput{
		StreamName:   aws.String(n.streamName),
		PartitionKey: aws.String(n.partitionKey),
		Data:         data,
	})
	if err != nil {
		return fmt.Errorf("kinesis: failed to put record: %w", err)
	}
	return nil
}
