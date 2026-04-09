package scanner

import (
	"context"
	"errors"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

// mockLister implements VaultLister for testing.
type mockLister struct {
	responses map[string]*vaultapi.Secret
	err       error
}

func (m *mockLister) List(_ context.Context, p string) (*vaultapi.Secret, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.responses[p], nil
}

func TestNew_NilLister(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil lister")
	}
}

func TestScan_FlatPaths(t *testing.T) {
	lister := &mockLister{
		responses: map[string]*vaultapi.Secret{
			"secret/": {
				Data: map[string]interface{}{
					"keys": []interface{}{"db-password", "api-key"},
				},
			},
		},
	}
	s, _ := New(lister)
	results, err := s.Scan(context.Background(), "secret/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(results[0].Keys))
	}
}

func TestScan_NestedPaths(t *testing.T) {
	lister := &mockLister{
		responses: map[string]*vaultapi.Secret{
			"secret/": {
				Data: map[string]interface{}{
					"keys": []interface{}{"apps/"},
				},
			},
			"secret/apps": {
				Data: map[string]interface{}{
					"keys": []interface{}{"token"},
				},
			},
		},
	}
	s, _ := New(lister)
	results, err := s.Scan(context.Background(), "secret/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Path != "secret/apps" {
		t.Errorf("expected path secret/apps, got %s", results[0].Path)
	}
}

func TestScan_ListerError(t *testing.T) {
	lister := &mockLister{err: errors.New("vault unavailable")}
	s, _ := New(lister)
	_, err := s.Scan(context.Background(), "secret/")
	if err == nil {
		t.Fatal("expected error from lister")
	}
}

func TestScan_EmptyRoot(t *testing.T) {
	lister := &mockLister{
		responses: map[string]*vaultapi.Secret{},
	}
	s, _ := New(lister)
	results, err := s.Scan(context.Background(), "secret/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
