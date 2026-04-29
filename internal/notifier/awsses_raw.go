package notifier

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// SESRawNotifier sends raw MIME email notifications via AWS SES.
type SESRawNotifier struct {
	client sesRawClient
	from   string
	to     string
	region string
}

type sesRawClient interface {
	SendRawEmail(ctx context.Context, params *ses.SendRawEmailInput, optFns ...func(*ses.Options)) (*ses.SendRawEmailOutput, error)
}

// NewSESRawNotifier creates a new SESRawNotifier using ambient AWS credentials.
func NewSESRawNotifier(from, to, region string) (*SESRawNotifier, error) {
	if from == "" {
		return nil, fmt.Errorf("ses_raw: from address is required")
	}
	if to == "" {
		return nil, fmt.Errorf("ses_raw: to address is required")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("ses_raw: failed to load AWS config: %w", err)
	}
	return &SESRawNotifier{
		client: ses.NewFromConfig(cfg),
		from:   from,
		to:     to,
		region: region,
	}, nil
}

func newSESRawNotifierWithClient(client sesRawClient, from, to string) *SESRawNotifier {
	return &SESRawNotifier{client: client, from: from, to: to}
}

// Notify sends a raw MIME email for the given secret event.
func (n *SESRawNotifier) Notify(ctx context.Context, secret *Secret) error {
	msg := FormatMessage(secret)
	raw := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain\r\n\r\n%s",
		n.from, n.to, msg.Subject, msg.Body,
	)
	_, err := n.client.SendRawEmail(ctx, &ses.SendRawEmailInput{
		RawMessage: &types.RawMessage{Data: []byte(raw)},
		Source:     aws.String(n.from),
		Destinations: []string{n.to},
	})
	if err != nil {
		return fmt.Errorf("ses_raw: failed to send email: %w", err)
	}
	return nil
}
