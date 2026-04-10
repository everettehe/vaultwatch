package notifier

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yourusername/vaultwatch/internal/vault"
)

// PrometheusNotifier pushes metrics to a Prometheus Pushgateway.
type PrometheusNotifier struct {
	pushgatewayURL string
	job            string
	client         *http.Client
}

// NewPrometheusNotifier creates a new PrometheusNotifier.
// pushgatewayURL is the base URL of the Prometheus Pushgateway.
// job is the job label used when pushing metrics.
func NewPrometheusNotifier(pushgatewayURL, job string) (*PrometheusNotifier, error) {
	if pushgatewayURL == "" {
		return nil, fmt.Errorf("prometheus: pushgateway URL is required")
	}
	if job == "" {
		job = "vaultwatch"
	}
	return &PrometheusNotifier{
		pushgatewayURL: strings.TrimRight(pushgatewayURL, "/"),
		job:            job,
		client:         &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Notify pushes a metric to the Pushgateway representing days until expiration.
func (p *PrometheusNotifier) Notify(secret *vault.Secret) error {
	days := secret.DaysUntilExpiration()
	metricName := "vaultwatch_secret_days_until_expiration"
	path := sanitizeLabelValue(secret.Path)

	body := fmt.Sprintf(
		"# HELP %s Days until the Vault secret expires\n"+
			"# TYPE %s gauge\n"+
			"%s{path=%q} %d\n",
		metricName, metricName, metricName, path, days,
	)

	url := fmt.Sprintf("%s/metrics/job/%s/instance/%s",
		p.pushgatewayURL, p.job, path)

	resp, err := p.client.Post(url, "text/plain", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("prometheus: failed to push metric: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("prometheus: pushgateway returned status %d", resp.StatusCode)
	}
	return nil
}

// sanitizeLabelValue replaces characters that are invalid in Prometheus label values.
func sanitizeLabelValue(s string) string {
	return strings.NewReplacer("/", "_", " ", "_").Replace(strings.Trim(s, "/"))
}
