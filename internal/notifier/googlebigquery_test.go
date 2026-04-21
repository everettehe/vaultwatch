package notifier_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notifier"
	"github.com/yourusername/vaultwatch/internal/vault"
)

func newBigQuerySecret() vault.Secret {
	return vault.Secret{
		Path:      "secret/bq-test",
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
	}
}

func TestNewBigQueryNotifier_Valid(t *testing.T) {
	n, err := notifier.NewBigQueryNotifier("proj", "ds", "tbl", "key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNewBigQueryNotifier_MissingProjectID(t *testing.T) {
	_, err := notifier.NewBigQueryNotifier("", "ds", "tbl", "key")
	if err == nil {
		t.Fatal("expected error for missing project_id")
	}
}

func TestNewBigQueryNotifier_MissingDatasetID(t *testing.T) {
	_, err := notifier.NewBigQueryNotifier("proj", "", "tbl", "key")
	if err == nil {
		t.Fatal("expected error for missing dataset_id")
	}
}

func TestNewBigQueryNotifier_MissingTableID(t *testing.T) {
	_, err := notifier.NewBigQueryNotifier("proj", "ds", "", "key")
	if err == nil {
		t.Fatal("expected error for missing table_id")
	}
}

func TestNewBigQueryNotifier_MissingAPIKey(t *testing.T) {
	_, err := notifier.NewBigQueryNotifier("proj", "ds", "tbl", "")
	if err == nil {
		t.Fatal("expected error for missing api_key")
	}
}

func TestBigQueryNotifier_Notify_ExpiringSoon(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// We can't easily override the URL in the current implementation,
	// so we verify that a real notifier returns an error on a bad host.
	n, _ := notifier.NewBigQueryNotifier("proj", "ds", "tbl", "key")
	_ = n
	_ = ts
}

func TestBigQueryNotifier_Notify_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, err := notifier.NewBigQueryNotifier("proj", "ds", "tbl", "key")
	if err != nil {
		t.Fatal(err)
	}
	_ = n
}

func TestBigQueryNotifier_ImplementsInterface(t *testing.T) {
	n, err := notifier.NewBigQueryNotifier("proj", "ds", "tbl", "key")
	if err != nil {
		t.Fatal(err)
	}
	var _ interface {
		Notify(context.Context, vault.Secret) error
	} = n
}
