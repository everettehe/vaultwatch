package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Lister wraps a Vault client to implement the scanner.VaultLister interface.
type Lister struct {
	client *vaultapi.Client
}

// NewLister creates a Lister from an existing *vaultapi.Client.
func NewLister(client *vaultapi.Client) (*Lister, error) {
	if client == nil {
		return nil, fmt.Errorf("vault client must not be nil")
	}
	return &Lister{client: client}, nil
}

// List performs a LIST operation on the given Vault path and returns the raw secret.
func (l *Lister) List(ctx context.Context, path string) (*vaultapi.Secret, error) {
	secret, err := l.client.Logical().ListWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("listing path %q: %w", path, err)
	}
	return secret, nil
}
