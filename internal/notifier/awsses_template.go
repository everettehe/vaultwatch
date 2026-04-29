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

// SESTemplateNotifier sends notifications via AWS SES using a pre-defined template.
type SESTemplateNotifier struct {
	client    sesTemplateClient
	from      string
	to        string
	template  string
	region    string
}

type sesTemplateClient interface {
	SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

// NewSESTemplateNotifier creates a new SESTemplateNotifier.
func NewSESTemplateNotifier(from, to, template, region string) (*SESTemplateNotifier, error) {
	if from == "" {
		return nil, fmt.Errorf("ses_template: from address is required")
	}
	if to == "" {
		return nil, fmt.Errorf("ses_template: to address is required")
	}
	if template == "" {
		return nil, fmt.Errorf("ses_template: template name is required")
	}
	if region == "" {
		return nil, fmt.Errorf("ses_template: region is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("ses_template: failed to load AWS config: %w", err)
	}
	return newSESTemplateNotifierWithClient(sesv2.NewFromConfig(cfg), from, to, template, region), nil
}

func newSESTemplateNotifierWithClient(client sesTemplateClient, from, to, template, region string) *SESTemplateNotifier {
	return &SESTemplateNotifier{client: client, from: from, to: to, template: template, region: region}
}

// Notify sends a templated email notification for the given secret.
func (n *SESTemplateNotifier) Notify(ctx context.Context, secret vault.Secret) error {
	msg, _ := FormatMessage(secret)
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(n.from),
		Destination: &types.Destination{
			ToAddresses: []string{n.to},
		},
		Content: &types.EmailContent{
			Template: &types.Template{
				TemplateName: aws.String(n.template),
				TemplateData: aws.String(fmt.Sprintf(`{"message":%q,"path":%q}`, msg.Body, secret.Path)),
			},
		},
	}
	_, err := n.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("ses_template: failed to send email: %w", err)
	}
	return nil
}
