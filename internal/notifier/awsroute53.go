package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// route53Client defines the subset of Route53 API used by the notifier.
type route53Client interface {
	ChangeResourceRecordSets(ctx context.Context, params *route53.ChangeResourceRecordSetsInput, optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error)
}

// Route53Notifier writes a TXT record to a Route53 hosted zone on secret expiry events.
type Route53Notifier struct {
	client  route53Client
	zoneID  string
	record  string
	ttl     int64
}

// NewRoute53Notifier creates a Route53Notifier using the default AWS credential chain.
func NewRoute53Notifier(zoneID, record string, ttl int64) (*Route53Notifier, error) {
	if zoneID == "" {
		return nil, fmt.Errorf("route53: hosted zone ID is required")
	}
	if record == "" {
		return nil, fmt.Errorf("route53: record name is required")
	}
	if ttl <= 0 {
		ttl = 300
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("route53: failed to load AWS config: %w", err)
	}
	return &Route53Notifier{
		client: route53.NewFromConfig(cfg),
		zoneID: zoneID,
		record: record,
		ttl:    ttl,
	}, nil
}

func newRoute53NotifierWithClient(client route53Client, zoneID, record string, ttl int64) *Route53Notifier {
	return &Route53Notifier{client: client, zoneID: zoneID, record: record, ttl: ttl}
}

// Notify upserts a TXT record in Route53 describing the secret expiration event.
func (n *Route53Notifier) Notify(ctx context.Context, s *vault.Secret) error {
	msg := FormatMessage(s)
	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(n.zoneID),
		ChangeBatch: &types.ChangeBatch{
			Comment: aws.String(fmt.Sprintf("vaultwatch alert at %s", time.Now().UTC().Format(time.RFC3339))),
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(n.record),
						Type: types.RRTypeTxt,
						TTL:  aws.Int64(n.ttl),
						ResourceRecords: []types.ResourceRecord{
							{Value: aws.String(fmt.Sprintf("%q", msg.Body))},
						},
					},
				},
			},
		},
	}
	_, err := n.client.ChangeResourceRecordSets(ctx, input)
	if err != nil {
		return fmt.Errorf("route53: failed to upsert record: %w", err)
	}
	return nil
}
