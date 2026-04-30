package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

type wafv2Client interface {
	UpdateIPSet(ctx context.Context, params *wafv2.UpdateIPSetInput, optFns ...func(*wafv2.Options)) (*wafv2.UpdateIPSetOutput, error)
}

// WAFNotifier tags an AWS WAF IP set with secret expiration metadata.
type WAFNotifier struct {
	client    wafv2Client
	ipSetID   string
	ipSetName string
	scope     types.Scope
	region    string
}

// NewWAFNotifier creates a WAFNotifier using the default AWS credential chain.
func NewWAFNotifier(ipSetID, ipSetName, scope, region string) (*WAFNotifier, error) {
	if ipSetID == "" {
		return nil, fmt.Errorf("waf: ip_set_id is required")
	}
	if region == "" {
		return nil, fmt.Errorf("waf: region is required")
	}
	s := types.ScopeRegional
	if scope == "CLOUDFRONT" {
		s = types.ScopeCloudfront
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("waf: failed to load AWS config: %w", err)
	}
	return newWAFNotifierWithClient(wafv2.NewFromConfig(cfg), ipSetID, ipSetName, s, region), nil
}

func newWAFNotifierWithClient(client wafv2Client, ipSetID, ipSetName string, scope types.Scope, region string) *WAFNotifier {
	return &WAFNotifier{
		client:    client,
		ipSetID:   ipSetID,
		ipSetName: ipSetName,
		scope:     scope,
		region:    region,
	}
}

// Notify updates the WAF IP set description with expiration info for the secret.
func (n *WAFNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	meta := map[string]string{
		"path":        secret.Path,
		"days_left":   fmt.Sprintf("%d", int(secret.DaysUntilExpiration())),
		"status":      severityLabel(secret),
	}
	payload, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("waf: failed to marshal metadata: %w", err)
	}
	_, err = n.client.UpdateIPSet(ctx, &wafv2.UpdateIPSetInput{
		Id:          aws.String(n.ipSetID),
		Name:        aws.String(n.ipSetName),
		Scope:       n.scope,
		Description: aws.String(string(payload)),
		Addresses:   []string{},
	})
	if err != nil {
		return fmt.Errorf("waf: update ip set failed: %w", err)
	}
	return nil
}

func severityLabel(s *vault.Secret) string {
	if s.IsExpired() {
		return "expired"
	}
	return "expiring"
}
