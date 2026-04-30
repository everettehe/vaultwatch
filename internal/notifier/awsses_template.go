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

type sesTemplateClient interface {
	SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

// SESTemplateNotifier sends alerts via AWS SES using a pre-defined template.
type SESTemplateNotifier struct {
	client       sesTemplateClient
	from         string
	to           string
	templateName string
}

// NewSESTemplateNotifier creates a new SESTemplateNotifier using ambient AWS credentials.
func NewSESTemplateNotifier(from, to, templateName, region string) (*SESTemplateNotifier, error) {
	if from == "" {
		return nil, fmt.Errorf("ses_template: from address is required")
	}
	if to == "" {
		return nil, fmt.Errorf("ses_template: to address is required")
	}
	if templateName == "" {
		return nil, fmt.Errorf("ses_template: template name is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("ses_template: failed to load AWS config: %w", err)
	}
	return newSESTemplateNotifierWithClient(sesv2.NewFromConfig(cfg), from, to, templateName)
}

func newSESTemplateNotifierWithClient(client sesTemplateClient, from, to, templateName string) (*SESTemplateNotifier, error) {
	if from == "" {
		return nil, fmt.Errorf("ses_template: from address is required")
	}
	if to == "" {
		return nil, fmt.Errorf("ses_template: to address is required")
	}
	if templateName == "" {
		return nil, fmt.Errorf("ses_template: template name is required")
	}
	return &SESTemplateNotifier{
		client:       client,
		from:         from,
		to:           to,
		templateName: templateName,
	}, nil
}

// Notify sends a templated SES email for the given secret.
func (n *SESTemplateNotifier) Notify(ctx context.Context, secret *vault.Secret) error {
	msg := FormatMessage(secret)
	templateData := fmt.Sprintf(`{"subject":%q,"body":%q}`, msg.Subject, msg.Body)
	_, err := n.client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(n.from),
		Destination: &types.Destination{
			ToAddresses: []string{n.to},
		},
		Content: &types.EmailContent{
			Template: &types.Template{
				TemplateName: aws.String(n.templateName),
				TemplateData: aws.String(templateData),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ses_template: failed to send email: %w", err)
	}
	return nil
}
