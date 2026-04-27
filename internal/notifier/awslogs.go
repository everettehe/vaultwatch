package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	vaultconfig "github.com/yourusername/vaultwatch/internal/config"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// CloudWatchLogsV2Notifier sends log events to AWS CloudWatch Logs.
type CloudWatchLogsV2Notifier struct {
	client    *cloudwatchlogs.Client
	logGroup  string
	logStream string
}

// NewCloudWatchLogsV2Notifier creates a CloudWatchLogsV2Notifier from config.
func NewCloudWatchLogsV2Notifier(cfg vaultconfig.CloudWatchLogsConfig) (*CloudWatchLogsV2Notifier, error) {
	if cfg.LogGroup == "" {
		return nil, fmt.Errorf("cloudwatchlogs: log group is required")
	}
	logStream := cfg.LogStream
	if logStream == "" {
		logStream = "vaultwatch"
	}
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("cloudwatchlogs: failed to load AWS config: %w", err)
	}
	return &CloudWatchLogsV2Notifier{
		client:    cloudwatchlogs.NewFromConfig(awsCfg),
		logGroup:  cfg.LogGroup,
		logStream: logStream,
	}, nil
}

// Notify sends a CloudWatch Logs event for the given secret.
func (n *CloudWatchLogsV2Notifier) Notify(ctx context.Context, s vault.Secret) error {
	msg := FormatMessage(s)
	_, err := n.client.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(n.logGroup),
		LogStreamName: aws.String(n.logStream),
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(msg.Body),
				Timestamp: aws.Int64(time.Now().UnixMilli()),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("cloudwatchlogs: failed to put log events: %w", err)
	}
	return nil
}
