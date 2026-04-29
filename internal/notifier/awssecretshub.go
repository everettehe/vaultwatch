package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/securityhub"
	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// securityHubClient defines the interface for AWS Security Hub operations.
type securityHubClient interface {
	BatchImportFindings(ctx context.Context, params *securityhub.BatchImportFindingsInput, optFns ...func(*securityhub.Options)) (*securityhub.BatchImportFindingsOutput, error)
}

// SecurityHubNotifier sends Vault secret expiration findings to AWS Security Hub.
type SecurityHubNotifier struct {
	client    securityHubClient
	accountID string
	region    string
	productARN string
}

// NewSecurityHubNotifier creates a SecurityHubNotifier using default AWS credentials.
func NewSecurityHubNotifier(accountID, region string) (*SecurityHubNotifier, error) {
	if accountID == "" {
		return nil, fmt.Errorf("securityhub: account ID is required")
	}
	if region == "" {
		return nil, fmt.Errorf("securityhub: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("securityhub: failed to load AWS config: %w", err)
	}
	return newSecurityHubNotifierWithClient(securityhub.NewFromConfig(cfg), accountID, region), nil
}

func newSecurityHubNotifierWithClient(client securityHubClient, accountID, region string) *SecurityHubNotifier {
	return &SecurityHubNotifier{
		client:     client,
		accountID:  accountID,
		region:     region,
		productARN: fmt.Sprintf("arn:aws:securityhub:%s:%s:product/%s/default", region, accountID, accountID),
	}
}

// Notify imports a Security Hub finding for the given secret.
func (n *SecurityHubNotifier) Notify(ctx context.Context, s *vault.Secret) error {
	severity := types.SeverityLabelMedium
	if s.IsExpired() {
		severity = types.SeverityLabelCritical
	} else if s.DaysUntilExpiration() <= 7 {
		severity = types.SeverityLabelHigh
	}

	finding := types.AwsSecurityFinding{
		AwsAccountId:  aws.String(n.accountID),
		Description:   aws.String(FormatMessage(s)),
		GeneratorId:   aws.String("vaultwatch"),
		Id:            aws.String(fmt.Sprintf("vaultwatch/%s/%s", n.region, s.Path)),
		ProductArn:    aws.String(n.productARN),
		SchemaVersion: aws.String("2018-10-08"),
		Title:         aws.String(fmt.Sprintf("Vault secret expiring: %s", s.Path)),
		Severity:      &types.Severity{Label: severity},
		Types:         []string{"Software and Configuration Checks/Vault Secret Expiration"},
	}

	_, err := n.client.BatchImportFindings(ctx, &securityhub.BatchImportFindingsInput{
		Findings: []types.AwsSecurityFinding{finding},
	})
	if err != nil {
		return fmt.Errorf("securityhub: failed to import finding: %w", err)
	}
	return nil
}
