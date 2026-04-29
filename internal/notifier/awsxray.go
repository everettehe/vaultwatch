package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	"github.com/aws/aws-sdk-go-v2/service/xray/types"
	"github.com/your-org/vaultwatch/internal/vault"
)

type xrayClient interface {
	PutTraceSegments(ctx context.Context, params *xray.PutTraceSegmentsInput, optFns ...func(*xray.Options)) (*xray.PutTraceSegmentsOutput, error)
}

// XRayNotifier sends vault secret expiration events as X-Ray trace segments.
type XRayNotifier struct {
	client  xrayClient
	service string
	region  string
}

// NewXRayNotifier creates an XRayNotifier using the provided region and service name.
func NewXRayNotifier(region, service string) (*XRayNotifier, error) {
	if region == "" {
		return nil, fmt.Errorf("xray: region is required")
	}
	if service == "" {
		service = "vaultwatch"
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("xray: failed to load AWS config: %w", err)
	}
	return &XRayNotifier{
		client:  xray.NewFromConfig(cfg),
		service: service,
		region:  region,
	}, nil
}

func newXRayNotifierWithClient(client xrayClient, region, service string) *XRayNotifier {
	return &XRayNotifier{client: client, region: region, service: service}
}

// Notify sends a trace segment to AWS X-Ray describing the secret expiration event.
func (n *XRayNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	now := float64(time.Now().UnixNano()) / 1e9
	msg := FormatMessage(secret)
	segment := fmt.Sprintf(
		`{"name":%q,"id":"vaultwatch-xray","start_time":%f,"end_time":%f,"annotations":{"secret_path":%q,"message":%q}}`,
		n.service, now, now, secret.Path, msg.Body,
	)
	_, err := n.client.PutTraceSegments(ctx, &xray.PutTraceSegmentsInput{
		TraceSegmentDocuments: []string{segment},
	})
	if err != nil {
		return fmt.Errorf("xray: failed to put trace segment: %w", err)
	}
	_ = types.EncryptionConfig{} // ensure types import is used
	return nil
}
