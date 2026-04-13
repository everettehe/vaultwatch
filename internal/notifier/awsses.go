package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/yourusername/vaultwatch/internal/vault"
)

// sesClient is the interface for AWS SES operations used in tests.
type sesClient interface {
	SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
}

// SESNotifier sends notifications via AWS Simple Email Service.
type SESNotifier struct {
	client sesClient
	from   string
	to     []string
	region string
}

// NewSESNotifier creates a new SESNotifier using the ambient AWS config.
func NewSESNotifier(from string, to []string, region string) (*SESNotifier, error) {
	if from == "" {
		return nil, fmt.Errorf("ses: from address is required")
	}
	if len(to) == 0 {
		return nil, fmt.Errorf("ses: at least one recipient is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("ses: failed to load AWS config: %w", err)
	}

	return &SESNotifier{
		client: ses.NewFromConfig(cfg),
		from:   from,
		to:     to,
		region: region,
	}, nil
}

func newSESNotifierWithClient(client sesClient, from string, to []string) *SESNotifier {
	return &SESNotifier{client: client, from: from, to: to}
}

// Notify sends an SES email alert for the given secret.
func (n *SESNotifier) Notify(ctx context.Context, secret vault.Secret) error {
	msg := FormatMessage(secret)

	toAddrs := make([]string, len(n.to))
	copy(toAddrs, n.to)

	_, err := n.client.SendEmail(ctx, &ses.SendEmailInput{
		Source: aws.String(n.from),
		Destination: &types.Destination{
			ToAddresses: toAddrs,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    aws.String(msg.Subject),
				Charset: aws.String("UTF-8"),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data:    aws.String(msg.Body),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("ses: failed to send email: %w", err)
	}
	return nil
}
