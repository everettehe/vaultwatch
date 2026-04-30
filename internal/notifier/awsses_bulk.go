package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// sesBulkClient defines the SES v2 bulk email interface used by SESSendBulkNotifier.
type sesBulkClient interface {
	SendBulkEmail(ctx context.Context, params *sesv2.SendBulkEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendBulkEmailOutput, error)
}

// SESSendBulkNotifier sends bulk templated emails via AWS SES v2.
type SESSendBulkNotifier struct {
	client       sesBulkClient
	fromAddress  string
	toAddresses  []string
	templateName string
}

// NewSESSendBulkNotifier creates a new SESSendBulkNotifier.
func NewSESSendBulkNotifier(fromAddress string, toAddresses []string, templateName, region string) (*SESSendBulkNotifier, error) {
	if fromAddress == "" {
		return nil, fmt.Errorf("ses_bulk: from address is required")
	}
	if len(toAddresses) == 0 {
		return nil, fmt.Errorf("ses_bulk: at least one to address is required")
	}
	if templateName == "" {
		return nil, fmt.Errorf("ses_bulk: template name is required")
	}
	if region == "" {
		return nil, fmt.Errorf("ses_bulk: region is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("ses_bulk: failed to load AWS config: %w", err)
	}

	return newSESSendBulkNotifierWithClient(sesv2.NewFromConfig(cfg), fromAddress, toAddresses, templateName), nil
}

func newSESSendBulkNotifierWithClient(client sesBulkClient, fromAddress string, toAddresses []string, templateName string) *SESSendBulkNotifier {
	return &SESSendBulkNotifier{
		client:       client,
		fromAddress:  fromAddress,
		toAddresses:  toAddresses,
		templateName: templateName,
	}
}

// Notify sends a bulk templated email alert for the given secret.
func (n *SESSendBulkNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg := FormatMessage(secret)

	destinations := make([]types.BulkEmailEntry, 0, len(n.toAddresses))
	for _, addr := range n.toAddresses {
		destinations = append(destinations, types.BulkEmailEntry{
			Destination: &types.Destination{
				ToAddresses: []string{addr},
			},
			ReplacementTemplateData: aws.String(fmt.Sprintf(`{"message":%q,"path":%q}`, msg.Body, secret.Path)),
		})
	}

	_, err := n.client.SendBulkEmail(ctx, &sesv2.SendBulkEmailInput{
		FromEmailAddress: aws.String(n.fromAddress),
		DefaultContent: &types.BulkEmailContent{
			Template: &types.Template{
				TemplateName: aws.String(n.templateName),
				TemplateData: aws.String(`{"message":"vault secret expiring","path":"unknown"}`),
			},
		},
		BulkEmailEntries: destinations,
	})
	if err != nil {
		return fmt.Errorf("ses_bulk: failed to send bulk email: %w", err)
	}
	return nil
}
