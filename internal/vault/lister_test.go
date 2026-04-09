package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func newTestVaultClient(t *testing.T, handler http.Handler) *vaultapi.Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	cfg := vaultapi.DefaultConfig()
	cfg.Address = ts.URL
	client, err := vaultapi.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	client.SetToken("test-token")
	return client
}

func TestNewLister_NilClient(t *testing.T) {
	_, err := NewLister(nil)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestNewLister_Valid(t *testing.T) {
	client := newTestVaultClient(t, http.NotFoundHandler())
	lister, err := NewLister(client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lister == nil {
		t.Fatal("expected non-nil lister")
	}
}

func TestLister_List_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"keys":["foo","bar/"]}}`))
	})
	client := newTestVaultClient(t, handler)
	lister, _ := NewLister(client)

	secret, err := lister.List(context.Background(), "secret/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secret == nil {
		t.Fatal("expected non-nil secret")
	}
}

func TestLister_List_ServerError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	client := newTestVaultClient(t, handler)
	lister, _ := NewLister(client)

	_, err := lister.List(context.Background(), "secret/")
	if err == nil {
		t.Fatal("expected error from server 500")
	}
}
