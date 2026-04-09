// Package runner provides the top-level execution loop for vaultwatch.
//
// A Runner is constructed with a Config, Scanner, and Monitor. It calls
// Scanner.Scan periodically (according to Config.ScanInterval) and passes
// each discovered secret to Monitor.Check, which decides whether an alert
// notification should be sent.
//
// Typical usage:
//
//	r, err := runner.New(cfg, s, m)
//	if err != nil { ... }
//	if err := r.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package runner
