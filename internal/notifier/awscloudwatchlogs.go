package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type cloudWatchLogsClient interface {
	PutLogEvents(ctx context.Context, params *cloudwatchlogs.PutLogEventsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.PutLogEventsOutput, error)
	CreateLogStream(ctx context.Context, params *cloudwatchlogs.CreateLogStreamInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogStreamOutput, error)
}

// CloudWatchLogsNotifier sends vault secret expiration events to AWS CloudWatch Logs.
type CloudWatchLogsNotifier struct {
	client     cloudWatchLogsClient
	logGroup   string
	logStream  string
	region     string
}

// NewCloudWatchLogsNotifier creates a new CloudWatchLogsNotifier.
func NewCloudWatchLogsNotifier(logGroup, logStream, region string) (*CloudWatchLogsNotifier, error) {
	if logGroup == "" {
		return nil, fmt.Errorf("cloudwatchlogs: log group is required")
	}
	if logStream == "" {
		logStream = "vaultwatch"
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs: failed to load AWS config: %w", err)
	}
	return newCloudWatchLogsNotifierWithClient(cloudwatchlogs.NewFromConfig(cfg), logGroup, logStream, region), nil
}

func newCloudWatchLogsNotifierWithClient(client cloudWatchLogsClient, logGroup, logStream, region string) *CloudWatchLogsNotifier {
	return &CloudWatchLogsNotifier{
		client:    client,
		logGroup:  logGroup,
		logStream: logStream,
		region:    region,
	}
}

// Notify sends a log event to CloudWatch Logs for the given secret.
func (n *CloudWatchLogsNotifier) Notify(ctx context.Context, secret vault.Secret) error {
	_, _ = n.client.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(n.logGroup),
		LogStreamName: aws.String(n.logStream),
	})

	msg := FormatMessage(secret)
	_, err := n.client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(n.logGroup),
		LogStreamName: aws.String(n.logStream),
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(msg),
				Timestamp: aws.Int64(time.Now().UnixMilli()),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("cloudwatchlogs: failed to put log events: %w", err)
	}
	return nil
}
