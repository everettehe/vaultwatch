package scanner

import (
	"context"
	"fmt"
	"path"

	vaultapi "github.com/hashicorp/vault/api"
)

// SecretMeta holds metadata about a discovered Vault secret path.
type SecretMeta struct {
	Path string
	Keys []string
}

// VaultLister is the interface used to list secret paths from Vault.
type VaultLister interface {
	List(ctx context.Context, path string) (*vaultapi.Secret, error)
}

// Scanner recursively discovers secret paths under a given prefix.
type Scanner struct {
	lister VaultLister
}

// New creates a new Scanner with the provided VaultLister.
func New(lister VaultLister) (*Scanner, error) {
	if lister == nil {
		return nil, fmt.Errorf("lister must not be nil")
	}
	return &Scanner{lister: lister}, nil
}

// Scan recursively lists all secret paths under the given root prefix.
// It returns a slice of SecretMeta for each leaf path discovered.
func (s *Scanner) Scan(ctx context.Context, root string) ([]SecretMeta, error) {
	var results []SecretMeta
	if err := s.walk(ctx, root, &results); err != nil {
		return nil, fmt.Errorf("scan failed at %q: %w", root, err)
	}
	return results, nil
}

func (s *Scanner) walk(ctx context.Context, prefix string, results *[]SecretMeta) error {
	secret, err := s.lister.List(ctx, prefix)
	if err != nil {
		return err
	}
	if secret == nil || secret.Data == nil {
		return nil
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil
	}

	var leafKeys []string
	for _, k := range keys {
		key, ok := k.(string)
		if !ok {
			continue
		}
		if len(key) > 0 && key[len(key)-1] == '/' {
			// Directory — recurse
			if err := s.walk(ctx, path.Join(prefix, key), results); err != nil {
				return err
			}
		} else {
			leafKeys = append(leafKeys, key)
		}
	}

	if len(leafKeys) > 0 {
		*results = append(*results, SecretMeta{Path: prefix, Keys: leafKeys})
	}
	return nil
}
