package notifier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

type mockSESSendTemplateClient struct {
	calledWith *sesv2.SendEmailInput
	err        error
}

func (m *mockSESSendTemplateClient) SendEmail(_ context.Context, params *sesv2.SendEmailInput, _ ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	m.calledWith = params
	return &sesv2.SendEmailOutput{}, m.err
}

func newSESSendTemplateSecret(days int) Secret {
	return &mockSecret{
		path:      "secret/db/password",
		expiry:    time.Now().Add(time.Duration(days) * 24 * time.Hour),
		expired:   days < 0,
		daysUntil: days,
	}
}

func TestNewSESSendTemplateNotifier_MissingFrom(t *testing.T) {
	_, err := NewSESSendTemplateNotifier("", "to@example.com", "MyTemplate", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestNewSESSendTemplateNotifier_MissingTo(t *testing.T) {
	_, err := NewSESSendTemplateNotifier("from@example.com", "", "MyTemplate", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing to")
	}
}

func TestNewSESSendTemplateNotifier_MissingTemplate(t *testing.T) {
	_, err := NewSESSendTemplateNotifier("from@example.com", "to@example.com", "", "us-east-1")
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestNewSESSendTemplateNotifier_MissingRegion(t *testing.T) {
	_, err := NewSESSendTemplateNotifier("from@example.com", "to@example.com", "MyTemplate", "")
	if err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestSESSendTemplateNotifier_Notify_ExpiringSoon(t *testing.T) {
	mock := &mockSESSendTemplateClient{}
	n := newSESSendTemplateNotifierWithClient(mock, "from@example.com", "to@example.com", "VaultAlert")

	if err := n.Notify(context.Background(), newSESSendTemplateSecret(5)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calledWith == nil {
		t.Fatal("expected SendEmail to be called")
	}
}

func TestSESSendTemplateNotifier_Notify_SendError(t *testing.T) {
	mock := &mockSESSendTemplateClient{err: errors.New("SES error")}
	n := newSESSendTemplateNotifierWithClient(mock, "from@example.com", "to@example.com", "VaultAlert")

	if err := n.Notify(context.Background(), newSESSendTemplateSecret(3)); err == nil {
		t.Fatal("expected error from SendEmail failure")
	}
}
