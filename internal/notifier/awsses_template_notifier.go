package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// sesSendTemplatedEmailAPI abstracts the SES v2 SendEmail call for templated messages.
type sesSendTemplatedEmailAPI interface {
	SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

// SESSendTemplateNotifier sends notifications via AWS SES using a named template.
type SESSendTemplateNotifier struct {
	client   sesSendTemplatedEmailAPI
	from     string
	to       string
	template string
}

// NewSESSendTemplateNotifier constructs a notifier that uses an AWS SES template.
func NewSESSendTemplateNotifier(from, to, template, region string) (*SESSendTemplateNotifier, error) {
	if from == "" {
		return nil, fmt.Errorf("ses_send_template: from address is required")
	}
	if to == "" {
		return nil, fmt.Errorf("ses_send_template: to address is required")
	}
	if template == "" {
		return nil, fmt.Errorf("ses_send_template: template name is required")
	}
	if region == "" {
		return nil, fmt.Errorf("ses_send_template: region is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("ses_send_template: load AWS config: %w", err)
	}

	return &SESSendTemplateNotifier{
		client:   sesv2.NewFromConfig(cfg),
		from:     from,
		to:       to,
		template: template,
	}, nil
}

func newSESSendTemplateNotifierWithClient(client sesSendTemplatedEmailAPI, from, to, template string) *SESSendTemplateNotifier {
	return &SESSendTemplateNotifier{
		client:   client,
		from:     from,
		to:       to,
		template: template,
	}
}

// Notify sends an alert using the configured SES template.
func (n *SESSendTemplateNotifier) Notify(ctx context.Context, secret Secret) error {
	msg, _ := FormatMessage(secret)
	templateData := fmt.Sprintf(`{"subject":%q,"body":%q}`, msg.Subject, msg.Body)

	_, err := n.client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(n.from),
		Destination: &types.Destination{
			ToAddresses: []string{n.to},
		},
		Content: &types.EmailContent{
			Template: &types.Template{
				TemplateName: aws.String(n.template),
				TemplateData: aws.String(templateData),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ses_send_template: send email: %w", err)
	}
	return nil
}
