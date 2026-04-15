package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type cloudWatchClient interface {
	PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error)
}

// CloudWatchNotifier publishes a custom metric to AWS CloudWatch when a secret
// is expiring or has expired.
type CloudWatchNotifier struct {
	client    cloudWatchClient
	namespace string
}

// NewCloudWatchNotifier creates a CloudWatchNotifier using the default AWS
// credential chain. namespace is the CloudWatch metric namespace to publish to.
func NewCloudWatchNotifier(namespace string) (*CloudWatchNotifier, error) {
	if namespace == "" {
		namespace = "VaultWatch"
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cloudwatch: load aws config: %w", err)
	}
	return newCloudWatchNotifierWithClient(cloudwatch.NewFromConfig(cfg), namespace), nil
}

func newCloudWatchNotifierWithClient(client cloudWatchClient, namespace string) *CloudWatchNotifier {
	return &CloudWatchNotifier{client: client, namespace: namespace}
}

// Notify publishes a DaysUntilExpiration metric to CloudWatch for the given secret.
func (n *CloudWatchNotifier) Notify(s *vault.Secret) error {
	days := s.DaysUntilExpiration()
	_, err := n.client.PutMetricData(context.Background(), &cloudwatch.PutMetricDataInput{
		Namespace: aws.String(n.namespace),
		MetricData: []types.MetricDatum{
			{
				MetricName: aws.String("DaysUntilExpiration"),
				Timestamp:  aws.Time(time.Now()),
				Value:      aws.Float64(float64(days)),
				Unit:       types.StandardUnitCount,
				Dimensions: []types.Dimension{
					{Name: aws.String("SecretPath"), Value: aws.String(s.Path)},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("cloudwatch: put metric data: %w", err)
	}
	return nil
}
