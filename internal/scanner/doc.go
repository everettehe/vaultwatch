// Package scanner provides recursive discovery of secret paths within a
// HashiCorp Vault instance. It walks a given root prefix and returns
// metadata for every leaf path found, which can then be passed to the
// monitor for expiration checking.
//
// Usage:
//
//	lister, _ := vault.NewLister(vaultClient)
//	s, _ := scanner.New(lister)
//	paths, err := s.Scan(ctx, "secret/")
package scanner
